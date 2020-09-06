package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	BuildDir      string
	ComponentsDir string
	LayoutsDir    string
	PagesDir      string
	StaticsDir    string
}

// ReadConfig from path.
func ReadConfig(path string) Config {
	var c Config
	if _, err := toml.DecodeFile(path, &c); err != nil {
		log.Println("Could not read config file, continuing with defaults.")
	}

	if c.BuildDir == "" {
		c.BuildDir = "docs"
	}
	if c.ComponentsDir == "" {
		c.ComponentsDir = "components"
	}
	if c.LayoutsDir == "" {
		c.LayoutsDir = "layouts"
	}
	if c.PagesDir == "" {
		c.PagesDir = "pages"
	}
	if c.StaticsDir == "" {
		c.StaticsDir = "statics"
	}

	return c
}
