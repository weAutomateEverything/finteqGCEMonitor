package selenium

import (
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	"github.com/pkg/errors"
	"fmt"
)

type Service interface {
	HandleSeleniumError(err error)
	WaitForWaitFor() error
	Driver() selenium.WebDriver
}

type chromeService struct {
	alert alert.Service
	driver selenium.WebDriver


}

func NewChromeService(service alert.Service, server string) Service {

	s :=  &chromeService{alert: service}
	err := s.newClient(server)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *chromeService) Driver() selenium.WebDriver{
	return s.driver
}
func (s *chromeService)WaitForWaitFor() error {
	return s.driver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByID, "ModalCalLabel")
		if err != nil {
			return true, nil
		}
		r, err := elem.IsDisplayed()
		return !r, nil
	})
}

func (s *chromeService) newClient(seleniumServer string)  error {
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	caps["chrome.switches"] = []string{"--ignore-certificate-errors"}

	var err error
	s.driver, err = selenium.NewRemote(caps, seleniumServer)
	if err != nil {
		return err
	}
	return nil

}

func (s chromeService) HandleSeleniumError(err error) {
	internal := true
	msg := err.Error()
	if err, ok := err.(*SeleniumnError); ok {
		internal = err.Internal
		msg = err.Message.Error()
	}

	if s.driver == nil {
		s.sendError(msg, nil, internal)
		return
	}
	bytes, error := s.driver.Screenshot()
	if error != nil {
		// Couldnt get a screenshot - lets end the original error
		s.sendError(msg, nil, internal)
		return
	}
	s.sendError(msg, bytes, internal)
}

func (s chromeService) sendError(message string, image []byte, internalError bool) error {

	if image != nil {
		if internalError {
			err := s.alert.SendImageToHeartbeatGroup(image)
			if err != nil {
				return err
			}
		} else {
			err := s.alert.SendImageToAlertGroup(image)
			if err != nil {
				return err
			}
		}
	}

	if internalError {
		s.alert.SendError(errors.New(message))
	} else {
		err := s.alert.SendAlert(message)
		if err != nil {
			return err
		}
	}
	return nil

}

type SeleniumnError  struct {
	Internal bool
	Message error
}

func (e *SeleniumnError) Error() string {
	return fmt.Sprintf("internal: %v, message: %v",e.Internal,e.Message)
}
