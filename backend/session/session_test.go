package session_test

import (
	"net/http"
	"testing"

	"github.com/smockyio/smocky/backend/session"
	"github.com/stretchr/testify/assert"
)

func TestSession_GetSet(t *testing.T) {
	sess := session.New()

	sess.Set("post", "/here", "age", 20)
	assert.Equal(t, 20, sess.Get("post", "/here", "age"))
	assert.Equal(t, 20, sess.GetInt("post", "/here", "age"))

	sess.Increase("post", "/here", "age")
	assert.Equal(t, 21, sess.GetInt("post", "/here", "age"))

	req, _ := http.NewRequest("GET", "https://here.com/there", nil)
	assert.Equal(t, 1, sess.IncreaseRequestNumber(req))
	assert.Equal(t, 1, sess.GetRequestNumber(req))

	sess.SetNextResponseIndex(req, 5)
	assert.Equal(t, 5, sess.NextResponseIndex(req))
}
