package main

import (
	"CameraFileCopy/args"
	filehandler "CameraFileCopy/fileHandler"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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

	err := filepath.WalkDir(cliArgs.OrigDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v", path, err)
			return nil
		}

		relPath, err := filepath.Rel(cliArgs.OrigDir, path)
		if err != nil {
			log.Printf("Error getting relative path: %v", err)
			return nil
		}

		destPath := filepath.Join(cliArgs.DestDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				log.Printf("Error creating directory %s: %v", destPath, err)
			}
			return nil
		}

		copyWg.Add(1)
		go filehandler.CopyFile(path, destPath, &copyWg, copyChannel)
		return nil
	})

	if err != nil {
		log.Printf("Error walking directory: %v", err)
		return
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
	cutoffTime := time.Now().AddDate(0, 0, -cliArgs.DaysOld)

	err := filepath.WalkDir(cliArgs.OrigDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		fi, err := d.Info()
		if err != nil {
			log.Printf("Error getting file info for %s: %v", path, err)
			return nil
		}

		if fi.ModTime().Before(cutoffTime) {
			cleanWg.Add(1)
			go filehandler.RemoveFile(path, &cleanWg, cleanChannel)
		}
		return nil
	})

	if err != nil {
		log.Printf("Error walking directory: %v", err)
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
