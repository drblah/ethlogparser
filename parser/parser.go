package parser

import (
	"fmt"
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

	// MSGInsertedForkedBlock : Inserted forked block                    number=1  hash=e68e79…6f23a5 diff=131072 elapsed=651.016µs txs=0 gas=0 uncles=0
	MSGInsertedForkedBlock

	// MSGChainSplitDetected : Chain split detected                     number=278 hash=75e1fa…a7ee0d drop=1 dropfrom=b1ad02…79f8fb add=1 addfrom=f692d6…226951
	MSGChainSplitDetected
	// Unimplemented
	/*
		DEBUG[10-24|12:31:26.434] Inserted new block                       number=279 hash=f692d6…226951 uncles=0 txs=0 gas=0 elapsed=17.166ms
		INFO [10-24|12:31:26.434] Imported new chain segment               blocks=1  txs=0 mgas=0.000 elapsed=17.214ms     mgasps=0.000 number=279 hash=f692d6…226951 cache=27.26kB
	*/
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

// PropagatedBlock1Data contains the parsed information from the data portion of a propatagedblock1 log
type PropagatedBlock1Data struct {
	Hash       string
	Recipients int
	Duration   string
}

// PropagatedBlock2Data contains the parsed information from the data portion of a propatagedblock2 log
type PropagatedBlock2Data struct {
	ID     string
	Conn   string
	Number int
	Hash   string
	TD     int64
}

// AnnouncedBlock1Data contains the parsed information from the data portion of an AnnouncedBlock1 log
type AnnouncedBlock1Data struct {
	Hash       string
	Recipients int
	Duration   string
}

// AnnouncedBlock2Data contains the parsed information from the data portion of an AnnouncedBlock2 log
type AnnouncedBlock2Data struct {
	ID     string
	Conn   string
	Number int
	Hash   string
}

// QueuedPropagatedBlockData contains the parsed information from the data portion of an QueuedPropagatedBlock log
type QueuedPropagatedBlockData struct {
	Peer   string
	Number int
	Hash   string
	Queued int
}

// ImportingPropBlockData contains the parsed information from the data portion of an ImportingPropBlock log
type ImportingPropBlockData struct {
	Peer   string
	Number int
	Hash   string
}

// InsertedForkedBlockData contains the parsed information from the data portion of an InsertedForkedBlock log
type InsertedForkedBlockData struct {
	Number  int
	Hash    string
	Diff    int
	Elapsed string
	Txs     int
	Gas     int
	Uncles  int
}

// ChainSplitDetectedData containse the parsed information from the data portion of ChainSplitDetected log
type ChainSplitDetectedData struct {
	Number   int
	Hash     string
	Drop     int
	Dropfrom string
	Add      int
	Addfrom  string
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
	var re = regexp.MustCompile(`number=(\d+) +hash=(\S+)`)
	numberHashString := re.FindStringSubmatch(str)

	blockNumber, err := strconv.Atoi(numberHashString[1])

	if err != nil {
		log.Fatal("Failed to parse block number. Offending string: ", numberHashString[1])
	}

	minedBlock.Number = blockNumber
	minedBlock.Hash = numberHashString[2]
	return minedBlock
}

// ParsePropagatedBlock1 parses PropagatedBlock1
func ParsePropagatedBlock1(str string) (propData PropagatedBlock1Data) {
	var re = regexp.MustCompile(`hash=(\S+) recipients=(\d+) +duration=([\d\S]+)`)
	numberHashString := re.FindStringSubmatch(str)

	recipients, err := strconv.Atoi(numberHashString[2])

	if err != nil {
		log.Fatal("Failed to parse recipients number. Offending string: ", numberHashString[1])
	}

	propData.Hash = numberHashString[1]
	propData.Recipients = recipients
	propData.Duration = numberHashString[3]

	return propData
}

// ParsePropagatedBlock2 parses PropagatedBlock2
func ParsePropagatedBlock2(str string) (propData PropagatedBlock2Data) {
	var re = regexp.MustCompile(`id=([\S\d]+) conn=(\S+) number=(\d+) +hash=([\S\d]+) td=(\d+)`)
	numberHashString := re.FindStringSubmatch(str)

	number, err := strconv.Atoi(numberHashString[3])

	if err != nil {
		log.Fatal("Failed to parse number. Offending string: ", numberHashString[3])
	}

	td, err := strconv.ParseInt(numberHashString[5], 10, 64)

	if err != nil {
		log.Fatal("Failed to parse td. Offending string: ", numberHashString[5])
	}

	propData.ID = numberHashString[1]
	propData.Conn = numberHashString[2]
	propData.Number = number
	propData.Hash = numberHashString[4]
	propData.TD = td

	return propData
}

// ParseAnnouncedBlock1 parses ParseAnnouncedBlock1
func ParseAnnouncedBlock1(str string) (annData AnnouncedBlock1Data) {
	var re = regexp.MustCompile(`hash=(\S+) recipients=(\d+) +duration=([\d\S]+)`)
	numberHashString := re.FindStringSubmatch(str)

	recipients, err := strconv.Atoi(numberHashString[2])

	if err != nil {
		log.Fatal("Failed to parse recipients number. Offending string: ", numberHashString[1])
	}

	annData.Hash = numberHashString[1]
	annData.Recipients = recipients
	annData.Duration = numberHashString[3]

	return annData
}

// ParseAnnouncedBlock2 parses AnnouncedBlock1Data
func ParseAnnouncedBlock2(str string) (annData AnnouncedBlock2Data) {
	var re = regexp.MustCompile(`id=([\S\d]+) conn=(\S+) number=(\d+) +hash=([\S\d]+)`)
	numberHashString := re.FindStringSubmatch(str)

	if len(numberHashString) < 4 {
		fmt.Println("String is: ", str)
		fmt.Println("After regex: ", numberHashString)
	}

	number, err := strconv.Atoi(numberHashString[3])

	if err != nil {
		log.Fatal("Failed to parse number. Offending string: ", numberHashString[3])
	}

	annData.ID = numberHashString[1]
	annData.Conn = numberHashString[2]
	annData.Number = number
	annData.Hash = numberHashString[4]

	return annData
}

// ParseQueuedPropagatedBlock parses QueuedPropagatedBlock
func ParseQueuedPropagatedBlock(str string) (queueData QueuedPropagatedBlockData) {

	var re = regexp.MustCompile(`peer=([\S\d]+) number=(\d+) +hash=([\S\d]+) queued=(\d+)`)
	numberHashString := re.FindStringSubmatch(str)

	if len(numberHashString) < 4 {
		fmt.Println("String is: ", str)
		fmt.Println("After regex: ", numberHashString)
	}

	number, err := strconv.Atoi(numberHashString[2])

	if err != nil {
		log.Fatal("Failed to parse number. Offending string: ", numberHashString[2])
	}

	queued, err := strconv.Atoi(numberHashString[4])

	if err != nil {
		log.Fatal("Failed to parse queued. Offending string: ", numberHashString[4])
	}

	queueData.Peer = numberHashString[1]
	queueData.Number = number
	queueData.Hash = numberHashString[3]
	queueData.Queued = queued

	return queueData
}

// ParseImportingPropBlock parses ImportingPropBlock
func ParseImportingPropBlock(str string) (importData ImportingPropBlockData) {

	var re = regexp.MustCompile(`peer=([\S\d]+) number=(\d+) +hash=([\S\d]+)`)
	numberHashString := re.FindStringSubmatch(str)

	number, err := strconv.Atoi(numberHashString[2])

	if err != nil {
		log.Fatal("Failed to parse number. Offending string: ", numberHashString[2])
	}

	importData.Peer = numberHashString[1]
	importData.Number = number
	importData.Hash = numberHashString[3]

	return importData
}

// ParseInsertedForkedBlock parses InsertedForkedBlock
func ParseInsertedForkedBlock(str string) (insertData InsertedForkedBlockData) {
	var re = regexp.MustCompile(`number=(\d+) +hash=([\d\S]+) +diff=(\d+) +elapsed=([\d\S]+) +txs=(\d+) +gas=(\d+) +uncles=(\d+)`)
	numberHashString := re.FindStringSubmatch(str)

	number, err := strconv.Atoi(numberHashString[1])

	if err != nil {
		log.Fatal("Failed to parse number. Offending string: ", numberHashString[1])
	}

	diff, err := strconv.Atoi(numberHashString[3])

	if err != nil {
		log.Fatal("Failed to parse diff. Offending string: ", numberHashString[3])
	}

	txs, err := strconv.Atoi(numberHashString[5])

	if err != nil {
		log.Fatal("Failed to parse txs. Offending string: ", numberHashString[5])
	}

	gas, err := strconv.Atoi(numberHashString[6])

	if err != nil {
		log.Fatal("Failed to parse gas. Offending string: ", numberHashString[6])
	}

	uncles, err := strconv.Atoi(numberHashString[7])

	if err != nil {
		log.Fatal("Failed to parse uncles. Offending string: ", numberHashString[7])
	}

	insertData.Number = number
	insertData.Hash = numberHashString[2]
	insertData.Diff = diff
	insertData.Elapsed = numberHashString[4]
	insertData.Txs = txs
	insertData.Gas = gas
	insertData.Uncles = uncles

	return insertData
}

// ParseChainSplitDetected parses ChainSplitDetected
func ParseChainSplitDetected(str string) (splitData ChainSplitDetectedData) {
	var re = regexp.MustCompile(`number=(\d+) hash=(\S+) drop=(\d+) dropfrom=(\S+) add=(\d+) addfrom=(\S+)`)
	numberHashString := re.FindStringSubmatch(str)

	number, err := strconv.Atoi(numberHashString[1])

	if err != nil {
		log.Fatal("Failed to parse number. Offending string: ", numberHashString[1])
	}

	drop, err := strconv.Atoi(numberHashString[3])

	if err != nil {
		log.Fatal("Failed to parse drop. Offending string: ", numberHashString[3])
	}

	add, err := strconv.Atoi(numberHashString[5])

	if err != nil {
		log.Fatal("Failed to parse add. Offending string: ", numberHashString[5])
	}

	splitData.Number = number
	splitData.Hash = numberHashString[2]
	splitData.Drop = drop
	splitData.Dropfrom = numberHashString[4]
	splitData.Add = add
	splitData.Addfrom = numberHashString[6]

	return splitData
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
	case strings.Contains(str, "Inserted forked block                    "):
		logType = MSGInsertedForkedBlock
	case strings.Contains(str, "Chain split detected"):
		logType = MSGChainSplitDetected
	default:
		logType = MSGUnknown
	}

	return
}
