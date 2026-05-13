package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
)

type EventParser struct{}

func (p EventParser) ParseFile(path string) ([]domain.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open events file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	var events []domain.Event
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		event, err := p.ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan events file: %w", err)
	}

	return events, nil
}

func (EventParser) ParseLine(line string) (domain.Event, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return domain.Event{}, fmt.Errorf("expected at least 3 fields")
	}

	timeField := fields[0]
	if len(timeField) != len("[HH:MM:SS]") || !strings.HasPrefix(timeField, "[") || !strings.HasSuffix(timeField, "]") {
		return domain.Event{}, fmt.Errorf("invalid time field %q", timeField)
	}

	eventTime, err := domain.ParseClock(strings.Trim(timeField, "[]"))
	if err != nil {
		return domain.Event{}, err
	}

	playerID, err := strconv.Atoi(fields[1])
	if err != nil {
		return domain.Event{}, fmt.Errorf("invalid player id %q: %w", fields[1], err)
	}

	eventID, err := strconv.Atoi(fields[2])
	if err != nil {
		return domain.Event{}, fmt.Errorf("invalid event id %q: %w", fields[2], err)
	}

	eventType, err := domain.ParseEventType(eventID)
	if err != nil {
		return domain.Event{}, err
	}

	return domain.Event{
		Time:     eventTime,
		PlayerID: playerID,
		Type:     eventType,
		Param:    strings.Join(fields[3:], " "),
	}, nil
}
