package route

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestOneBot11PostMux_ServeHTTP(t *testing.T) {
	mux := NewOneBot11PostMux()
	var data string

	mux.HandleGroup("echo :message", func(ctx *onebotContext) {
		data = ctx.params["message"]
	})

	var case1s = []groupMessage{
		{Message: `echo a`, MessageType: `group`, RawMessage: `a`},
		{Message: `echo some`, MessageType: `group`, RawMessage: `some`},
		{Message: `echo  echo`, MessageType: `group`, RawMessage: `some`},
	}

	for _, message := range case1s {
		msg, _ := json.Marshal(message)
		req := httptest.NewRequest("POST", "/", bytes.NewReader(msg))
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		assert.Equal(t, message.RawMessage, data)
	}
}
