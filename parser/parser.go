package parser

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// MSGUnknown message type
	MSGUnknown = iota

	// MSGMinedBlock : 🔨 mined potential block                  number=10 hash=75d8ad…0f4a6c
	MSGMinedBlock

	// MSGPropagatedBlock1 : Propagated block                         hash=75d8ad…0f4a6c recipients=3 duration=2562047h47m16.854s
	MSGPropagatedBlock1

	// MSGPropagatedBlock2 : Propagated block                         id=d9c2b87e4525fab9 conn=inbound number=10 hash=75d8ad…0f4a6c td=1444032
	MSGPropagatedBlock2

	// MSGQueuedPropagatedBlock : Queued propagated block                  peer=d9c2b87e4525fab9 number=10 hash=75d8ad…0f4a6c queued=1
	MSGQueuedPropagatedBlock

	// MSGAnnouncedBlock1 : Announced block                          hash=75d8ad…0f4a6c recipients=9 duration=2562047h47m16.854s
	MSGAnnouncedBlock1

	// MSGAnnouncedBlock2 : Announced block                          id=c465b03a2b2aee96 conn=inbound number=10 hash=75d8ad…0f4a6c
	MSGAnnouncedBlock2

	// MSGImportingPropBlock : Importing propagated block               peer=d9c2b87e4525fab9 number=10 hash=75d8ad…0f4a6c
	MSGImportingPropBlock
)

// Header contains the parsed information from a log header
type Header struct {
	Status    string
	TimeStamp time.Time
}

// MinedBlockData contains the parsed information from the data portion of a MinedBlock log
type MinedBlockData struct {
	Number int
	Hash   string
}

// SplitByCol takes one line from getl log and splits it into columns
func SplitByCol(str string) []string {
	var re = regexp.MustCompile(`(?m)(?P<col1>^.+?\]) (?P<col2>.{1,40}) (?P<col3>.+$)?`)

	return re.FindStringSubmatch(str)
}

// ParseLogHeader takes the first column of the geth log and parses it into a status and a timestamp
func ParseLogHeader(str string) (header Header) {
	var re = regexp.MustCompile(`(^.+)\[(\d\d-\d\d\|\d\d:\d\d:\d\d.\d\d\d)\]`)
	statusTimeString := re.FindStringSubmatch(str)

	status := statusTimeString[1]

	timeStamp, err := time.Parse("01-02|15:04:05.000", statusTimeString[2])

	if err != nil {
		log.Println("Unable to parse timeStamp. Offending string was ", statusTimeString[2])
		log.Fatal(err)
	}

	return Header{status, timeStamp}
}

// ParseMinedBlock parses the data portion of a MinedBlock log
func ParseMinedBlock(str string) (minedBlock MinedBlockData) {
	var re = regexp.MustCompile(`number=(\d+) hash=(\S+)`)
	numberHashString := re.FindStringSubmatch(str)

	blockNumber, err := strconv.Atoi(numberHashString[1])

	if err != nil {
		log.Fatal("Failed to parse block number. Offending string: ", numberHashString[1])
	}

	minedBlock.Number = blockNumber
	minedBlock.Hash = numberHashString[2]
	return minedBlock
}

// ClassifyLogType determines the type of log
func ClassifyLogType(str string) (logType int) {

	switch {
	case strings.Contains(str, "🔨 mined potential block"):
		logType = MSGMinedBlock
	case strings.Contains(str, "Propagated block                         hash="):
		logType = MSGPropagatedBlock1
	case strings.Contains(str, "Propagated block                         id="):
		logType = MSGPropagatedBlock2
	case strings.Contains(str, "Queued propagated block"):
		logType = MSGQueuedPropagatedBlock
	case strings.Contains(str, "Announced block                          hash"):
		logType = MSGAnnouncedBlock1
	case strings.Contains(str, "Announced block                          id="):
		logType = MSGAnnouncedBlock2
	case strings.Contains(str, "Importing propagated block               peer="):
		logType = MSGImportingPropBlock
	default:
		logType = MSGUnknown
	}

	return
}
