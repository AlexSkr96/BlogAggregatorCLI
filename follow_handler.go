package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/database"
)

func HandlerNewFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("invalid number of arguments, expected URL")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), sql.NullString{String: cmd.Args[0], Valid: true})
	if err != nil {
		return err
	}
	feedFollowParams := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	}
	feedFollowData, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		err2 := s.db.DeleteFeed(context.Background(), feed.ID)
		if err2 != nil {
			return fmt.Errorf("failed to delete feed %v: %v\nfailed to create feed follow: %v", feed.Name.String, err2, err)
		}
		return err
	}
	fmt.Printf("added feed %v for user %v", feedFollowData.Name_2.String, feedFollowData.Name)
	return nil
}

func HandlerGetFollowsForUser(s *state, cmd command, user database.User) error {
	followsData, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Subscriptions for %v:\n", user.Name)
	for _, followRow := range followsData {
		fmt.Printf("  * %v\n", followRow.FeedName.String)
	}
	return nil
}

func HandlerDeleteFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("feed url expected")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), sql.NullString{String: cmd.Args[0], Valid: true})
	if err != nil {
		return err
	}
	feedFollowParams := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}
	fmt.Printf("deleted feed %v for user %v", feed.Name.String, user.Name)
	return nil
}
