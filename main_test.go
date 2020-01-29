package main

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestParser_00(t *testing.T) {
	in := `
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total counter
http_requests_total{method="post",code="200"} Inf 1395066363000
http_requests_total{method="post",code="400"}    3 1395066363000`
	arr, err := parsePrometheus(in)
	ok(t, err)
	equals(t, 2, len(arr))
	equals(t, 1, countInf(arr))
}

func TestParser_01(t *testing.T) {
	in := `
# HELP http_requests_total The total number of HTTP requests.
http_requests_total{method="post",code="200"} Inf 1395066363000
http_requests_total{method="post",code="400"} -Inf 1395066363000`
	arr, err := parsePrometheus(in)
	ok(t, err)
	equals(t, 2, len(arr))
	equals(t, 2, countInf(arr))
}

func TestParser_02(t *testing.T) {
	in := `
# HELP http_requests_total The total number of HTTP requests.
http_requests_total{method="post",code="200"} Inf
http_requests_total{method="post",code="400"} 0.0 1395066363000`
	arr, err := parsePrometheus(in)
	ok(t, err)
	equals(t, 2, len(arr))
	equals(t, 1, countInf(arr))
}
