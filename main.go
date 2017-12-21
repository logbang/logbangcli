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
var ctxPre = color.New(color.BgCyan).Sprint(" ")
var unparsedPre = color.New(color.BgYellow, color.FgBlack).Sprint("SKP")
var errorPre = color.New(color.BgRed, color.FgBlack).Sprint("ERR")
var warningPre = color.New(color.BgYellow, color.FgBlack).Sprint("WRN")
var infoPre = color.New(color.BgGreen, color.FgBlack).Sprint("INF")
var debugPre = color.New(color.BgBlue, color.FgBlack).Sprint("DBG")
var levelPre = []string{errorPre, warningPre, infoPre, debugPre}

// Context holds additional info in a log payload about the log event
type Context map[string]interface{}

func (c Context) String() string {
	pairs := make([]string, len(c))
	i := 0
	for k, v := range c {
		pairs[i] = fmt.Sprintf("\t%s %s = %v", ctxPre, k, v)
		i++
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "\n")
}

// Payload is a log event
type Payload struct {
	Format   string  `json:"__f"`
	Language string  `json:"__l"`
	Message  string  `json:"msg"`
	Time     int64   `json:"time"`
	Level    int64   `json:"level"`
	Logger   string  `json:"logger"`
	Context  Context `json:"ctx"`
}

func (p *Payload) String() string {
	sec := p.Time / 1000
	nsec := p.Time % 1000 * 1e6
	ts := time.Unix(sec, nsec).Format(time.RFC3339Nano)
	return strings.TrimSpace(
		fmt.Sprintf(
			"%s [%s] %s -> %s\n%s",
			levelPre[p.Level], ts, blue(p.Logger), p.Message, p.Context,
		),
	)
}

func parseLine(line []byte) {
	p := &Payload{}
	if err := json.Unmarshal(line, p); err == nil {
		fmt.Println(p)
	} else {
		fmt.Printf("%s %s\n", unparsedPre, line)
	}
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
