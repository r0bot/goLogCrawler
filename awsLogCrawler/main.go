package main

import (
	"fmt"
	"os"
	"bufio"
	"regexp"
	"io/ioutil"
//"encoding/json"
	"compress/gzip"
	"strconv"
)

func checkIfError(e error) {
	if e != nil {
		fmt.Println("Searching for " + e.Error())
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Log must be located in ./ in gz format.")
	fmt.Println("enter string to search for")
	fmt.Println("---------------------")

	for {
		fmt.Print("-> ")
		text,_, _ := reader.ReadLine()
		searchString := string(text)
		outputFile, err := os.Create("./output"+searchString+".txt")
		checkIfError(err)

		fmt.Println("Searching for " + searchString)
		//readWriteAccessInt := syscall.O_RDWR
		files, err := ioutil.ReadDir("./logs")
		checkIfError(err)

		var filesScanned int
		var positiveMatchedLines int
		isLineAppRelated := regexp.MustCompile(searchString)
		for _, f := range files {
			filesScanned++

			file, err := os.Open("./logs/" + f.Name())
			checkIfError(err)
			defer file.Close()

			gz, err := gzip.NewReader(file)

			defer gz.Close()

			scanner := bufio.NewScanner(gz)

			// Assemble regexp


			for scanner.Scan() {

				lineOfLog := scanner.Text()
				jumpOverLine := isLineAppRelated.FindString(lineOfLog)
				//Process line as part of migration
				if jumpOverLine != "" {
					positiveMatchedLines++

					outputFile.WriteString(lineOfLog)
					outputFile.WriteString("\n")
					outputFile.WriteString("\n")
					outputFile.WriteString("\n")

				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}

		}
		fmt.Println("Files Scanned " + strconv.Itoa(filesScanned))
		fmt.Println("Matched Lines " + strconv.Itoa(positiveMatchedLines))
		outputFile.Sync()
		outputFile.Close()
	}

}