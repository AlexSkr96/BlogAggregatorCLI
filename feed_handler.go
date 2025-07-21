package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

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
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", rssFeed)
	return nil
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
		return fmt.Errorf("error while creatin feed follows: %v", err)
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
