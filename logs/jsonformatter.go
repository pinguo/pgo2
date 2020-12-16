package logs

import (
	"encoding/json"
	"strings"

	"github.com/pinguo/pgo2/iface"
)

var (
	// JSONFormatter log JSON formatter
	JSONFormatter IFormatter = jsonFormatter{}

	// JSONAccessFormatter access log with JSON formatter
	JSONAccessFormatter iface.IAccessLogFormat = jsonAccessFormatter{}
)

type (
	message string

	logItem struct {
		When    string  `json:"time"`
		ID      string  `json:"id"`
		Level   string  `json:"level"`
		Trace   string  `json:"trace,omitempty"`
		Message message `json:"message"`
	}

	accessItem struct {
		Method   string `json:"method"`
		Path     string `json:"path"`
		Status   int    `json:"status"`
		Size     int    `json:"size"`
		Elapse   int    `json:"elapse"`
		Push     string `json:"push"`
		Profile  string `json:"profile"`
		Counting string `json:"counting"`
	}

	jsonFormatter struct{}

	jsonAccessFormatter struct{}
)

func (j jsonFormatter) Format(item *LogItem) string {
	w := &strings.Builder{}
	json.NewEncoder(w).Encode(createLogItem(item))
	return w.String()
}

func (ja jsonAccessFormatter) Format(ctx iface.IContext) string {
	item := accessItem{
		Method:   ctx.Method(),
		Path:     ctx.Path(),
		Status:   ctx.Status(),
		Size:     ctx.Size(),
		Elapse:   ctx.ElapseMs(),
		Push:     ctx.PushLogString(),
		Profile:  ctx.ProfileString(),
		Counting: ctx.CountingString(),
	}
	w := &strings.Builder{}
	json.NewEncoder(w).Encode(item)
	return w.String()
}

func (m message) MarshalJSON() ([]byte, error) {
	b := []byte(m)
	if json.Valid(b) {
		return b, nil
	}
	return json.Marshal(string(m))
}

func createLogItem(i *LogItem) *logItem {
	return &logItem{
		When:    i.When.Format("2006/01/02 15:04:05.000"),
		ID:      i.LogId,
		Level:   LevelToString(i.Level),
		Trace:   i.Trace,
		Message: message(i.Message),
	}
}
