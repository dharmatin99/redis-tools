package main

import (
	"errors"
	"log"
	"os"
	"sort"

	"github.com/dharmatin99/redis-tools/command"
	"github.com/dharmatin99/redis-tools/lib"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "Redis Tools",
		Usage:       "Managing your redis apps",
		Description: "Redis tool for managing ee",
		Commands: []*cli.Command{
			{
				Name:    "copy",
				Aliases: []string{"c"},
				Usage:   "copy redis data from source to target",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "source_addr", Aliases: []string{"sa"}, Usage: "source address <host:port>"},
					&cli.StringFlag{Name: "target_addr", Aliases: []string{"ta"}, Usage: "target address <host:port>"},
					&cli.IntFlag{Name: "db", Aliases: []string{"n"}, Usage: "database target to copy"},
					&cli.StringFlag{Name: "pattern", Aliases: []string{"p"}, Usage: "key pattern, default(*)", Value: "*"},
				},
				Action: func(c *cli.Context) error {
					log.Print("copy command started")
					if c.String("sa") == "" {
						return errors.New("source address can't be empty")
					}

					if c.String("ta") == "" {
						return errors.New("target address can't be empty")
					}

					if c.String("db") == "" {
						return errors.New("please define your db")
					}
					copier := &command.Copier{
						Ctx:          c.Context,
						SourceClient: lib.CreateRedisClient(c.String("sa"), c.Int("db")),
						TargetClient: lib.CreateRedisClient(c.String("ta"), c.Int("db")),
						Pattern:      c.String("p"),
					}
					return copier.Copy()
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
