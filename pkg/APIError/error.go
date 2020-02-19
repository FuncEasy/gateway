package APIError

import "net/http"

type Error struct {
	Err  error
	Msg  string
	Code int
}

func PanicError(err error, msg string, code int) {
	panic(&Error{
		Err:  err,
		Msg:  msg,
		Code: code,
	})
}

func Panic(err error) {
	if err != nil {
		PanicError(err, "Server Error", http.StatusInternalServerError)
	}
}
