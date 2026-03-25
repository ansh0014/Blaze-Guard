package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	file, _ := os.Create("structure.txt")
	defer file.Close()

	root := "./"

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Fprintln(file, path)
		}
		return nil
	})
}