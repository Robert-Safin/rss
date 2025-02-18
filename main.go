package main

import (
	"database/sql"
	"fmt"
	"os"
	"rss/internal/commands"
	"rss/internal/config"
	"rss/internal/database"
	"rss/internal/state"
	"rss/middleware"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://rob:@localhost:5432/gator?sslmode=disable")
	dbQueries := database.New(db)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config, err := config.Read()
	if err != nil {
		fmt.Println("config error")
		os.Exit(1)
	}
	appState := state.State{Config: &config, Db: dbQueries}
	cmds := commands.Commands{
		Commands: make(map[string]func(*state.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerReset)
	cmds.Register("users", commands.HandlerUsers)
	cmds.Register("agg", commands.HandlerAgg)
	cmds.Register("addfeed", middleware.MiddlewareLoggedIn(commands.HandlerAddFeed))
	cmds.Register("feeds", commands.HandlerListFeeds)
	cmds.Register("follow", middleware.MiddlewareLoggedIn(commands.HandlerFollow))
	cmds.Register("following", middleware.MiddlewareLoggedIn(commands.HandlerFollowing))

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no command given")
		os.Exit(1)
	}
	commandName := args[0]
	args = args[1:]

	newCommand := commands.Command{
		Name: commandName,
		Args: args,
	}

	err = cmds.Run(&appState, newCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
