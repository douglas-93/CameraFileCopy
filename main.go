package main

import (
	"CameraFileCopy/args"
	filehandler "CameraFileCopy/fileHandler"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

func main() {
	cliArgs := args.ParseArgs()

	if flag.NFlag() == 0 {
		args.HelpMenu()
		return
	}

	_, err := os.Stat(cliArgs.DestDir)
	if os.IsNotExist(err) {
		log.Printf("Error: Destination Directory not found: %s", cliArgs.DestDir)
		var confirm string
		fmt.Println("Do you want to create the directory? (y/n)")
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			return
		}
		os.MkdirAll(cliArgs.DestDir, 0755)
	}

	_, err = os.Stat(cliArgs.OrigDir)
	if os.IsNotExist(err) {
		log.Printf("Error: Source Directory not found: %s", cliArgs.OrigDir)
		return
	}

	if cliArgs.Clean && cliArgs.DaysOld <= 0 {
		log.Printf("Days must be greater than 0")
		return
	}

	runCopy(cliArgs)

	if cliArgs.Clean {
		runClean(cliArgs)
	}
}

func runCopy(cliArgs *args.CliArgs) {
	var copyWg sync.WaitGroup
	copyChannel := make(chan filehandler.Result, cliArgs.MaxItens)

	files, err := os.ReadDir(cliArgs.OrigDir)

	if err != nil {
		log.Printf("Error reading directory: %v", err)
		return
	}

	fmt.Println("****** Start copy process ******")
	for _, f := range files {
		if !f.IsDir() {
			copyWg.Add(1)
			go filehandler.CopyFile(path.Join(cliArgs.OrigDir, f.Name()), path.Join(cliArgs.DestDir, f.Name()), &copyWg, copyChannel)
		}
	}

	go func() {
		copyWg.Wait()
		close(copyChannel)
	}()

	for result := range copyChannel {
		if result.Success {
			// Verifica se há um erro (aviso) mesmo com sucesso
			if result.Error != nil {
				log.Printf("[ WARN ] %s: %v", result.FileName, result.Error)
			} else {
				fmt.Printf("[  OK  ] %s\n", result.FileName)
			}
		} else {
			log.Printf("[ FAIL ] %s: %v", result.FileName, result.Error)
		}
	}
	fmt.Println("****** Copy process finish ******")
}

func runClean(cliArgs *args.CliArgs) {
	var cleanWg sync.WaitGroup
	cleanChannel := make(chan filehandler.Result, cliArgs.MaxItens)

	fmt.Println("****** Start clean process ******")
	files, err := os.ReadDir(cliArgs.OrigDir)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -cliArgs.DaysOld)

	for _, f := range files {
		fi, e := f.Info()
		if e != nil {
			log.Printf("Error getting file info for %s: %v", f.Name(), e)
			continue
		}

		if !fi.IsDir() && fi.ModTime().Before(cutoffTime) {
			cleanWg.Add(1)
			go filehandler.RemoveFile(path.Join(cliArgs.OrigDir, f.Name()), &cleanWg, cleanChannel)
		}
	}

	go func() {
		cleanWg.Wait()
		close(cleanChannel)
	}()

	for result := range cleanChannel {
		if result.Success {
			fmt.Printf("[  OK  ] %s\n", result.FileName)
		} else {
			log.Printf("[ FAIL ] %s: %v", result.FileName, result.Error)
		}
	}
	fmt.Println("****** Clean process finish ******")
}
