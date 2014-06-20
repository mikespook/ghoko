package ghoko

import "net/http"

var (
	ErrSyncNeeded = &HttpError{http.StatusBadRequest, "`Ghoko-sync` header needed"}
	ErrForbidden  = &HttpError{http.StatusForbidden, "Incorrect `_secret` parameter"}
	ErrNotFound   = &HttpError{http.StatusNotFound, "Request path was not found"}
)

type HttpError struct {
	status  int
	message string
}

func (err *HttpError) Error() string {
	return err.message
}
