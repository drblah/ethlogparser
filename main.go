package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/akamensky/argparse"

	"github.com/drblah/ethlogparser/parser"
)

type logLine struct {
	timeStamp time.Time
	data      string
}

func makeLogString(timeStamp time.Time, miner string, msgType string, blockNumber int, hash string) string {
	return fmt.Sprintf("%s;%s;%s;%d;%s\n", timeStamp.Format("01-02-15:04:05.000"), miner, msgType, blockNumber, hash)
}

func getInputList(dir string) (logFiles []os.FileInfo) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		log.Fatal("Unable to open input dir. ", err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".txt" {
			logFiles = append(logFiles, f)
		}
	}

	return logFiles
}

func main() {

	// Command line argument parser
	cliParser := argparse.NewParser("ethlogparser", "Parses the output of geth and outputs it as csv.")

	var concat = cliParser.Flag("c", "concat", &argparse.Options{Help: "Concatenates all logs into one output."})
	var usestdout = cliParser.Flag("s", "stdout", &argparse.Options{Help: "Output to stdout instead of file."})
	var inputPath = cliParser.String("i", "input", &argparse.Options{Help: "Path to directory which contains the raw log files."})

	err := cliParser.Parse(os.Args)
	if err != nil {
		log.Fatal(cliParser.Usage(err))
	}

	if *concat {
		log.Fatal("We concat now")
	}

	if len(*inputPath) == 0 {
		log.Println("Input path not set. Defaulting to ./logs/")
		*inputPath = "./logs/"
	}

	inputFiles := getInputList(*inputPath)

	for _, f := range inputFiles {
		fileName := fmt.Sprintf("./logs/%s", f.Name())
		minerName := strings.TrimSuffix(f.Name(), "_log.txt")

		fmt.Println("Processing: ", fileName)

		file, err := os.Open(fileName)

		if err != nil {
			log.Fatal("Unable to open file ", fileName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		outName := fmt.Sprintf("output/%s.csv", minerName)
		outputFile, err := os.Create(outName)

		if err != nil {
			log.Fatal("Failed to open output file, ", outName)
		}
		defer outputFile.Close()

		for scanner.Scan() {
			var logLine string

			line := scanner.Text()

			logType := parser.ClassifyLogType(line)

			switch logType {
			case parser.MSGMinedBlock:
				//fmt.Println("Type: MinedBlock :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParseMinedBlock(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], logData.Number, logData.Hash)
				logLine = newLine
			case parser.MSGPropagatedBlock1:
				//fmt.Println("Type: PropagatedBlock1 :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParsePropagatedBlock1(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], -1, logData.Hash)
				logLine = newLine

			case parser.MSGPropagatedBlock2:
				//fmt.Println("Type: PropagatedBlock2 :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParsePropagatedBlock2(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], logData.Number, logData.Hash)
				logLine = newLine

			case parser.MSGQueuedPropagatedBlock:
				//fmt.Println("Type: QueuedPropagatedBlock :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParseQueuedPropagatedBlock(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], logData.Number, logData.Hash)
				logLine = newLine

			case parser.MSGAnnouncedBlock1:
				//fmt.Println("Type: AnnouncedBlock1 :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParseAnnouncedBlock1(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], -1, logData.Hash)
				logLine = newLine

			case parser.MSGAnnouncedBlock2:
				//fmt.Println("Type: AnnouncedBlock2 :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParseAnnouncedBlock2(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], logData.Number, logData.Hash)
				logLine = newLine

			case parser.MSGImportingPropBlock:
				//fmt.Println("Type: ImportingPropBlock :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParseImportingPropBlock(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], logData.Number, logData.Hash)
				logLine = newLine

			case parser.MSGInsertedForkedBlock:
				//fmt.Println("Type: MSGInsertedForkedBlock :: ", line)

				columns := parser.SplitByCol(line)
				header := parser.ParseLogHeader(columns[1])
				logData := parser.ParseInsertedForkedBlock(columns[3])

				newLine := makeLogString(header.TimeStamp, minerName, columns[2], logData.Number, logData.Hash)
				logLine = newLine

			}

			if *usestdout == true {
				fmt.Print(logLine)
			} else {
				outputFile.WriteString(logLine)
			}

		}

	}

}
