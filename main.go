package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
)

var blue = color.New(color.FgBlue).SprintFunc()
var bgblue = color.New(color.BgBlue).SprintfFunc()

// Context holds additional info in a log payload about the log event
type Context map[string]interface{}

func (c Context) String() string {
	pairs := make([]string, len(c))
	i := 0
	for k, v := range c {
		pairs[i] = fmt.Sprintf("\t%s %s = %v", bgblue(" "), k, v)
		i++
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "\n")
}

// Payload is a log event
type Payload struct {
	Format  string  `json:"__f"`
	Message string  `json:"msg"`
	Time    int64   `json:"time"`
	Level   int64   `json:"level"`
	Logger  string  `json:"logger"`
	Context Context `json:"ctx"`
}

func (p *Payload) String() string {
	sec := p.Time / 1000
	nsec := p.Time % 1000 * 1e6
	ts := time.Unix(sec, nsec).Format(time.RFC3339Nano)
	return fmt.Sprintf("[%s] %s -> %s\n%s", ts, blue(p.Logger), p.Message, p.Context)
}

func parseLine(line []byte) {
	p := Payload{}
	err := json.Unmarshal(line, &p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[log!] error parsing log line: %s", err)
		return
	}
	fmt.Println(&p)
}

func main() {
	var flagNoColor = flag.Bool("no-color", false, "Disable color output")
	if *flagNoColor {
		color.NoColor = true
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if (stat.Mode() & os.ModeNamedPipe) == 0 {
		panic(fmt.Errorf("nothing to read on stdin"))
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parseLine(scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("reading standard input: %s", err))
	}
}
