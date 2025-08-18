package main

import (
	"code"
	"context"
	"fmt"
	"log"
	"os"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() == 0 {
				return fmt.Errorf("path argument is required")
			}
			if cmd.NArg() < 2 {
				return fmt.Errorf("two file paths are required")
			}
			path1 := cmd.Args().First()
			path2 := cmd.Args().Get(1)
			map1 := code.Parsing(path1)
			map2 := code.Parsing(path2)
			fmt.Println(code.GenDiff(map1, map2))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
