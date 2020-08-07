package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"ssssg"
)

const (
	BuildDirConfigName      = "buildDir"
	ComponentsDirConfigName = "componentsDir"
	LayoutsDirConfigName    = "layoutsDir"
	PagesDirConfigName      = "pagesDir"
	StaticsDirConfigName    = "staticsDir"
)

func main() {
	viper.SetConfigName("ssssg")
	viper.AddConfigPath(".")

	viper.SetDefault(BuildDirConfigName, "docs")
	viper.SetDefault(ComponentsDirConfigName, "components")
	viper.SetDefault(LayoutsDirConfigName, "layouts")
	viper.SetDefault(PagesDirConfigName, "pages")
	viper.SetDefault(StaticsDirConfigName, "statics")

	if err := viper.ReadInConfig(); err != nil {
		printFatal("Could not read config file.", err)
	}

	if err := ssssg.Build(ssssg.BuildOptions{
		BuildDir:      viper.GetString(BuildDirConfigName),
		ComponentsDir: viper.GetString(ComponentsDirConfigName),
		LayoutsDir:    viper.GetString(LayoutsDirConfigName),
		PagesDir:      viper.GetString(PagesDirConfigName),
		StaticsDir:    viper.GetString(StaticsDirConfigName),
	}); err != nil {
		printFatal("Error building site:", err)
	}
}

func printFatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}
