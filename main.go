package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Tensai75/nzbparser"
)

var (
	appName    = "NZBFileCleaner"
	appVersion = "" // Github tag version, e.g. "v1.0.0"
	startTime  time.Time
	wg         sync.WaitGroup
)

func init() {
	startTime = time.Now()
	parseArguments()
	if args.Verbose {
		fmt.Println(args.Version())
	}
}

func main() {
	nzbFileList, sourcePath, err := loadNZBFiles(args.NZBFile)
	if err != nil {
		exit(err)
	}
	destPath := sourcePath
	if args.DestPath != "" {
		if !pathExists(args.DestPath) {
			if err := askToCreatePath(args.DestPath); err != nil {
				exit(err)
			}
		}
		destPath = args.DestPath
	}
	for _, nzbFile := range nzbFileList {
		wg.Add(1)
		go processNZBFile(nzbFile, sourcePath, destPath)
	}
	wg.Wait()
	duration := time.Since(startTime)
	if args.Verbose {
		fmt.Printf("Processing of %v NZB files completed in %v seconds (%v file/s)\n", len(nzbFileList), duration.Seconds(), float64(len(nzbFileList))/duration.Seconds())
	}
}

// loadNZBFiles checks if the provided path is an NZB file or a directory containing NZB files.
func loadNZBFiles(nzbPath string) ([]string, string, error) {
	var nzbFiles []string

	// check if the path exists
	info, err := os.Stat(nzbPath)
	if err != nil {
		return nil, "", fmt.Errorf("file or path '%s' does not exist", nzbPath)
	}

	// if it's a file, check if it ends with .nzb
	if !info.IsDir() {
		if strings.HasSuffix(strings.ToLower(info.Name()), ".nzb") {
			return []string{info.Name()}, filepath.Dir(nzbPath), nil
		}
		return nil, "", fmt.Errorf("provided file '%s' is not an NZB file", nzbPath)
	}

	// if it's a directory, scan for .nzb files
	files, err := os.ReadDir(nzbPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read directory '%s': %v", nzbPath, err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".nzb") {
			nzbFiles = append(nzbFiles, file.Name())
		}
	}

	if len(nzbFiles) == 0 {
		return nil, "", fmt.Errorf("no NZB files found in directory '%s'", nzbPath)
	}

	return nzbFiles, nzbPath, nil
}

// processNZBFile processes the NZB file based on the provided arguments and writes the modified file to disk.
func processNZBFile(nzbFile string, sourcePath string, destPath string) {
	defer wg.Done()

	// read the content of the NZB file
	if args.Verbose {
		fmt.Printf("Processing NZB file '%s'\n", nzbFile)
	}
	content, err := os.ReadFile(filepath.Join(sourcePath, nzbFile))
	if err != nil {
		fmt.Printf("Error: %v\n", fmt.Errorf("failed to read file '%s': %v", nzbFile, err))
		return
	}

	filename, password := extractFilenameAndPassword(nzbFile)

	// parse the content of the NZB file
	nzb, err := nzbparser.ParseString(string(content))
	if err != nil {
		fmt.Printf("Error: %v\n", fmt.Errorf("failed to parse NZB file '%s': %v", nzbFile, err))
		return
	}

	// process the NZB file based on the provided arguments
	if password != "" && args.AddPwToMeta {
		if args.Verbose {
			fmt.Printf("Adding password '%s' to metadata of NZB file '%s'\n", password, nzbFile)
		}
		nzb.Meta["password"] = password
	}
	if args.RemovePwFromMeta {
		if args.Verbose {
			fmt.Printf("Removing password from metadata of NZB file '%s'\n", nzbFile)
		}
		delete(nzb.Meta, "password")
	}
	if args.AddTitleToMeta {
		if args.Verbose {
			fmt.Printf("Adding title '%s' to metadata of NZB file '%s'\n", filename, nzbFile)
		}
		nzb.Meta["title"] = filename
	}
	if args.RemoveTitleFromMeta {
		if args.Verbose {
			fmt.Printf("Removing title from metadata of NZB file '%s'\n", nzbFile)
		}
		delete(nzb.Meta, "title")
	}
	if nzb.Meta["title"] != "" && args.UseTitleForFilename {
		if args.Verbose {
			fmt.Printf("Using title '%s' from metadata as filename for NZB file '%s'\n", nzb.Meta["title"], nzbFile)
		}
		if !isValidFilename(nzb.Meta["title"]) {
			fmt.Printf("Error: %v\n", fmt.Errorf("title '%s' for NZB file '%s' contains invalid characters", nzb.Meta["title"], nzbFile))
		} else {
			filename = nzb.Meta["title"]
		}
	}
	if nzb.Meta["password"] != "" && args.AddPwToFilename {
		if args.Verbose {
			fmt.Printf("Adding password '%s' from metadata to the filename of NZB file '%s'\n", nzb.Meta["password"], nzbFile)
		}
		if !isValidFilename(nzb.Meta["password"]) {
			fmt.Printf("Error: %v\n", fmt.Errorf("password '%s' for NZB file '%s' contains invalid characters", nzb.Meta["password"], nzbFile))
		} else {
			password = nzb.Meta["password"]
		}
	}
	if args.RemovePwFromFilename {
		if args.Verbose {
			fmt.Printf("Removing the password from the filename of NZB file '%s'\n", nzbFile)
		}
		password = ""
	}

	// write the modified NZB file back to disk
	outputFile := filepath.Join(destPath, filename)
	if password != "" {
		outputFile += "{{" + password + "}}"
	}
	outputFile += ".nzb"
	if args.Verbose {
		fmt.Printf("Writing new NZB file '%s' to disk\n", outputFile)
	}
	nzbString, err := nzbparser.WriteString(nzb)
	if err != nil {
		fmt.Printf("Error: %v\n", fmt.Errorf("failed to generate new NZB file '%s': %v", nzbFile, err))
		return
	}
	err = os.WriteFile(outputFile, []byte(nzbString), 0644)
	if err != nil {
		fmt.Printf("Error: %v\n", fmt.Errorf("failed to write new NZB file '%s': %v", outputFile, err))
		return
	}
}

// extractFilenameAndPassword checks if the NZB file path has the structure filename{{password}}.nzb
// and returns the filename and password.
func extractFilenameAndPassword(nzbFile string) (string, string) {
	// define a regex pattern to match the structure filename{{password}}.nzb
	pattern := `^(.*)\{\{(.+)\}\}\.nzb$`
	re := regexp.MustCompile(pattern)

	// extract the base name of the file from the full path
	baseName := filepath.Base(nzbFile)

	// extract matches
	matches := re.FindStringSubmatch(baseName)
	if len(matches) == 3 {
		// extract filename and password
		filename := matches[1]
		password := matches[2]
		return filename, password
	}

	// if the structure is not present, return the filename without the .nzb extension and an empty password
	filename := strings.TrimSuffix(baseName, ".nzb")
	return filename, ""
}

// isValidFilename checks if a string contains invalid characters for filenames using a regular expression.
func isValidFilename(input string) bool {
	// define a regex pattern to match invalid characters for filenames
	// windows invalid characters: \ / : * ? " < > |
	pattern := `[\\/:*?"<>|]`
	re := regexp.MustCompile(pattern)

	// check if the input contains any invalid characters
	return !re.MatchString(input)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false // path does not exist
	}
	return err == nil // path exists if no error occurred
}

func askToCreatePath(path string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Destination path '%s' does not exist.\n", args.DestPath)
	fmt.Println("Do you want to create it? (y/n): ")

	// read user input
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	// trim whitespace and convert to lowercase
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "y" {
		// create the directory
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create path '%s': %v", path, err)
		}
		fmt.Printf("Path '%s' created successfully.\n", path)
	} else if input == "n" {
		return fmt.Errorf("path creation declined by user")
	} else {
		fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		return askToCreatePath(path) // Retry on invalid input
	}
	return nil
}

func exit(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
