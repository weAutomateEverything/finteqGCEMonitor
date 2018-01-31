package gceservices

import (
	"github.com/tebeka/selenium"
	"strings"
	"log"
	"fmt"
	selenium2 "github.com/CardFrontendDevopsTeam/FinteqGCEMonitor/selenium"
	"github.com/pkg/errors"
)

type Service interface{
	RunServiceCheck()
}

type service struct {
	selenium selenium2.Service
	inward bool
}

func NewService(svc selenium2.Service, inward bool) Service{
	return &service{svc, inward}
}

type gceError struct {
	uid            string
	filename       string
	status         string
	date           string
	runno          string
	runid          string
	attempts       string
	maxattempts    string
	responseCode   string
	nonStdFileName string
	description    string
}

func (s service)RunServiceCheck(){
	err := s.checkServices()
	if err != nil {
		s.selenium.HandleSeleniumError(err)
	}
}

func (s *service)checkServices() error {
	driver := s.selenium.Driver()
	elem, err := driver.FindElement(selenium.ByPartialLinkText, "Monitor Services")
	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	elem.MoveTo(10, 10)

	link := ""
	if s.inward {
		link = "Monitor Inward Services"
	} else {
		link = "Monitor Outward Services"
	}

	elem, err = driver.FindElement(selenium.ByPartialLinkText, link)
	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	err = elem.Click()

	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	err = s.selenium.WaitForWaitFor()
	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	if s.inward {
		link = "INWARD"
	} else {
		link = "OUTWARD"
	}
	err = driver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByPartialLinkText, link)
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	})
	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	elems, err := driver.FindElements(selenium.ByClassName, "ErrorTreeView")
	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	//Expand the first level of trees
	for _, elem = range elems {
		id, err := elem.GetAttribute("id")
		if err != nil {
			return &selenium2.SeleniumnError{true,err}
		}
		if strings.Index(id, "^") > 0 {
			continue
		}
		id = strings.Replace(id, "TD", "P", 1)

		elem, err = driver.FindElement(selenium.ByID, id)
		if err != nil {
			return &selenium2.SeleniumnError{true,err}
		}

		displayed, err := elem.IsDisplayed()
		if err != nil {
			return &selenium2.SeleniumnError{true,err}
		}

		if displayed {
			err = elem.Click()
			if err != nil {
				return &selenium2.SeleniumnError{true,err}
			}
		}
	}

	//Expand the Second Level of Trees
	elems, err = driver.FindElements(selenium.ByClassName, "ErrorTreeView")
	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	for _, elem = range elems {
		id, err := elem.GetAttribute("id")
		if err != nil {
			return &selenium2.SeleniumnError{true,err}
		}

		if strings.Index(id, "^") > 0 {
			id = strings.Replace(id, "TD", "P", 1)
			log.Printf("using id %v", id)

			elem, err = driver.FindElement(selenium.ByID, id)
			if err != nil {
				return &selenium2.SeleniumnError{true,err}
			}
			displayed, err := elem.IsDisplayed()
			if err != nil {
				return &selenium2.SeleniumnError{true,err}
			}
			if displayed {
				err = elem.Click()
				if err != nil {
					return &selenium2.SeleniumnError{true,err}
				}
			} else {
				log.Printf("%v is not displayed", id)
			}
		}

	}

	//Click on each of the flags and get the data

	elems, err = driver.FindElements(selenium.ByCSSSelector, "TD.ErrorTreeView > A.TreeviewText")
	if err != nil {
		return &selenium2.SeleniumnError{true,err}
	}

	for _, elem = range elems {

		displayed, err := elem.IsDisplayed()
		if err != nil {
			s.selenium.HandleSeleniumError(&selenium2.SeleniumnError{true,err})
			continue
		}

		if !displayed {
			continue
		}

		err = elem.Click()
		if err != nil {
			s.selenium.HandleSeleniumError(&selenium2.SeleniumnError{true,err})
			continue
		}

		err = s.selenium.WaitForWaitFor()
		if err != nil {
			s.selenium.HandleSeleniumError(&selenium2.SeleniumnError{true,err})
			continue
		}

		grids, err := driver.FindElements(selenium.ByClassName, "GridCellError")
		if err != nil {
			s.selenium.HandleSeleniumError(&selenium2.SeleniumnError{true,err})
			continue
		}
		log.Println(elem.Text())

		if len(grids) > 11 {
			//We are only going to report the top row.
			error := gceError{
				uid:            getText(grids[1]),
				filename:       getText(grids[2]),
				status:         getText(grids[3]),
				date:           getText(grids[4]),
				runno:          getText(grids[5]),
				runid:          getText(grids[6]),
				attempts:       getText(grids[7]),
				maxattempts:    getText(grids[8]),
				responseCode:   getText(grids[9]),
				nonStdFileName: getText(grids[10]),
				description:    getText(grids[11]),
			}

			msg := fmt.Sprintf("*GCE Error*\n UID: %v\n File Name: %v\n Status: %v\n Timestamp: %v\n Run Number: %v\n  "+
				"Run Indicator: %v\n Attempts: %v\n Max Attempts: %v\n Response Code: %v\n Non Standdard file name: %v\n Description: %v", error.uid,
				error.filename, error.status, error.date, error.runno, error.runid, error.attempts, error.maxattempts,
				error.responseCode, error.nonStdFileName, error.description)

			s.selenium.HandleSeleniumError(&selenium2.SeleniumnError{false,errors.New(msg)})

		}
	}
	return nil
}


func getText(element selenium.WebElement) string {
	str, _ := element.Text()
	return str
}



