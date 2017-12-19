package service

import (
	"github.com/tebeka/selenium"
	"strings"
	"log"
	"fmt"
)

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

func checkServices(driver selenium.WebDriver, inward bool) {
	elem, err := driver.FindElement(selenium.ByPartialLinkText, "Monitor Services")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	elem.MoveTo(10, 10)

	link := ""
	if inward{
		link = "Monitor Inward Services"
	} else {
		link = "Monitor Outward Services"
	}

	elem, err = driver.FindElement(selenium.ByPartialLinkText, link)
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	err = elem.Click()

	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	err = waitForWaitFor(driver)
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	if inward {
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
		handleSeleniumError(err, driver)
		return
	}

	elems, err := driver.FindElements(selenium.ByClassName, "ErrorTreeView")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	//Expand the first level of trees
	for _, elem = range elems {
		id, err := elem.GetAttribute("id")
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}
		if strings.Index(id, "^") > 0 {
			continue
		}
		id = strings.Replace(id, "TD", "P", 1)

		elem, err = driver.FindElement(selenium.ByID, id)
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}

		displayed, err := elem.IsDisplayed()
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}

		if displayed {
			err = elem.Click()
			if err != nil {
				handleSeleniumError(err, driver)
				return
			}
		}
	}

	//Expand the Second Level of Trees
	elems, err = driver.FindElements(selenium.ByClassName, "ErrorTreeView")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	for _, elem = range elems {
		id, err := elem.GetAttribute("id")
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}

		if strings.Index(id, "^") > 0 {
			id = strings.Replace(id, "TD", "P", 1)
			log.Printf("using id %v", id)

			elem, err = driver.FindElement(selenium.ByID, id)
			if err != nil {
				handleSeleniumError(err, driver)
				return
			}
			displayed, err := elem.IsDisplayed()
			if err != nil {
				handleSeleniumError(err, driver)
				return
			}
			if displayed {
				err = elem.Click()
				if err != nil {
					handleSeleniumError(err, driver)
					return
				}
			} else {
				log.Printf("%v is not displayed", id)
			}
		}

	}

	//Click on each of the flags and get the data

	elems, err = driver.FindElements(selenium.ByCSSSelector, "TD.ErrorTreeView > A.TreeviewText")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	for _, elem = range elems {

		displayed, err := elem.IsDisplayed()
		if err != nil {
			handleSeleniumError(err, driver)
			continue
		}

		if !displayed{
			continue
		}


		err = elem.Click()
		if err != nil {
			handleSeleniumError(err, driver)
			continue
		}

		err = waitForWaitFor(driver)
		if err != nil {
			handleSeleniumError(err, driver)
			continue
		}

		grids, err := driver.FindElements(selenium.ByClassName, "GridCellError")
		if err != nil {
			handleSeleniumError(err, driver)
			continue
		}
		log.Println(elem.Text())

		// We have 12 items per row
		rows := len(grids) / 12
		row := 0
		for row < rows {
			offset := 12 * row
			error := gceError{
				uid:            getText(grids[1+offset]),
				filename:       getText(grids[2+offset]),
				status:         getText(grids[3+offset]),
				date:           getText(grids[4+offset]),
				runno:          getText(grids[5+offset]),
				runid:          getText(grids[6+offset]),
				attempts:       getText(grids[7+offset]),
				maxattempts:    getText(grids[8+offset]),
				responseCode:   getText(grids[9+offset]),
				nonStdFileName: getText(grids[10+offset]),
				description:    getText(grids[11+offset]),
			}

			msg := fmt.Sprintf("*GCE Error*\n UID: %v\n File Name: %v\n Status: %v\n Timestamp: %v\n Run Number: %v\n  " +
				"Run Indicator: %v\n Attempts: %v\n Max Attempts: %v\n Response Code: %v\n Non Standdard file name: %v\n Description: %v",error.uid,
					error.filename,error.status,error.date,error.runno,error.runid,error.attempts,error.maxattempts,
						error.responseCode,error.nonStdFileName,error.description)

			img, _ := driver.Screenshot()
			sendError(msg,img,false)
			row++
		}
	}
}

func getText(element selenium.WebElement) string {
	str, _ := element.Text()
	return str
}
