package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return dirTreeRecursive(out, path, printFiles, make([]bool, 0))
}

func dirTreeRecursive(out io.Writer, path string, printFiles bool, levels []bool) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		entries = filter(entries, func(e os.DirEntry) bool { return e.IsDir() })
	}

	levelsLength := len(levels)

	var prefixBuilder strings.Builder
	for _, printStick := range levels {
		if printStick {
			prefixBuilder.WriteRune('│')
		}
		prefixBuilder.WriteRune('\t')
	}
	preparedPrefix := prefixBuilder.String()

	newLevels := append(levels, true)
	entriesLength := len(entries)
	for idx, e := range entries {
		hinge := '├'
		if idx == entriesLength-1 {
			hinge = '└'
			newLevels[levelsLength] = false
		}

		if e.IsDir() {
			fmt.Fprintf(out, "%s%c───%s\n", preparedPrefix, hinge, e.Name())
			newPath := filepath.Join(path, e.Name())
			if err := dirTreeRecursive(out, newPath, printFiles, newLevels); err != nil {
				return err
			}
		} else {
			fileInfo, _ := e.Info()
			var size string
			if fileInfo.Size() == 0 {
				size = "empty"
			} else {
				size = fmt.Sprintf("%db", fileInfo.Size())
			}
			fmt.Fprintf(out, "%s%c───%s (%s)\n", preparedPrefix, hinge, fileInfo.Name(), size)
		}
	}

	return nil
}

func filter(elements []os.DirEntry, test func(os.DirEntry) bool) []os.DirEntry {
	result := make([]os.DirEntry, 0)
	for _, s := range elements {
		if test(s) {
			result = append(result, s)
		}
	}
	return result
}
