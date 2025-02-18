package middleware

import (
	"context"
	"errors"
	"fmt"
	"rss/internal/commands"
	"rss/internal/database"
	"rss/internal/state"
)

func MiddlewareLoggedIn(handler func(s *state.State, cmd commands.Command, user database.User) error) func(*state.State, commands.Command) error {
	return func(s *state.State, cmd commands.Command) error {

		if s.Config.CurrentUserName == "" {
			return errors.New("no user is logged in")
		}

		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("not logged in, %w", err)
		}

		return handler(s, cmd, user)
	}
}
