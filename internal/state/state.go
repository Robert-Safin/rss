package state

import (
	"context"
	"errors"
	"fmt"
	"os"
	"rss/internal/config"
	"rss/internal/database"
	"rss/internal/rssfeed"
	"time"

	"github.com/google/uuid"
)

type State struct {
	Config *config.Config
	Db     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	if c.Commands == nil {
		c.Commands = make(map[string]func(*State, Command) error)
	}

	if _, present := c.Commands[name]; present {
		return errors.New("command name already registered")
	}

	c.Commands[name] = f
	return nil
}
func (c *Commands) Run(s *State, cmd Command) error {
	handler := c.Commands[cmd.Name]
	if handler == nil {
		return fmt.Errorf("%v command does not exist", cmd.Name)
	}

	if err := handler(s, cmd); err != nil {
		return err
	}

	return nil
}

func HandlerLogin(state *State, command Command) error {
	if len(command.Args) == 0 {
		return errors.New(" handler expects a single argument, the username")
	}

	_, err := state.Db.GetUser(context.Background(), command.Args[0])
	if err != nil {
		fmt.Println("username doesn't exist")
		os.Exit(1)
	}

	err = state.Config.SetUser(command.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username set to: ", command.Args[0])
	return nil
}

func HandlerRegister(state *State, command Command) error {

	if len(command.Args) == 0 {
		return errors.New("no username given")
	}
	_, err := state.Db.GetUser(context.Background(), command.Args[0])
	if err == nil {
		fmt.Println("username taken")
		os.Exit(1)
	}

	_, err = state.Db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      command.Args[0],
	})

	if err != nil {
		return fmt.Errorf("error crateing user %w", err)
	} else {
		fmt.Printf("created user %v", command.Args[0])
		state.Config.SetUser(command.Args[0])
	}

	return nil
}

func HandlerReset(state *State, command Command) error {

	err := state.Db.DeleteManyUsers(context.Background())

	if err != nil {
		return fmt.Errorf("Failed to delete all suers: %w", err)
	}

	fmt.Println("Reset users")

	return nil
}

func HandlerUsers(state *State, command Command) error {
	users, err := state.Db.GetAllUsers(context.Background())

	if err != nil {
		return fmt.Errorf("Error checking users: %w", err)
	}

	for _, user := range users {
		if user.Name == state.Config.CurrentUserName {
			fmt.Printf("# %v (current) \n", user.Name)
		} else {
			fmt.Printf("# %v \n", user.Name)
		}
	}

	return nil
}

func HandlerAgg(state *State, command Command) error {
	feed, err := rssfeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func HandlerAddFeed(state *State, command Command) error {
	if len(command.Args) < 2 {
		return errors.New("no argument provided")
	}

	user, err := state.Db.GetUser(context.Background(), state.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("logged in as invalid user: %w", err)
	}

	feed, err := state.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      command.Args[0],
		Url:       command.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %w", err)
	}

	fmt.Printf("Created new feed: %v by user %v", feed.Name, user.Name)

	return nil
}

func HandlerListFeeds(state *State, command Command) error {

	feedRecords, err := state.Db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}

	for i, item := range feedRecords {
		user, err := state.Db.FindFeedUser(context.Background(), item.UserID)
		if err != nil {
			return fmt.Errorf("error fetching feed's author: %w", err)
		}
		fmt.Printf("Name: %v, URL: %v, Username: %v \n", feedRecords[i].Name, feedRecords[i].Url, user.Name)
	}

	return nil
}
