package main

import (
	"fmt"
	"os"

	"github.com/rikut0904/starter/internal/interface/cli"
)

func main() {
	if err := cli.RootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
