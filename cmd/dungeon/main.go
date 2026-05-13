package main

import (
	"os"

	"github.com/6ermvH/yadro-dungeon-task/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr))
}
