package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sanesdotio/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}



func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}

	html.UnescapeString(feed.Channel.Title)
	html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Items {
		html.UnescapeString(feed.Channel.Items[i].Title)
		html.UnescapeString(feed.Channel.Items[i].Description)
	}
	
	return &feed, nil
}

func scrapeFeeds(ctx context.Context, db *database.Queries) error {
	nextFeed, err := db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,}

	if err := db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: currentTime,
		UpdatedAt:     time.Now(),
	}); err != nil {
		return err
	}

	feed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Items {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: currentTime,
			FeedID:    nextFeed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	return nil
}