package main

import (
	hybs "github.com/hayabusa-cloud/hayabusa"
	"github.com/hayabusa-cloud/hybs-server/application"
)

func main() {
	hybs.StartService(application.Up, application.Down)
}
