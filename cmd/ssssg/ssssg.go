package main

import (
	"flag"
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
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			printFatal("Could not read config file.", err)
		}
	}

	options := ssssg.BuildOptions{
		BuildDir:      viper.GetString(BuildDirConfigName),
		ComponentsDir: viper.GetString(ComponentsDirConfigName),
		LayoutsDir:    viper.GetString(LayoutsDirConfigName),
		PagesDir:      viper.GetString(PagesDirConfigName),
		StaticsDir:    viper.GetString(StaticsDirConfigName),
	}

	var err error
	flag.Parse()
	switch flag.Arg(0) {
	case "init":
		err = ssssg.Init(options)
	case "build":
		err = ssssg.Build(options)
	case "version":
		ssssg.Version()
	default:
		printFatal("Usage: ssssg init|build|version")
	}

	if err != nil {
		printFatal("Error:", err)
	}
}

func printFatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}
