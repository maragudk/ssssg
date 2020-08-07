package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"ssssg"
)

const (
	LayoutsDirConfigName = "layoutsDir"
	PagesDirConfigName   = "pagesDir"
	BuildDirConfigName   = "buildDir"
)

func main() {
	viper.SetConfigName("ssssg")
	viper.AddConfigPath(".")

	viper.SetDefault(LayoutsDirConfigName, "layouts")
	viper.SetDefault(PagesDirConfigName, "pages")
	viper.SetDefault(BuildDirConfigName, "docs")

	if err := viper.ReadInConfig(); err != nil {
		printFatal("Could not read config file.", err)
	}

	if err := ssssg.Build(ssssg.BuildOptions{
		LayoutsDir: viper.GetString(LayoutsDirConfigName),
		PagesDir:   viper.GetString(PagesDirConfigName),
		BuildDir:   viper.GetString(BuildDirConfigName),
	}); err != nil {
		printFatal("Error building site:", err)
	}
}

func printFatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}
