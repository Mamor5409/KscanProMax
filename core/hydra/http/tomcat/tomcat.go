package tomcat

import "errors"

var (
	LoginFailedError  = errors.New("login failed")
	LoginTimeoutError = errors.New("login timeout")
)
