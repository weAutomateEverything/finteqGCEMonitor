package gceSelenium

import (
	"github.com/zamedic/go2hal/halSelenium"
	"github.com/zamedic/go2hal/alert"
	"github.com/tebeka/selenium"
)

type Service interface{
	HandleSeleniumError(internal bool, err error)
	Driver() selenium.WebDriver
	WaitForWaitFor() error
}

type service struct{
	halSelenium halSelenium.Service
}

func NewService(alert alert.Service, seleniumEndpoint string)Service{
	sel := halSelenium.NewChromeService(alert,seleniumEndpoint)
	return &service{sel}
}

func(s *service)HandleSeleniumError(internal bool, err error){
	s.halSelenium.HandleSeleniumError(internal,err)
}

func (s * service)Driver() selenium.WebDriver{
	return s.halSelenium.Driver()
}

func (s *service)WaitForWaitFor() error {
	return s.halSelenium.Driver().Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByID, "ModalCalLabel")
		if err != nil {
			return true, nil
		}
		r, err := elem.IsDisplayed()
		return !r, nil
	})
}



