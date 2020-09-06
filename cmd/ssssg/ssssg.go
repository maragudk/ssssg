package main

import (
	"flag"
	"fmt"
	"os"

	"ssssg"
)

func main() {
	c := ReadConfig("./ssssg.toml")

	var err error
	flag.Parse()
	switch flag.Arg(0) {
	case "init":
		err = ssssg.Init(c)
	case "build":
		err = ssssg.Build(c)
	case "serve":
		err = ssssg.Serve(c.BuildDir)
	case "version":
		ssssg.Version()
	case "watch":
		err = ssssg.Watch(c)
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
