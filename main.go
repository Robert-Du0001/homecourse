package main

import (
	"homecourse/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}
