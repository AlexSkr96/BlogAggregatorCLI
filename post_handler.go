package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlexSkr96/BlogAggregatorCLI/internal/database"
)

func HandlerGetPosts(s *state, cmd command, user database.User) error {
	limit := int32(2)
	if len(cmd.Args) > 0 {
		// string to int
		parsedLimit, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(parsedLimit)
	}
	posts, err := s.db.GetPostsForUser(context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  limit,
		},
	)
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("%s:\n  %s\n", post.Title.String, post.Description.String)
	}
	return nil
}
