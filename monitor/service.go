package monitor

import (
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/cutofftimes"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/gceservices"
	selenium2 "github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/selenium"
	"github.com/zamedic/go2hal/alert"
	"os"
	"time"
)

type Service interface {
}

type service struct {
	alert            alert.Service
	cutoffStore      cutofftimes.Store
	seleniumEndpoint string
}

func NewService(alert alert.Service, cutoffStore cutofftimes.Store, seleniumEndpoint string) Service {
	s := &service{alert,cutoffStore,seleniumEndpoint}
	go func() {
		s.startMonitor()
	}()

	return s

}

func (s *service) startMonitor() {
	for true {
		s.doCheck()
		time.Sleep(10 * time.Minute)

	}

}

func (s *service) doCheck() {
	selenium := selenium2.NewChromeService(s.alert, s.seleniumEndpoint)
	defer selenium.Driver().Quit()

	err := selenium.Driver().Get(endpoint())
	if err != nil {
		selenium.HandleSeleniumError(&selenium2.SeleniumnError{true, err})
		return
	}

	err = selenium.WaitForWaitFor()

	if err != nil {
		selenium.HandleSeleniumError(&selenium2.SeleniumnError{true, err})
		return
	}

	service := gceservices.NewService(selenium, true)
	service.RunServiceCheck()

	service = gceservices.NewService(selenium, false)
	service.RunServiceCheck()

	cutoff := cutofftimes.NewService(s.cutoffStore, selenium, true)
	cutoff.DoCheck()

	cutoff = cutofftimes.NewService(s.cutoffStore, selenium, false)
	cutoff.DoCheck()

}

func endpoint() string {
	return os.Getenv("gce_endpoint")
}
