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
	return fmt.Sprintf("%s/%s-%s/count", r.SessionID, r.HTTPRequest.Method, r.HTTPRequest.URL)
}

func (r Context) SequenceID() string {
	return fmt.Sprintf("%s/%s-%s/sequence", r.SessionID, r.HTTPRequest.Method, r.HTTPRequest.URL)
}
