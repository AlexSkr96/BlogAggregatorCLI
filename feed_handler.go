package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
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
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rssFeed := RSSFeed{}
	err = xml.Unmarshal(body, &rssFeed)
	if err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for _, rssItem := range rssFeed.Channel.Item {
		rssItem.Title = html.UnescapeString(rssItem.Title)
		rssItem.Description = html.UnescapeString(rssItem.Description)
	}

	return &rssFeed, nil
}

func HandlerAggregate(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("expecting time between requests as argument, e.g. 1s, 1m, 1h")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeed(s)
	}
}

func HandlerNewFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}
	feedParams := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   sql.NullString{String: cmd.Args[0], Valid: true},
		Url:    sql.NullString{String: cmd.Args[1], Valid: true},
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	}
	feedID, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Printf("created feed %v with ID %v\n", cmd.Args[0], feedID)
	cmd.Args = []string{cmd.Args[1]}
	err = middlewareLoggedIn(HandlerNewFeedFollow)(s, cmd)
	if err != nil {
		return fmt.Errorf("error while creating feed follows: %v", err)
	}
	return nil
}

func HandlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Printf("  * Feed: %v\n", feed.Name.String)
		fmt.Printf("    - URL: %v\n", feed.Url.String)
		fmt.Printf("    - User: %v\n", feed.Name_2)
	}
	return nil
}

func scrapeFeed(s *state) error {
	feed, err := s.db.FetchNextFeed(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url.String)
	if err != nil {
		return err
	}
	for _, rssItem := range rssFeed.Channel.Item {
		postId, err := s.db.CreatePost(
			context.Background(),
			uuid.New(),
			rssItem.Title,
			rssItem.Link,
			rssItem.Description,
			time.Parse(rssItem.PubDate),
			feed.ID,
		)
		// If error is duplicate URL, it's fine, just skip
		if strings.Contains(err, "duplicate") {
			continue
		} else if err != nil {
			log.Fatal(err)
		}

	}
	//fmt.Printf("Feed %v:\n", rssFeed.Channel.Title)
	//for _, rssItem := range rssFeed.Channel.Item {
	//	fmt.Printf("  - %v\n", rssItem.Title)
	//}
	return nil
}
