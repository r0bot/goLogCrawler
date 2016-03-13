package main

import (
    "fmt"
    "os"
    "bufio"
    "regexp"
    "io/ioutil"
	//"encoding/json"
	"errors"
)

func checkIfError(e error) {
    if e != nil {
        panic(e)
    }
}

type MigrationLog struct {
	Name string
	ApiKey string
	LogEntries string
	ErrorSeverity int
	processed bool
}

func findLogForApp (mapToSearch map[string]MigrationLog, keyToSearch string) (MigrationLog, error) {
	for k, e := range mapToSearch {
		if k == keyToSearch {
			return e, nil
		}
	}
	return MigrationLog{}, errors.New("No entry")
}

func parseAppMigrationLog () map[string]MigrationLog{
	files, err := ioutil.ReadDir("E:\\googleDrive\\Telerik\\migrationLog")
	checkIfError(err)
	problematicAppsList := map[string]MigrationLog{};

	isProblematicAppRegexp := regexp.MustCompile("Error while migrating app: (.{16})(?:.*)")
	for _, f := range files {
		file, err := os.Open("E:\\googleDrive\\Telerik\\migrationLog\\" + f.Name())
		checkIfError(err)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineOfLog := scanner.Text()
			isProblematicApp := isProblematicAppRegexp.FindAllStringSubmatch(lineOfLog, -2)
			if len(isProblematicApp) != 0 {
				apiKey := isProblematicApp[0][1]
				problematicAppsList[apiKey] = MigrationLog{apiKey, apiKey, "", 0, true}
				fmt.Println(isProblematicApp[0][1]);
			}
		}
	}
	return problematicAppsList
}

func main() {
    files, err := ioutil.ReadDir("E:\\googleDrive\\Telerik\\log")
	checkIfError(err)
	problematicAppsList := parseAppMigrationLog();

    for _, f := range files {
		file, err := os.Open("E:\\googleDrive\\Telerik\\log\\" + f.Name())
		checkIfError(err)

        scanner := bufio.NewScanner(file)



        // Assemble regexp
        isLineMigrationRelated := regexp.MustCompile("app-migration")
        migrationStartRegex := regexp.MustCompile("Start migration of application:(.{16})")
        migrationEndRegex := regexp.MustCompile("Finished migration of app:")

        migrationStartFound := false
	var currentLog = MigrationLog{}
	for scanner.Scan() {

		lineOfLog := scanner.Text()
		jumpOverLine := isLineMigrationRelated.FindString(lineOfLog)
		//Process line as part of migration
		if jumpOverLine != "" {
			migrationStartData := migrationStartRegex.FindAllStringSubmatch(lineOfLog, -2)

			if len(migrationStartData) != 0 {
				migrationStartFound = true
				for k, e := range problematicAppsList {
					if k == migrationStartData[0][1] {
						currentLog = e
					}
				}
			}

			//If log entry is started and we have selected the correct log struct add to its logEntries
			if migrationStartFound == true {
				if currentLog.processed == true {
					currentLog.LogEntries += lineOfLog
				}

				//If this is the end of log default all local vars and write to File
				migrationEndData := migrationEndRegex.FindString(lineOfLog)
				if migrationEndData != "" {
					//out, err := json.Marshal(currentLog)
					//if err != nil {
					//	panic (err)
					//}
					//TODO actually write this to a file
					fmt.Println(currentLog.LogEntries)

					//default local vars
					migrationStartFound = false
					currentLog = MigrationLog{}
				}
			}




		}
	}


	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	defer file.Close()
	}
}