package ghoko

import (
	"fmt"
	"net/http"
)

var (
	ErrForbidden        = hokoErr{http.StatusForbidden, "Access Deny"}
	ErrMethodNotAllowed = hokoErr{http.StatusMethodNotAllowed, "Method Not Allowed"}
	ErrSyncNeeded       = fmt.Errorf("`sync` param needed")
)

type hokoErr struct {
	status int
	err    string
}

func (err hokoErr) Error() string {
	return err.err
}
