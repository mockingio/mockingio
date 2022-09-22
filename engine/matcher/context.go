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
	return fmt.Sprintf("%s/count", r.Key())
}

func (r Context) SequenceID() string {
	return fmt.Sprintf("%s/sequence", r.Key())
}

func (r Context) Key() string {
	return fmt.Sprintf("%s/%s-%s", r.SessionID, r.HTTPRequest.Method, r.HTTPRequest.URL)
}
