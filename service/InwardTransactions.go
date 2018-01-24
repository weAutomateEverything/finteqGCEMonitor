package service

import (
	"github.com/tebeka/selenium"
	"fmt"
	"strings"
	"github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/database"
	"strconv"
	"bufio"
	"io"
	"log"
	"bytes"
)

var sodOk = map[string]struct{}{"SOD : ACK RECEIVED": {}}
var eodOk = map[string]struct{}{"EOD : ACK RECEIVED": {}}

func doCheck(webDriver selenium.WebDriver, inward bool) {
	v := getData(webDriver, inward)
	var e []string
	for _, x := range v {
		if database.CutoffExists(x.service+x.subservice, x.destination) {
			if database.IsInStartOfDay(x.service+x.subservice, x.destination) {
				_, ok := sodOk[x.status]
				if !ok {
					e = append(e, fmt.Sprintf("invalid status for service %v, sub service %v, status %v", x.service+x.subservice, x.destination, x.status))
				}
			}
		} else {
			log.Printf("No records found for service %v, subservice %v", x.service+x.subservice, x.destination)
		}
	}
	if len(e) > 0 {
		b := bytes.Buffer{}
		for _,s := range e {
			b.WriteString(s)
			b.WriteString("\n")
		}
		sendError(b.String(),nil,false)
	}
}

/*
ParseInwardCutttoffTimes parses the inward data to setup the cuttoff times

For example
C2;BI;08H00;19H00;Mon - Fri;08H00;13H00;Sat - Sun

The above has 2 cuttoff times.
Start of day is 08:00 to 19:00 monday to Friday and 08:00 to 13:00 saturday and sunday.

 */
func ParseInwardCutttoffTimes(i io.Reader) {
	scanner := bufio.NewScanner(i)

	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Parsing Line: %v",line)
		tokens := strings.Split(line, ";")
		sodHour := tokens[2]
		sodHour = strings.TrimSpace(sodHour)
		if len(sodHour) == 0 {
			continue
		}
		parseBlock(tokens[0], tokens[1], tokens[2], tokens[3], tokens[4])
		parseBlock(tokens[0], tokens[1], tokens[5], tokens[6], tokens[7])
	}

}

func parseBlock(service, subservice, sodTime, eodTime, days string) {
	sodTime = strings.TrimSpace(sodTime)
	eodTime = strings.TrimSpace(eodTime)

	if len(sodTime) == 0 {
		return
	}

	sodTime = strings.Replace(sodTime, "A ", "", 1)
	eodTime = strings.Replace(eodTime, "A ", "", 1)

	days = strings.TrimSpace(days)
	days = strings.Replace(days, "(ph)", "", -1)

	c := database.CutoffTime{Service: service, SubService: subservice}

	sod := strings.Split(sodTime, "H")
	c.SodHour, _ = strconv.Atoi(sod[0])
	c.SodMinute, _ = strconv.Atoi(sod[1])

	eod := strings.Split(eodTime, "H")
	c.EodHour, _ = strconv.Atoi(eod[0])
	c.EodMinute, _ = strconv.Atoi(eod[1])

	if days == "Mon - Sun" {
		i := 0
		for i < 7 {
			c.DayOfWeek = i
			database.SaveCutoff(c)
			i++
		}
		return
	}

	if days == "Mon - Fri" {
		i := 1
		for i < 6 {
			c.DayOfWeek = i
			database.SaveCutoff(c)
			i++
		}
		return
	}

	//If its not Monday - Sunday or Monday to Friday, it must be Sat - Sun

	c.DayOfWeek = 6
	database.SaveCutoff(c)

	c.DayOfWeek = 0
	database.SaveCutoff(c)
}

func getData(webDriver selenium.WebDriver, inward bool) []inwardService {

	err := webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByPartialLinkText, "Service Options")
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	})

	if err != nil {
		handleSeleniumError(err, webDriver)
		return nil
	}

	elem, err := webDriver.FindElement(selenium.ByPartialLinkText, "Service Options")
	if err != nil {
		handleSeleniumError(err, webDriver)
		return nil
	}

	err = elem.Click()
	if err != nil {
		handleSeleniumError(err, webDriver)
		return nil
	}

	err = waitForWaitFor(webDriver)
	if err != nil {
		handleSeleniumError(err, webDriver)
		return nil
	}

	link := ""
	if inward {
		link = "INWARD SERVICE OPTIONS"

	} else {
		link = "OUTWARD SERVICE OPTIONS"
	}

	webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByPartialLinkText, link)
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	})

	elem, err = webDriver.FindElement(selenium.ByPartialLinkText, link)
	if err != nil {
		handleSeleniumError(err, webDriver)
		return nil
	}

	err = elem.Click()
	if err != nil {
		handleSeleniumError(err, webDriver)
		return nil
	}

	err = waitForWaitFor(webDriver)
	if err != nil {
		handleSeleniumError(err, webDriver)
		return nil
	}

	v := checkTable(webDriver, inward)
	elem, err = webDriver.FindElement(selenium.ByPartialLinkText, "2")
	if err != nil {
		handleSeleniumError(err, webDriver)
		return v
	}
	err = elem.Click()
	if err != nil {
		handleSeleniumError(err, webDriver)
		return v
	}

	err = waitForWaitFor(webDriver)
	if err != nil {
		handleSeleniumError(err, webDriver)
		return v
	}

	b := checkTable(webDriver, inward)
	for _, x := range b {
		v = append(v, x)
	}
	return v
}

func checkTable(webDriver selenium.WebDriver, inward bool) []inwardService {

	table := ""

	if inward {
		table = "TABLEINWARDSERVICES"

	} else {
		table = "TABLEOUTWARDSERVICES"
	}

	var v []inwardService
	webDriver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByXPATH, "//table[@id='"+table+"']/tbody/tr[1]/td[13]")
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	})
	i := 1
	for i < 50 {
		var service string

		service, err := getTableElement(i, 2, table, webDriver)
		if err != nil {
			return v
		}
		subService, err := getTableElement(i, 3, table, webDriver)
		if err != nil {
			return v
		}
		destinationCode, err := getTableElement(i, 4, table, webDriver)
		if err != nil {
			return v
		}
		status, err := getTableElement(i, 13, table, webDriver)
		if err != nil {
			return v
		}
		v = append(v, inwardService{service: service, destination: destinationCode, subservice: subService, status: status})
		i++
	}
	return v
}

func getTableElement(row, column int, table string, webDriver selenium.WebDriver) (string, error) {
	elem, err := webDriver.FindElement(selenium.ByXPATH, fmt.Sprintf("//table[@id='"+table+"']/tbody/tr[%v]/td[%v]", row, column))
	if err != nil {
		return "", err
	}
	return elem.Text()
}
