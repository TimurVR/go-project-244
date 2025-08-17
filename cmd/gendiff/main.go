package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "hexlet-path-size",
		Usage: "print size of a file or directory",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "help",
				Aliases: []string{"h"},
				Usage:   "show help",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Bool("help") {
				Info()
				return nil
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
func Info(){
	fmt.Println(`NAME:
   gendiff - Compares two configuration files and shows a difference.

USAGE:
   gendiff [global options]

GLOBAL OPTIONS:
   --help, -h                  show help
`)
}