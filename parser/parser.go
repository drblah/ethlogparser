package parser

import (
	"log"
	"regexp"
	"time"
)

const (
	MinedBlock = iota
	PropagatedBlock
	AnnouncedBlock
)

// SplitByCol takes one line from getl log and splits it into columns
func SplitByCol(str string) []string {
	var re = regexp.MustCompile(`(?m)(?P<col1>^.+?\]) (?P<col2>.{1,40}) (?P<col3>.+$)?`)

	return re.FindStringSubmatch(str)
}

// ParseStatusTimestamp takes the first column of the geth log and parses it into a status and a timestamp
func ParseStatusTimestamp(str string) (string, time.Time) {
	var re = regexp.MustCompile(`(^\S+)\[(\d\d-\d\d\|\d\d:\d\d:\d\d.\d\d\d)\]`)
	statusTimeString := re.FindStringSubmatch(str)

	timeStamp, err := time.Parse("01-02|15:04:05.000", statusTimeString[2])

	if err != nil {
		log.Println("Unable to parse timeStamp. Offending string was ", statusTimeString[2])
		log.Fatal(err)
	}

	return statusTimeString[1], timeStamp
}
