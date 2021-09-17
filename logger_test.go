package rome

import (
	"bytes"
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/amermelao/rome/data"
)

type logging struct{ goLog *log.Logger }

func (l *logging) Panic(v ...interface{}) {
	l.goLog.Print(prepearMsg("panic: ", v)...)
}
func (l *logging) Fatal(v ...interface{}) {
	l.goLog.Print(prepearMsg("fatal: ", v)...)
}
func (l *logging) Error(v ...interface{}) {
	l.goLog.Print(prepearMsg("error: ", v)...)
}
func (l *logging) Warning(v ...interface{}) {
	l.goLog.Print(prepearMsg("warning: ", v)...)
}
func (l *logging) Info(v ...interface{}) {
	l.goLog.Print(prepearMsg("info: ", v)...)
}
func (l *logging) Debug(v ...interface{}) {
	l.goLog.Print(prepearMsg("debug: ", v)...)
}
func (l *logging) Trace(v ...interface{}) {
	l.goLog.Print(prepearMsg("trace: ", v)...)
}

func prepearMsg(lvl string, v []interface{}) []interface{} {
	values := []interface{}{}
	values = append(values, lvl)
	values = append(values, v...)
	return values
}

func TestConcurrentLogger(t *testing.T) {
	var b bytes.Buffer
	standartLogger := log.New(&b, "", 0)
	centralLogger := NewCentrCentralLogger(&logging{standartLogger})
	messages := []data.Messages{
		{
			{Level: "panic", Content: "1"},
			{Level: "fatal", Content: "2"},
			{Level: "error", Content: "3"},
		},
		{
			{Level: "warning", Content: "a1"},
			{Level: "info", Content: "a2"},
			{Level: "debug", Content: "a3"},
		},
		{
			{Level: "trace", Content: "b1"},
			{Level: "panic", Content: "b2"},
			{Level: "info", Content: "b3"},
		},
	}
	var wg sync.WaitGroup
	for _, v := range messages {
		v := v
		wg.Add(1)
		go func() {
			centralLogger.Log(v)
			wg.Done()
		}()
	}
	wg.Wait()
	centralLogger.Close()

	values := b.String()
	splitted := strings.Split(values, "\n")
	expected := map[string][2]string{
		"panic: 1":    {"fatal: 2", "error: 3"},
		"warning: a1": {"info: a2", "debug: a3"},
		"trace: b1":   {"panic: b2", "info: b3"},
	}

	for cont, value := range splitted[0 : len(splitted)-1] {
		if cont%3 == 0 {
			tmpExpected := expected[value]
			if tmpExpected[0] != splitted[cont+1] &&
				tmpExpected[1] != splitted[cont+2] {
				t.Error("failed test")
			}
		}
	}
}
