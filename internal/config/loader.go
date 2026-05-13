package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
	"github.com/go-playground/validator/v10"
)

type fileConfig struct {
	Floors   int    `json:"Floors" validate:"gte=2"`
	Monsters int    `json:"Monsters" validate:"gte=0"`
	OpenAt   string `json:"OpenAt" validate:"required,clock"`
	Duration int    `json:"Duration" validate:"gt=0"`
}

type Loader struct{}

func (Loader) Load(path string) (domain.DungeonConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.DungeonConfig{}, fmt.Errorf("read config: %w", err)
	}

	var raw fileConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return domain.DungeonConfig{}, fmt.Errorf("parse config json: %w", err)
	}

	if err := validateRawConfig(raw); err != nil {
		return domain.DungeonConfig{}, err
	}

	openAt, err := domain.ParseClock(raw.OpenAt)
	if err != nil {
		return domain.DungeonConfig{}, fmt.Errorf("parse OpenAt: %w", err)
	}

	config := domain.DungeonConfig{
		Floors:   raw.Floors,
		Monsters: raw.Monsters,
		OpenAt:   openAt,
		Duration: raw.Duration,
	}
	return config, nil
}

func validateRawConfig(config fileConfig) error {
	validate := validator.New()
	if err := validate.RegisterValidation("clock", func(field validator.FieldLevel) bool {
		_, err := domain.ParseClock(field.Field().String())
		return err == nil
	}); err != nil {
		return fmt.Errorf("register config validator: %w", err)
	}

	if err := validate.Struct(config); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	return nil
}
