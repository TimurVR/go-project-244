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
		Name:  "gendiff",
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
			if cmd.NArg() == 0 {
				return fmt.Errorf("path argument is required")
			}

			path1 := cmd.Args().First()
			path2 := cmd.Args().Get(1)
			fmt.Println(path1,path2)
			map1:=code.Parsing(path1)
			map2:=code.Parsing(path2)
			fmt.Println(code.GenDiff(map1,map2))
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
func Info(){
	fmt.Print(`NAME:
   gendiff - Compares two configuration files and shows a difference.

USAGE:
   gendiff [global options]

GLOBAL OPTIONS:
   --format string, -f string  output format (default: "stylish")
   --help, -h                  show help`)
}