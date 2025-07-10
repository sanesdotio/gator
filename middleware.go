package main

import (
	"context"

	"github.com/sanesdotio/gator/internal/database"
)

func isLoggedIn(handler func(state *state, cmd command, user database.User) error) func(*state, command) error {
	return func(state *state, cmd command) error {
		// Check if the user is logged in
		user, err := state.db.GetUser(context.Background(), state.config.CurrentUserName)
		if err != nil {
			return err
		}

		// Call the handler with the user
		return handler(state, cmd, user)
	}
}