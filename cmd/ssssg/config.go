package main

import (
	"log"

	"github.com/BurntSushi/toml"

	"ssssg"
)

// ReadConfig from path.
func ReadConfig(path string) ssssg.Config {
	var c ssssg.Config
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
	if c.Data == nil {
		c.Data = map[string]string{}
	}

	return c
}
