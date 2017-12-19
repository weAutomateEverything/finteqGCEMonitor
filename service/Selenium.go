package service

import (
	"github.com/tebeka/selenium"
	"encoding/base64"
	"net/http"
	"encoding/json"
	"bytes"
	"log"
	"io/ioutil"
)

type alertMessage struct {
	Message, Image string
}

type inwardService struct {
	service, subservice, destination, status string
}

func DoSelenium() {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	caps["chrome.switches"] = []string{"--ignore-certificate-errors"}

	if webDriver, err = selenium.NewRemote(caps, seleniumServer()); err != nil {
		handleSeleniumError(err, nil)
		return
	}

	defer webDriver.Quit()

	err = webDriver.Get(endpoint())
	if err != nil {
		handleSeleniumError(err, webDriver)
	}


	err = waitForWaitFor(webDriver)

	if err != nil {
		handleSeleniumError(err, webDriver)
		return
	}



	doInwardCheck(webDriver)
	checkServices(webDriver)

}

func waitForWaitFor(webDriver selenium.WebDriver) error {
	return webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByID, "ModalCalLabel")
		if err != nil {
			return true, nil
		}
		r, err := elem.IsDisplayed()
		return !r, nil
	})
}

func handleSeleniumError(err error, driver selenium.WebDriver) {
	if driver == nil {
		sendError(err.Error(), nil)
		return
	}
	bytes, error := driver.Screenshot()
	if error != nil {
		// Couldnt get a screenshot - lets end the original error
		sendError(err.Error(), nil)
		return
	}
	sendError(err.Error(), bytes)
}

func sendError(message string, image []byte) {
	a := alertMessage{Message: message}
	if image != nil {
		a.Image = base64.StdEncoding.EncodeToString(image)
	}

	request, _ := json.Marshal(a)

	response, err := http.Post(errorEndpoint(), "application/json", bytes.NewReader(request))
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Println(ioutil.ReadAll(response.Body))
	}

}

func endpoint() string {
	return "http://c1592023:trendweb@10.187.5.61/GCEControlCentre/MonitorOutwardServices.aspx"
}

func seleniumServer() string {
	return "http://card-devops-selenium-service.legion.sbsa.local/wd/hub"
}

func errorEndpoint() string {
	return "http://localhost:8000/alert/image"
}
