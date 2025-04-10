package main

import (
	"github.com/arrow2nd/memento/app"
)

var (
	appName = "memento"
	version = "develop"
)

func main() {
	app := app.New(appName, version)
	app.Run()
}
