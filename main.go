package main

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DataPoint present a single data point
type DataPoint struct {
	Name   string
	Value  float64
	Time   time.Time
	Labels map[string]string
}

func main() {
	http.HandleFunc("/", metricServer)
	http.ListenAndServe(":8080", nil)
}

func metricServer(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	s := buf.String()
	arr, err := parsePrometheus(s)
	if err != nil {
		w.WriteHeader(404)
	} else {
		c := countInf(arr)
		fmt.Fprintf(w, "Hello, %d data points are received and %d of them are Inf!\n", len(arr), c)
	}
}

func countInf(dps []DataPoint) int {
	counter := 0
	for _, dp := range dps {
		if math.Inf(1) == dp.Value || math.Inf(-1) == dp.Value {
			counter++
		}
	}
	return counter
}

// Prometheus text format can be found here:
//   https://github.com/prometheus/docs/blob/master/content/docs/instrumenting/exposition_formats.md#text-format-details
func parsePrometheus(str string) ([]DataPoint, error) {
	lines := strings.Split(str, "\n")
	result := make([]DataPoint, 0)

	for _, line := range lines {
		dp, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		if dp == nil {
			continue
		}
		result = append(result, *dp)
	}

	return result, nil
}

func parseLine(line string) (*DataPoint, error) {
	if len(line) == 0 || strings.HasPrefix(line, "#") {
		return nil, nil
	}

	fields := strings.Fields(line)
	if len(fields) < 2 || len(fields) > 3 {
		return nil, fmt.Errorf("failed to parse line: %s", line)
	}

	t := time.Now()
	if len(fields) == 3 {
		tInt, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time stamp %s", fields[2])
		}
		t = time.Unix(0, tInt*int64(time.Millisecond))
	}

	v, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}

	index := strings.IndexAny(fields[0], "{")
	name := fields[0][:index]
	labels := parseLabels(fields[0][index:])
	return &DataPoint{
		Name:   name,
		Time:   t,
		Value:  v,
		Labels: labels,
	}, nil
}

func parseLabels(field string) map[string]string {
	field = strings.Trim(field, "{} ")
	tags := strings.Split(field, ",")
	result := make(map[string]string)
	regex := regexp.MustCompile(`(.*)="(.*)"`)
	for _, tag := range tags {
		matches := regex.FindStringSubmatch(tag)
		result[matches[1]] = matches[2]
	}
	return result
}
