package model

import (
	"time"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	ID string `yaml:"id"`

	Name        string `yaml:"name"`
	Regex       string `yaml:"regex"`
	Level       string `yaml:"level"`
	Description string `yaml:"description"`
}

type LogEvent struct {
	TimeStamp time.Time

	Level        string
	RuleID       string
	RuleName     string
	SourceFile   string
	OriginalLine string
}
