package sllmbench

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"git.fractalqb.de/fractalqb/sllm/v3"
	jsoniter "github.com/json-iterator/go"
)

// Real-world example:
// rsyslog.service: Sent signal SIGHUP to main process 1611 (rsyslogd) on client request.

func ExampleWrite() {
	sllm.Fprint(os.Stdout, sllmForm, sllmArgs)
	fmt.Println()
	json.NewEncoder(os.Stdout).Encode(jsonDynamic)
	json.NewEncoder(os.Stdout).Encode(jsonStatic)
	// Output:
	// `service:rsyslog`: Sent `signal:SIGHUP` to main `process:1611` (`name:rsyslogd`) on client request.
	// {"msg":"Sent signal to main process on client request.","process":1611,"process-name":"rsyslogd","service":"rsyslog","signal":"SIGHUP"}
	// {"Msg":"Sent signal to main process on client request.","Service":"rsyslog","Signal":"SIGHUP","process-name":"rsyslogd","Process":1611}
}

func ExampleParse() {
	var tmpl bytes.Buffer
	args, _ := sllm.ParseMap(
		"`service:rsyslog`: Sent `signal:SIGHUP` to main `process:1611` (`name:rsyslogd`) on client request.",
		&tmpl,
	)
	var jsonData map[string]any
	json.Unmarshal([]byte(jsonMsg), &jsonData)
	fmt.Println(tmpl.String(), args)
	fmt.Println(jsonData)
	// Output:
	// `service`: Sent `signal` to main `process` (`name`) on client request. map[name:[rsyslogd] process:[1611] service:[rsyslog] signal:[SIGHUP]]
	// map[msg:Sent signal to main process on client request. process:1611 process-name:rsyslogd service:rsyslog signal:SIGHUP]
}

var (
	sllmForm = "`service`: Sent `signal` to main `process` (`name`) on client request."
	sllmArgs = sllm.IdxArgsDefault("???", "rsyslog", "SIGHUP", 1611, "rsyslogd")
)

var result int

func BenchmarkSllmAppend(b *testing.B) {
	var buf []byte
	for i := 0; i < b.N; i++ {
		buf, _ = sllm.Append(buf[:0], sllmForm, sllmArgs)
	}
	result = len(buf)
}

func BenchmarkSllmByteBuffer(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		buf.Reset()
		sllm.Fprint(&buf, sllmForm, sllmArgs)
	}
	result = buf.Len()
}

func BenchmarkSllmStringBuilder(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < b.N; i++ {
		sb.Reset()
		sllm.Fprint(&sb, sllmForm, sllmArgs)
	}
	result = sb.Len()
}

type staticMsg struct {
	Msg, Service, Signal string
	ProcessName          string `json:"process-name"`
	Process              int
}

var (
	jsonDynamic = map[string]any{
		"msg":          "Sent signal to main process on client request.",
		"service":      "rsyslog",
		"signal":       "SIGHUP",
		"process":      1611,
		"process-name": "rsyslogd",
	}
	jsonStatic = staticMsg{
		Msg:         "Sent signal to main process on client request.",
		Service:     "rsyslog",
		Signal:      "SIGHUP",
		Process:     1611,
		ProcessName: "rsyslogd",
	}
)

const jsonMsg = `{"msg":"Sent signal to main process on client request.","process":1611,"process-name":"rsyslogd","service":"rsyslog","signal":"SIGHUP"}`

func BenchmarkGoJSONDynamic(b *testing.B) {
	var buf []byte
	for i := 0; i < b.N; i++ {
		buf, _ = json.Marshal(jsonDynamic)
	}
	result = len(buf)
}

func BenchmarkGoJSONStatic(b *testing.B) {
	var buf []byte
	for i := 0; i < b.N; i++ {
		buf, _ = json.Marshal(jsonStatic)
	}
	result = len(buf)
}

func BenchmarkJSONiterDynamic(b *testing.B) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	var buf []byte
	for i := 0; i < b.N; i++ {
		buf, _ = json.Marshal(jsonDynamic)
	}
	result = len(buf)
}

func BenchmarkJSONiterStatic(b *testing.B) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	var buf []byte
	for i := 0; i < b.N; i++ {
		buf, _ = json.Marshal(jsonStatic)
	}
	result = len(buf)
}

func BenchmarkSllmParseDynamic(b *testing.B) {
	var (
		tmpl bytes.Buffer
		m    map[string][]any
	)
	for i := 0; i < b.N; i++ {
		m, _ = sllm.ParseMap(
			"`service:rsyslog`: Sent `signal:SIGHUP` to main `process:1611` (`name:rsyslogd`) on client request.",
			&tmpl,
		)
	}
	result = len(m)
}

func BenchmarkGoJSONparseDynamic(b *testing.B) {
	msg := []byte(jsonMsg)
	var data map[string]any
	for i := 0; i < b.N; i++ {
		json.Unmarshal(msg, &data)
	}
	result = len(data)
}

func BenchmarkGoJSONparseStatic(b *testing.B) {
	msg := []byte(jsonMsg)
	var data staticMsg
	for i := 0; i < b.N; i++ {
		json.Unmarshal(msg, &data)
	}
	result = len(data.Msg)
}

func BenchmarkGoJSONiterParseDynamic(b *testing.B) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	msg := []byte(jsonMsg)
	var data map[string]any
	for i := 0; i < b.N; i++ {
		json.Unmarshal(msg, &data)
	}
	result = len(data)
}

func BenchmarkGoJSONiterParseStatic(b *testing.B) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b.ResetTimer()
	msg := []byte(jsonMsg)
	var data staticMsg
	for i := 0; i < b.N; i++ {
		json.Unmarshal(msg, &data)
	}
	result = len(data.Msg)
}
