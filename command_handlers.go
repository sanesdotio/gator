package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sanesdotio/gator/internal/database"
)

func loginHandler(state *state, cmd command) error {
	// Check if the user provided a username argument
	if len(cmd.args) != 1 {
		os.Exit(1)
		return fmt.Errorf("login command requires a username argument")

	}

	// Check if the user exists in the database
	user, err := state.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Printf("User %s does not exist. Please register first.\n", cmd.args[0])
		os.Exit(1)
	}

	// Set the user in the state configuration
	if err := state.config.SetUser(user.Name); err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	
	fmt.Printf("%s logged in successfully\n", cmd.args[0])
	return nil
}

func registerHandler(state *state, cmd command) error {
	// Check if the user provided a username argument
	// If not, print an error message and exit
	if len(cmd.args) != 1 {
		fmt.Println("register command requires a username argument")
		os.Exit(1)
	}

	// Check if the user already exists in the database
	// If so, print an error message and exit
	if user, err := state.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		fmt.Printf("User %s already exists\n", user.Name)
		os.Exit(1)
	}

	// Create a new user in the database
	newUser, err := state.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	// Set the user in the state configuration
	if err := state.config.SetUser(newUser.Name); err != nil {
		return fmt.Errorf("error setting user in config: %w", err)
	}


	fmt.Printf("User %s created successfully\n", newUser.Name)
	fmt.Println(newUser)
	return nil

}

func resetHandler(state *state, cmd command) error {
	fmt.Println("Resetting application data...")
	err := state.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting database: %w", err)
	}
	fmt.Println("Application data reset successfully")
	os.Exit(0)
	return nil
}

func usersHandler(state *state, cmd command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	fmt.Println("Users:")
	for _, user := range users {
		if user.Name == state.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func aggHandler(state *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("agg command requires a time between requests argument")
		os.Exit(1)
	}
	
	timeBetweenRequests := cmd.args[0]
	
	duration, err := time.ParseDuration(timeBetweenRequests)
	if err != nil {
		return fmt.Errorf("error parsing time duration: %w", err)
	}
	
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(context.Background() ,state.db)
	}

}

func addFeedHandler(state *state, cmd command, user database.User) error {

	if len(cmd.args) != 2 {
		fmt.Println("add command requires a feed name and URL argument")
		os.Exit(1)
  	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	feed, err := state.db.CreateFeed(context.Background(), database.CreateFeedParams{
  		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	follow, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}

	fmt.Printf("Feed %s added successfully with URL %s\n", feed.Name, feed.Url)
	fmt.Printf("Followed feed %s\n", follow.FeedName)

	return nil
}

func feedsHandler(state *state, cmd command) error {

	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	
	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		feedOwner, err := state.db.GetFeedOwner(context.Background(), feed.ID)
		if err != nil {
			return fmt.Errorf("error getting feed owner: %w", err)
		}
		fmt.Printf("* %s, %s, %s\n", feed.Name, feed.Url, feedOwner)
	}
	return nil
}

func followHandler(state *state, cmd command, user database.User) error {

	if len(cmd.args) != 1 {
		fmt.Println("follow command requires a feed URL argument")
		os.Exit(1)
	}

	feedURL := cmd.args[0]

	feed, err := state.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("error getting feed by URL: %w", err)
	}

	follow, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}
	
	fmt.Printf("Followed feed %s by %s\n", follow.FeedName, follow.UserName)
	return nil
}

func followingHandler(state *state, cmd command, user database.User) error {

	// Get the feed follows for the current user
	feedFollows, err := state.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error getting feed follows for user: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("You are not following any feeds")
		return nil
	}

	fmt.Println("Following feeds:")
	for _, follow := range feedFollows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}

func unfollowHandler(state *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		fmt.Println("unfollow command requires a feed URL argument")
		os.Exit(1)
	}

	feedURL := cmd.args[0]

	err := state.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    feedURL,
	})
	if err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
	}

	fmt.Printf("Unfollowed feed with URL %s\n", feedURL)
	return nil
}

func browseHandler(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}