package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// Logline represents one line from a log file for easy sorting by timestamp
type Logline struct {
	Timestamp time.Time
	Body      string
}

// Logfile represents a log file, its current state, and the most recently read line
// from that logfile
type Logfile struct {
	Name   string
	Handle *os.File
	Scan   *bufio.Scanner
	//Open   bool
	Logline
}

// ByTime is something like an alias for a slice of Logfiles for easy sorting
type ByTime []Logfile

func (t ByTime) Len() int           { return len(t) }
func (t ByTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTime) Less(i, j int) bool { return t[i].Timestamp.Before(t[j].Timestamp) }

// parse log lines that will all look like
// Jan 26 11:53:05 combo sendmail: sendmail shutdown failed
// Jan 26 11:53:06 combo sendmail: sm-client shutdown failed
var logfiles []Logfile

func main() {
	files, err := ioutil.ReadDir("./data")
	if err != nil {
		log.Fatal(err)
	}

	// scan the target dir and init the logfiles slice
	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), ".log") {
			ll := Logfile{}
			ll.Name = "data/" + f.Name()
			logfiles = append(logfiles, ll)
		}
	}
	// initial - open the log files && fetch head line from each file
	for e := range logfiles {
		h, herr := os.Open(logfiles[e].Name)
		if herr != nil {
		}
		logfiles[e].Handle = h
		logfiles[e].Scan = bufio.NewScanner(logfiles[e].Handle)
		defer logfiles[e].Handle.Close()
		getNextLine(&logfiles[e])
	}

	for len(logfiles) > 0 {
		sort.Sort(ByTime(logfiles))
		fmt.Println(logfiles[0].Name, " : ", logfiles[0].Body, " (", logfiles[0].Timestamp.String(), ")")
		status := getNextLine(&logfiles[0])
		if !status {
			// if logfiles has only 1 entry remaining, we're done
			if len(logfiles) == 1 {
				break
			} else {
				logfiles = append(logfiles[:0], logfiles[1:]...)
			}
		}
	}
}

func getNextLine(lf *Logfile) bool {
	var err error

	cont := lf.Scan.Scan()
	if !cont {
		return cont
	}
	lf.Body = lf.Scan.Text()
	fields := strings.Fields(lf.Body)
	timeBuf := fmt.Sprintf("%s %s %s %s", fields[0], fields[1], fields[2], "2017")
	lf.Timestamp, err = time.Parse("Jan 2 15:04:05 2006", timeBuf)
	if err != nil {
		fmt.Println("err: ", err)
		return false
	}
	return true
}
