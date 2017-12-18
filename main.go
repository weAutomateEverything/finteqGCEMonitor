package main

import (
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/rest"
	"time"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/service"
)

func main() {
	rest.Router()
	service.DoSelenium()
	for true {
		time.Sleep(10 * time.Minute)
	}
}
