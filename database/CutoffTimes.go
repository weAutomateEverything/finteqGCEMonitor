package database

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

/*
CutoffTime database object for start of day and end of day
 */
type CutoffTime struct {
	Service    string
	SubService string
	SodHour    int
	SodMinute  int
	EodHour    int
	EodMinute  int
	DayOfWeek int
}

/*
SaveCutoff adds a new cutoff time to the DB
 */
func SaveCutoff(time CutoffTime){
	c := database.C("CutoffTimes")
	c.Insert(time)
}

/*
CutoffExists returns true if records exist for the service
 */
func CutoffExists(service, subservice string) bool{
	c := database.C("CutoffTimes")
	var r []CutoffTime
	c.Find(bson.M{"service":service,"subservice":subservice,"dayofweek":time.Now().Weekday()}).All(&r)
	return len(r) > 0
}

/*
IsInStartOfDay returns true if a record is found where the current time > start of day but < end of day
 */
func IsInStartOfDay(service, subservice string) bool {
	c := database.C("CutoffTimes")
	var r []CutoffTime
	t := time.Now()
	// Look for all records there the start hour is less than or equal to the current hour, and the eod hour is greather than or equal to the current hour.
	c.Find(bson.M{"service":service,"subservice":subservice,"sodhour":bson.M{"$lte":t.Hour()},"eodhour": bson.M{"$gte":t.Hour()},"dayofweek":t.Weekday() }).All(&r)

	for _,c := range r {
		//For example, its 14:05, but the SOD is only 14:30 - we need to check the minute
		if c.SodHour == t.Hour() {
			//Now, we need to check that the minute is greater than or equal to SOD min
			return c.SodMinute >= t.Minute()
		}

		//We need to do the same for the EOD, if its 15:01 and the cutoff time is 15:00 thjen the window is closed
		if c.EodHour == t.Hour(){
			return c.EodMinute <= t.Minute()
		}

		//log.Printf("Found Start of day. Service: %v, subservice: %v, start: %02d:%02d, end: %02d:%02d, day: %v  ",c.Service,c.SubService,c.SodHour, c.SodMinute,c.EodHour,c.EodMinute,c.DayOfWeek)

		return true
	}

	return false
}



