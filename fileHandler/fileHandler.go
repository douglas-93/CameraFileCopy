package filehandler

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"
)

type Result struct {
	FileName string
	Success  bool
	Error    error
}

func CopyFile(src, dst string, wg *sync.WaitGroup, r chan Result) {
	defer wg.Done()

	// 1. Obter informações do arquivo de origem (necessário para metadados)
	sourceInfo, err := os.Stat(src)
	if err != nil {
		r <- Result{FileName: path.Base(src), Success: false, Error: fmt.Errorf("failed to stat source: %w", err)}
		return
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		r <- Result{FileName: path.Base(src), Success: false, Error: fmt.Errorf("failed to open source: %w", err)}
		return
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		r <- Result{FileName: path.Base(src), Success: false, Error: fmt.Errorf("failed to create destination: %w", err)}
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		r <- Result{FileName: path.Base(src), Success: false, Error: fmt.Errorf("failed to copy content: %w", err)}
		return
	}

	// Forçar a gravação no disco
	if err := destFile.Sync(); err != nil {
		r <- Result{FileName: path.Base(src), Success: false, Error: fmt.Errorf("failed to sync file to disk: %w", err)}
		return
	}

	// Copiar Metadados (Data de Modificação)
	// Usa os.Chtimes para definir a hora de acesso e a hora de modificação
	modTime := sourceInfo.ModTime()
	if err := os.Chtimes(dst, modTime, modTime); err != nil {
		logMsg := fmt.Errorf("warning: failed to set modification time (copy successful): %w", err)
		r <- Result{FileName: path.Base(src), Success: true, Error: logMsg}
		return
	}

	r <- Result{FileName: path.Base(src), Success: true, Error: nil}
}

func RemoveFile(file string, wg *sync.WaitGroup, r chan Result) {
	defer wg.Done()
	err := os.Remove(file)
	if err != nil {
		r <- Result{FileName: path.Base(file), Success: false, Error: fmt.Errorf("failed to remove file: %w", err)}
		return
	}

	r <- Result{FileName: path.Base(file), Success: true, Error: nil}
}
