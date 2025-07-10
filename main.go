package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/sanesdotio/gator/internal/config"
	"github.com/sanesdotio/gator/internal/database"
)

type state struct {
	db *database.Queries
	config *config.Config
}

func main() {
	// Initialize the configuration
	// Read the configuration from the config file
	cfg := config.Read();

	// Connect to the database using the configuration
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
		return
	}
	dbQueries := database.New(db)

	// Update the state with the database queries and configuration
	currentState := &state{
		db: dbQueries,
		config: cfg,
	}

	// Initialize the commands
	cmds := &commands {
		Commands: map[string]func(*state, command) error {},
	}

	// Register the commands
	if err := cmds.register("login", loginHandler); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("register", registerHandler); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("reset", resetHandler); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("users", usersHandler); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("agg", aggHandler); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("addfeed", isLoggedIn(addFeedHandler)); err != nil {
  		fmt.Printf("error registering command: %v", err)
 	}
	
	if err := cmds.register("feeds", feedsHandler); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("follow", isLoggedIn(followHandler)); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("following", isLoggedIn(followingHandler)); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("unfollow", isLoggedIn(unfollowHandler)); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if err := cmds.register("browse", isLoggedIn(browseHandler)); err != nil {
		fmt.Printf("error registering command: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("No arguments provided.\n")
		return
	} 

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		Name: cmdName,
		args: cmdArgs,
	}

	// Run the command
	if err := cmds.run(currentState, cmd); err != nil {
		fmt.Printf("error running command: %v", err)
	}

	
}