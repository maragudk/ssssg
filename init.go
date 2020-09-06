package ssssg

import (
	"os"

	"ssssg/assets"
	"ssssg/errors2"
)

// Init a new site, creating the necessary directories and default assets.
func Init(config Config) error {
	for _, dir := range []string{config.BuildDir, config.ComponentsDir, config.LayoutsDir, config.PagesDir, config.StaticsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors2.Wrap(err, "could not create directory %v", dir)
		}
	}

	if err := assets.RestoreAssets(".", ""); err != nil {
		return errors2.Wrap(err, "could not copy default files")
	}
	return nil
}
