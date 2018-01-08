package service

/*
Ping will send a sample message to HAL and return the response
 */
func Ping() error {
	return sendError("Ping... Pong",nil,true)
}
