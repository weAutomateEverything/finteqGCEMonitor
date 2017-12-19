package service

import (
	"github.com/tebeka/selenium"
	"strings"
	"log"
	"errors"
)

func checkServices(driver selenium.WebDriver) {
	elem, err := driver.FindElement(selenium.ByPartialLinkText, "Monitor Services")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	elem.MoveTo(10, 10)

	elem, err = driver.FindElement(selenium.ByPartialLinkText, "Monitor Outward Services")
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

	err = driver.Wait(func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByPartialLinkText, "OUTWARD")
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	})
	log.Println("done with waiting for OUTWARD")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	log.Println("Getting first level of elements")
	elems, err := driver.FindElements(selenium.ByClassName, "ErrorTreeView")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	log.Printf("Found %v first level elements", len(elems))

	//Expand the first level of trees
	for _, elem = range elems {
		id, err := elem.GetAttribute("id")
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}
		if strings.Index(id, "|") > 0 {
			continue
		}
		id = strings.Replace(id, "TD", "P", 1)

		log.Printf("searching for element %v", id)
		elem, err = driver.FindElement(selenium.ByID, id)
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}
		log.Printf("Clicking %v", id)

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

	log.Println("Moving to second level")

	//Expand the Second Level of Trees
	elems, err = driver.FindElements(selenium.ByClassName, "ErrorTreeView")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	log.Printf("found %v second level items",len(elems))

	for _, elem = range elems {
		id, err := elem.GetAttribute("id")
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}

		if strings.Index(id, "|") > 0 {
			id = strings.Replace(id, "TD", "P", 1)
			log.Printf("using id %v",id)

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
				log.Printf("%v is not displayed",id)
			}
		}

	}

	log.Println("Moving to Tree")

	//Click on each of the flags and get the data

	elems, err = driver.FindElements(selenium.ByCSSSelector, "TD.ErrorTreeView > A.TreeviewText")
	if err != nil {
		handleSeleniumError(err, driver)
		return
	}

	log.Printf("%v flags found",len(elems))
	handleSeleniumError(errors.New("test"),driver)

	for _, elem = range elems {
		err = elem.Click()
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}

		err := waitForWaitFor(driver)
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}

		grids, err := driver.FindElements(selenium.ByClassName, "GridCellError")
		if err != nil {
			handleSeleniumError(err, driver)
			return
		}
		log.Println(elem.Text())

		for _, grid := range grids {
			log.Println(grid.Text())
		}
	}

	log.Println("All Done")
}
