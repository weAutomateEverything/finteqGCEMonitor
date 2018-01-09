package service

/*
Ping will send a sample message to HAL and return the response
 */
func Ping() error {
	err := sendError("Ping... Pong",nil,false)
	if err != nil {
		return err
	}
	err = getSessions()
	if err != nil {
		return err
	}
	return nil
}
