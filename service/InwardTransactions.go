package service

import (
	"github.com/tebeka/selenium"
	"fmt"
)

func doInwardCheck(webDriver selenium.WebDriver) {
	getInwardData(webDriver)
	
}

func getInwardData(webDriver selenium.WebDriver) []inwardService {
	elem, err := webDriver.FindElement(selenium.ByPartialLinkText, "Service Options")
	if err != nil {
		handleSeleniumError(err, webDriver)
	}

	err = elem.Click()
	if err != nil {
		handleSeleniumError(err, webDriver)
	}

	webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByPartialLinkText, "INWARD SERVICE OPTIONS")
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	})

	elem, err = webDriver.FindElement(selenium.ByPartialLinkText, "INWARD SERVICE OPTIONS")
	if err != nil {
		handleSeleniumError(err, webDriver)
	}

	err = elem.Click()
	if err != nil {
		handleSeleniumError(err, webDriver)
	}
	v := checkInwardTable(webDriver)
	elem, err = webDriver.FindElement(selenium.ByPartialLinkText, "2")
	if err != nil {
		handleSeleniumError(err, webDriver)
	}
	err = elem.Click()
	if err != nil {
		handleSeleniumError(err, webDriver)
	}

	b := checkInwardTable(webDriver)
	for _, x := range b {
		v = append(v, x)
	}
	return v
}

func checkInwardTable(webDriver selenium.WebDriver) []inwardService {
	var v []inwardService
	webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByXPATH, "//table[@id='TABLEINWARDSERVICES']/tbody/tr[1]/td[13]")
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	})
	i := 1
	for i < 50 {
		var service string

		service, err := getTableElement(i, 2, webDriver)
		if err != nil {
			return v
		}
		subService, err := getTableElement(i, 3, webDriver)
		if err != nil {
			return v
		}
		destinationCode, err := getTableElement(i, 4, webDriver)
		if err != nil {
			return v
		}
		status, err := getTableElement(i, 13, webDriver)
		if err != nil {
			return v
		}
		v = append(v, inwardService{service: service, destination: destinationCode, subservice: subService, status: status})
		i++
	}
	return v
}

func getTableElement(row, column int, webDriver selenium.WebDriver) (string, error) {
	elem, err := webDriver.FindElement(selenium.ByXPATH, fmt.Sprintf("//table[@id='TABLEINWARDSERVICES']/tbody/tr[%v]/td[%v]", row, column))
	if err != nil {
		handleSeleniumError(err, webDriver)
		return "", err
	}
	return elem.Text()
}
