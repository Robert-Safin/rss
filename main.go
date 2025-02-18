package main

import (
	"database/sql"
	"fmt"
	"os"
	"rss/internal/config"
	"rss/internal/database"
	"rss/internal/state"

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
	commands := state.Commands{
		Commands: make(map[string]func(*state.State, state.Command) error),
	}

	commands.Register("login", state.HandlerLogin)
	commands.Register("register", state.HandlerRegister)
	commands.Register("reset", state.HandlerReset)
	commands.Register("users", state.HandlerUsers)
	commands.Register("agg", state.HandlerAgg)

	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no command given")
		os.Exit(1)
	}
	commandName := args[0]
	args = args[1:]

	newCommand := state.Command{
		Name: commandName,
		Args: args,
	}

	err = commands.Run(&appState, newCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
