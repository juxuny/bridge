package log

import (
	"encoding/json"
	"fmt"
	"time"
)

type logRecord struct {
	Time  time.Time   `json:"time"`
	Level string      `json:"level"`
	Data  interface{} `json:"data"`
}

type JsonLogger struct {
	w      *writer
	prefix string
}

func (t *JsonLogger) SetPrefix(s string) {
	DefaultLogger.Error("no implement")
}

func (t *JsonLogger) Println(v ...interface{}) {
	if len(v) > 0 {
		t.Output(t.newRecord("INFO", v[0]))
	}
}

func (t *JsonLogger) Error(v ...interface{}) {
	if len(v) > 0 {
		t.Output(t.newRecord("ERROR", v[0]))
	}
}

func (t *JsonLogger) Warn(v ...interface{}) {
	if len(v) > 0 {
		t.Output(t.newRecord("WARN", v[0]))
	}
}

func (t *JsonLogger) Info(v ...interface{}) {
	if len(v) > 0 {
		t.Output(t.newRecord("INFO", v[0]))
	}
}

func (t *JsonLogger) Debug(v ...interface{}) {
	if len(v) > 0 {
		t.Output(t.newRecord("DEBUG", v[0]))
	}
}

func (t *JsonLogger) Printf(format string, v ...interface{}) {
	t.Output(t.newRecord("INFO ", fmt.Sprintf(format, v...)))
}

func (t *JsonLogger) SQL(v ...interface{}) {
	t.Output(t.newRecord("SQL", v))
}

func (t *JsonLogger) Consuming(v ...interface{}) {
	t.Output(t.newRecord("CONSUMING", v))
}

func (t *JsonLogger) Flush() {
	if err := t.w.Flush(); err != nil {
		DefaultLogger.Error(err)
	}
}

func (t *JsonLogger) newRecord(level string, data interface{}) logRecord {
	return logRecord{
		Time:  time.Now(),
		Level: level,
		Data:  data,
	}
}

func (t *JsonLogger) Print(v ...interface{}) {
	if len(v) > 0 {
		t.Output(t.newRecord("INFO", v[0]))
	}
}

func (t *JsonLogger) Output(v ...interface{}) {
	if len(v) > 0 {
		data, _ := json.Marshal(v[0])
		data = append(data, '\n')
		_, _ = t.w.Write(data)
	}
}

func NewJsonLogger(dir string, filePrefix string) ILogger {
	ret := JsonLogger{
		w: NewWriter(
			dir,
			filePrefix,
		),
	}
	return &ret
}
