package main

import (
	"flag"
	"fmt"
	"os"

	"ssssg"
)

func main() {
	c := ReadConfig("./ssssg.toml")

	options := ssssg.BuildOptions{
		BuildDir:      c.BuildDir,
		ComponentsDir: c.ComponentsDir,
		LayoutsDir:    c.LayoutsDir,
		PagesDir:      c.PagesDir,
		StaticsDir:    c.StaticsDir,
	}

	var err error
	flag.Parse()
	switch flag.Arg(0) {
	case "init":
		err = ssssg.Init(options)
	case "build":
		err = ssssg.Build(options)
	case "serve":
		err = ssssg.Serve(c.BuildDir)
	case "version":
		ssssg.Version()
	case "watch":
		err = ssssg.Watch(options)
	default:
		printFatal("Usage: ssssg init|build|watch|serve|version")
	}

	if err != nil {
		printFatal("Error:", err)
	}
}

func printFatal(v ...interface{}) {
	fmt.Println(v...)
	os.Exit(1)
}
