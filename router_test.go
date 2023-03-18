package web

import (
	"bytes"
	"fmt"
)

type testLogger struct {
	out bytes.Buffer
}

func (tl *testLogger) Debug(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Info(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Warn(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Error(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
func (tl *testLogger) Fatal(data ...interface{}) {
	tl.out.Write([]byte(fmt.Sprint(data...)))
}
