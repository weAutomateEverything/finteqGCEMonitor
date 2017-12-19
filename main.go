package main

import (
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/rest"
	"time"
)

func main() {
	rest.Router()
	for true {
		time.Sleep(10 * time.Minute)
	}
}
