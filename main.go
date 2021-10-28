package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var inputDir string
	var outputDir string
	var dataFile string

	flag.StringVar(&inputDir, "i", "", "Specify the input directory which holds the perk images.")
	flag.StringVar(&outputDir, "o", "", "Specify the output directory where the renamed files should be copied to.")
	flag.StringVar(&dataFile, "d", "", "Specify the config file holding the new perk names. (format: \\d\\d\\_(\\w) - $1 is searched for in the input files.")

	flag.Parse()

	checkParameters(inputDir, outputDir, dataFile)

	_, _ = pterm.DefaultSpinner.Start("Fetching input files.")

	var inputFiles []string
	err := filepath.Walk(inputDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			inputFiles = append(inputFiles, path)
			return nil
		})
	if err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}
	pb, _ := pterm.DefaultProgressbar.
		WithTotal(len(inputFiles)).
		WithTitle("Matching file names").
		WithShowCount(true).
		Start()

	pattern, err := os.ReadFile(dataFile)
	if err != nil {
		pterm.Error.Println(err.Error())
		os.Exit(1)
	}

	var copyTable [][]string
	var missing [][]string
	nameData := strings.Split(strings.ReplaceAll(string(pattern), "\r", ""), "\n")

	for _, nd := range nameData {
		needle := strings.ToLower(nd[3:])
		target := nd
		needles := strings.Split(nd, ";")
		if len(needles) > 1 {
			// in case of different names in the source, use target;source in the txt file
			needle = strings.ToLower(needles[1][3:])
			target = needles[0]
		}

		var result []string
		for _, f := range inputFiles {
			haystack := strings.ToLower(f)
			if strings.Contains(haystack, needle) {
				result = []string{f, outputDir + "/" + target + ".png"}
				break
			}
		}
		if len(result) > 0 {
			copyTable = append(copyTable, result)
			continue
		}
		missing = append(missing, needles)
	}

	_, _ = pb.Stop()
	if len(missing) > 0 {
		pterm.DefaultTable.WithBoxed().WithData(missing).Render()
		pterm.Error.Printfln("%d files could not be matched, please fix the entries above.", len(missing))
		os.Exit(1)
	}
	pterm.Info.Println("Matched:", len(copyTable))
	pbc, _ := pterm.DefaultProgressbar.
		WithTotal(len(copyTable)).
		WithTitle("Copying file").
		WithShowCount(true).
		Start()
	var copyErrs [][]string
	for _, pair := range copyTable {
		_, err := copy(pair[0], pair[1])
		if err != nil {
			copyErrs = append(copyErrs, []string{err.Error()})
		}
		pbc.Increment()
	}
	pbc.Stop()
	if len(copyErrs) > 0 {
		pterm.DefaultTable.WithBoxed().WithData(copyErrs).Render()
		pterm.Error.Println("Some files could not be copied, review list above")
	}
	pterm.Info.Println("Copied:", len(copyTable)-len(copyErrs))
}

func checkParameters(inputDir, outputDir, dataFile string) {
	if inputDir == "" || outputDir == "" {
		pterm.DefaultBox.Println(
			pterm.Info.Sprintln("Input and output folder needed."),
			pterm.Sprintln(),
			pterm.Sprintln("Usage: dpconv -i /path/to/input/folder -o /path/to/output/folder"),
			pterm.Sprintln(),
			pterm.Sprintln("Input Folder: Specify the input directory which holds the perk images."),
			pterm.Sprint("Output Folder: Specify the output directory where the renamed files should be copied to."),
		)
		os.Exit(0)
	}

	if _, err := os.Stat(dataFile); errors.Is(err, os.ErrNotExist) {
		pterm.Error.Println("Data file does not exist")
		os.Exit(1)
	}

	if _, err := os.Stat(inputDir); errors.Is(err, os.ErrNotExist) {
		pterm.Error.Println("Input folder does not exist")
		os.Exit(1)
	}

	if _, err := os.Stat(outputDir); errors.Is(err, os.ErrNotExist) {
		pterm.Error.Println("Output folder does not exist")
		os.Exit(1)
	}
}

// https://opensource.com/article/18/6/copying-files-go
func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
