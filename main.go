//go:generate go-winres make --product-version=git-tag

package main

import (
	"github.com/arrow2nd/memento/app"
)

func main() {
	app := app.New()
	app.Run()
}
