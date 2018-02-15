package monitor

import (
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/cutofftimes"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/gceSelenium"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/gceservices"
	"github.com/zamedic/go2hal/alert"
	"os"
	"time"
)

type Service interface {
}

type service struct {
	alert            alert.Service
	cutoffStore      cutofftimes.Store
}

func NewService(alert alert.Service, cutoffStore cutofftimes.Store) Service {
	s := &service{alert, cutoffStore}
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
	selenium := gceSelenium.NewService(s.alert)

	err := selenium.NewClient()
	if err != nil {
		s.alert.SendError(err)
		return
	}
	defer selenium.Driver().Quit()

	err = selenium.Driver().Get(endpoint())
	if err != nil {
		selenium.HandleSeleniumError(true, err)
		return
	}

	err = selenium.WaitForWaitFor()

	if err != nil {
		selenium.HandleSeleniumError(true, err)
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
