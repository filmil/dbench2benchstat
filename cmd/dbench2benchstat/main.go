package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"text/template"
	"time"
)

// msgLine is the regexp of the line to match on the input.
// Example:
// basic/increment -> avg 0.00927782162588819ms out of 2534 samples. (std dev 0.00161196365508829, min 0.008, max 0.082)
var msgLine = regexp.MustCompile(`\s*(\S+)\s+-> avg (\S+) out of (\d+) samples. \(std dev\s+(\S+), min (\S+), max (\S+)\)`)

var funcs = template.FuncMap{
	"ns": asNanoseconds,
}

// outLine is the template for an output line that is emitted.
var outLine = template.Must(
	template.New("dbench").
		Funcs(funcs).
		Parse("Benchmark{{.Name}}\t{{.NumSamples}}\t{{ns .Average}} ns/op"))

const (
	// NameGroup is the group matching the test name.
	NameGroup = 1
	// AvgGroup is the group matching the average test time.
	AvgGroup = 2
	// NumSamplesGroup is the group matching the number of samples collected.
	NumSamplesGroup = 3
	// StdDevGroup is the group matching the standard deviation of the samples.
	StdDevGroup = 4
	// MinGroup is the group matching the minimum time.
	MinGroup = 5
	// MaxGroup is the group matching the maximum time.
	MaxGroup = 6
)

// Result is the parsed result out of a message line.
type Result struct {
	Name       string
	Average    time.Duration
	NumSamples int
	StdDev     time.Duration
	Min        time.Duration
	Max        time.Duration
}

func asNanoseconds(d time.Duration) int64 {
	return d.Nanoseconds()
}

// parseFloat parses out a duration from a string that has no units,
// like: "0.12345".  The default units in this case are milliseconds.
func parseFloat(s string) (time.Duration, error) {
	var (
		raw float64
		d   time.Duration
	)
	raw, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return d, fmt.Errorf("could not parse as float64: %v: %v", s, err)
	}
	rawNs := int(float64(time.Millisecond/time.Nanosecond) * raw)
	d, err = time.ParseDuration(fmt.Sprintf("%vns", rawNs))
	if err != nil {
		return d, fmt.Errorf("could not parse duration from: %d", rawNs)
	}
	return d, nil
}

// Get result parses a Result out of a single line of text.
func GetResult(t string) (Result, error) {
	var r Result
	m := msgLine.FindStringSubmatch(t)
	if len(m) != MaxGroup+1 {
		// If a line does not parse, skipBuffer it.
		return r, fmt.Errorf("could not parse: %q: %+v", t, m)
	}

	r.Name = m[NameGroup]
	var err error
	r.Average, err = time.ParseDuration(m[AvgGroup])
	if err != nil {
		return r, fmt.Errorf("could not parse average duration from: %q: %v", t, err)
	}
	rawSamples, err := strconv.ParseInt(m[NumSamplesGroup], 10, 32)
	if err != nil {
		return r, fmt.Errorf("could not parse number of samples from: %q: %v", t, err)
	}
	r.NumSamples = int(rawSamples)
	r.StdDev, err = parseFloat(m[StdDevGroup])
	if err != nil {
		return r, fmt.Errorf("could not parse std dev from: %q: %v", t, err)
	}
	r.Min, err = parseFloat(m[MinGroup])
	if err != nil {
		return r, fmt.Errorf("could not parse min duration from: %q: %v", t, err)
	}
	r.Max, err = parseFloat(m[MaxGroup])
	if err != nil {
		return r, fmt.Errorf("could not parse max duration from: %q: %v", t, err)
	}
	return r, nil
}

func Process(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)

	for s.Scan() {
		t := s.Text()
		res, err := GetResult(t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse: %q: %v\n", t, err)
			continue
		}
		if err = outLine.Execute(w, res); err != nil {
			return fmt.Errorf("could not output a line for: %q: %v", t, err)
		}
		fmt.Fprintln(w, "")
	}

	return nil
}

func main() {
	if err := Process(os.Stdin, os.Stdout); err != nil {
		fmt.Errorf("Could not process input: %v", err)
	}
}
