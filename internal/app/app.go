package app

import (
	"fmt"
	"io"

	"github.com/6ermvH/yadro-dungeon-task/internal/config"
	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
	"github.com/6ermvH/yadro-dungeon-task/internal/output"
	"github.com/6ermvH/yadro-dungeon-task/internal/parser"
	"github.com/6ermvH/yadro-dungeon-task/internal/simulator"
)

const (
	defaultConfigPath = "task/config.json"
	defaultEventsPath = "task/events"
)

func Run(args []string, stdout, stderr io.Writer) int {
	configPath, eventsPath, err := paths(args)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}

	configLoader := config.Loader{}
	dungeonConfig, err := configLoader.Load(configPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "config error: %v\n", err)
		return 1
	}

	eventParser := parser.EventParser{}
	events, err := eventParser.ParseFile(eventsPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "events error: %v\n", err)
		return 1
	}

	dungeon := domain.NewDungeon(dungeonConfig)
	result := simulator.New(dungeon).Run(events)

	formatter := output.Formatter{}
	if err := formatter.Write(result, stdout); err != nil {
		_, _ = fmt.Fprintf(stderr, "output error: %v\n", err)
		return 1
	}

	return 0
}

func paths(args []string) (string, string, error) {
	switch len(args) {
	case 0:
		return defaultConfigPath, defaultEventsPath, nil
	case 2:
		return args[0], args[1], nil
	default:
		return "", "", fmt.Errorf("usage: dungeon [config_path events_path]")
	}
}
