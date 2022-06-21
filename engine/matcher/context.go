package matcher

import (
	"fmt"
	"net/http"
)

type Context struct {
	HTTPRequest *http.Request
	SessionID   string
}

func (r Context) CountID() string {
	return fmt.Sprintf("%s-%s-%s-count", r.HTTPRequest.Method, r.HTTPRequest.URL, r.SessionID)
}

func (r Context) SequenceID() string {
	return fmt.Sprintf("%s-%s-%s-sequence", r.HTTPRequest.Method, r.HTTPRequest.URL, r.SessionID)
}
