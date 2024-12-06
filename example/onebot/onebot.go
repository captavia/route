package onebot

import (
	"encoding/json"
	"github.com/captavia/route"
	"io"
	"net/http"
	"sync"
)

type OneBot11PostMux struct {
	router    *route.Router[*onebotContext]
	pool      *sync.Pool
	delimiter rune
}

func NewOneBot11PostMux() *OneBot11PostMux {
	var delimiter = ' '
	return &OneBot11PostMux{
		router: route.NewRouter[*onebotContext](route.WithDelimiter[*onebotContext](delimiter)),
		pool: &sync.Pool{
			New: func() interface{} {
				return new(onebotContext)
			},
		},
		delimiter: delimiter,
	}
}

func (c *OneBot11PostMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := c.pool.Get().(*onebotContext)
	defer c.pool.Put(ctx)

	bytes, readByteErr := io.ReadAll(r.Body)
	if readByteErr != nil {
		return
	}

	var msgType = new(messageType)
	if jsonErr := json.Unmarshal(bytes, msgType); jsonErr != nil {
		return
	}
	ctx.message = msgType.Message

	c.router.Serve(msgType.MessageType+string(c.delimiter)+msgType.Message, func() *onebotContext {
		return ctx
	})
}

func (c *OneBot11PostMux) HandleGroup(path string, handler route.Handler[*onebotContext]) {
	c.router.Handle("group"+string(c.delimiter)+path, handler)
}

type onebotContext struct {
	message string
	params  map[string]string
}

func (c *onebotContext) NotMatch() {

}

func (c *onebotContext) WithParam(params map[string]string) {
	c.params = params
}

type messageType struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	Message     string `json:"message"`
	RawMessage  string `json:"raw_message"`
}

type groupMessage struct {
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	SubType     string `json:"sub_type"`
	MessageId   int64  `json:"message_id"`
	GroupId     int64  `json:"group_id"`
	UserId      int64  `json:"user_id"`
	Anonymous   struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Flag int64  `json:"flag"`
	} `json:"anonymous"`
	Message    string `json:"message"`
	RawMessage string `json:"raw_message"`
	Sender     struct {
		UserId   int64  `json:"user_id"`
		NickName string `json:"nickname"`
		Card     string `json:"card"`
		Sex      string `json:"sex"`
		Age      int32  `json:"age"`
		Area     string `json:"area"`
		Level    string `json:"level"`
		Role     string `json:"role"`
		Title    string `json:"title"`
	}
}
