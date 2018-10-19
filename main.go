package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/drblah/ethlogparser/parser"
)

type logLine struct {
	timeStamp time.Time
	data      string
}

func makeLogString(timeStamp time.Time, miner string, msgType int, blockNumber int, hash string) string {
	return fmt.Sprintf("%s;%s;%d;%d;%s\n", timeStamp.Format("01-02-15:04:05.000"), miner, msgType, blockNumber, hash)
}

func main() {

	minerName := "minerNameChangeMe"

	fileName := "./logs/miner1_log.txt"

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal("Unable to open file ", fileName)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var logLines []string

	for scanner.Scan() {
		line := scanner.Text()

		logType := parser.ClassifyLogType(line)

		switch logType {
		case parser.MSGMinedBlock:
			//fmt.Println("Type: MinedBlock :: ", line)

			columns := parser.SplitByCol(line)
			header := parser.ParseLogHeader(columns[1])
			logData := parser.ParseMinedBlock(columns[3])

			newLine := makeLogString(header.TimeStamp, minerName, parser.MSGMinedBlock, logData.Number, logData.Hash)
			logLines = append(logLines, newLine)
		case parser.MSGPropagatedBlock1:
			//fmt.Println("Type: PropagatedBlock1 :: ", line)

			columns := parser.SplitByCol(line)
			header := parser.ParseLogHeader(columns[1])
			logData := parser.ParsePropagatedBlock1(columns[3])

			newLine := makeLogString(header.TimeStamp, minerName, parser.MSGPropagatedBlock1, -1, logData.Hash)
			logLines = append(logLines, newLine)

			fmt.Println(logLines)
			os.Exit(0)
		case parser.MSGPropagatedBlock2:
			fmt.Println("Type: PropagatedBlock2 :: ", line)
		case parser.MSGQueuedPropagatedBlock:
			fmt.Println("Type: QueuedPropagatedBlock :: ", line)
		case parser.MSGAnnouncedBlock1:
			fmt.Println("Type: AnnouncedBlock1 :: ", line)
		case parser.MSGAnnouncedBlock2:
			fmt.Println("Type: AnnouncedBlock2 :: ", line)
		case parser.MSGImportingPropBlock:
			fmt.Println("Type: ImportingPropBlock :: ", line)
		}

		//split := parser.SplitByCol(line)

		//fmt.Println(split[1])

		//logHeader := parser.ParseLogHeader(split[1])

		//fmt.Println(logHeader)

		//logFiles = append(logFiles, logLine{logHeader.TimeStamp, line})
	}

}
