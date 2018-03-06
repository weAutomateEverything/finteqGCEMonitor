package gceSelenium

import (
	"github.com/tebeka/selenium"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/halSelenium"
	"os"
)

type Service interface {
	HandleSeleniumError(internal bool, err error)
	NewClient() error
	Driver() selenium.WebDriver
	WaitForWaitFor() error
}

type service struct {
	halSelenium halSelenium.Service
}

func NewService(alert alert.Service) Service {
	sel := halSelenium.NewChromeService(alert)
	return &service{sel}
}

func (s *service) HandleSeleniumError(internal bool, err error) {
	s.halSelenium.HandleSeleniumError(internal, err)
}

func (s *service) Driver() selenium.WebDriver {
	return s.halSelenium.Driver()
}

func (s *service) WaitForWaitFor() error {
	return s.halSelenium.Driver().Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByID, "ModalCalLabel")
		if err != nil {
			return true, nil
		}
		r, err := elem.IsDisplayed()
		return !r, nil
	})
}

func (s *service) NewClient() error {
	return s.halSelenium.NewClient(seleniumServer())
}

func seleniumServer() string {
	return os.Getenv("SELENIUM_SERVER")
}
