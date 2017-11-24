package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"unicode"

	ansi "github.com/k0kubun/go-ansi"
	translit "github.com/macroblock/zl/text"
)

type meta struct {
	Project       string      `json:"project"`
	ProjectEn     string      `json:"project_en"`
	Season        int         `json:"season"`
	SeasonTitle   string      `json:"season_title"`
	Episode       int         `json:"episode"`
	EpisodeGlobal int         `json:"episode_global"`
	Title         interface{} `json:"title"`
	Description   string      `json:"description"`
	Pg            string      `json:"pg"`
	Duration      string      `json:"duration"`
	Files         struct {
		Mp4 interface{} `json:"mp4"`
		Md5 string      `json:"md5"`
		Mxf interface{} `json:"mxf"`
	} `json:"files"`
	ProviderID struct {
		ProjectID int    `json:"project_id"`
		SeasonID  int    `json:"season_id"`
		ProgramID string `json:"program_id"`
		ContentID string `json:"content_id"`
	} `json:"provider_id"`
	EfirDate       string   `json:"efir_date"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date"`
	Countries      []string `json:"countries"`
	GenerationDate string   `json:"generation_date"`
}

func help() {
	ansi.Println("\x1b[31;1mNo arguments were provided.\x1b[0m")
	ansi.Println("Pass one text file with a list of file pathes to parse.")
	ansi.Println("USAGE: tnt inputFile")
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLinesFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	args := os.Args[1:]

	// Show help info and exit if no arguments were passed.
	if len(args) < 1 {
		help()
		os.Exit(0)
	}

	files, err := readLinesFromFile(args[0])
	if err != nil {
		ansi.Println("  \x1b[31;1m", err, "\x1b[0m")
		return
	}
	argsLen := len(files)

	// Iterate over files.
	for i, filePath := range files {
		// Remove path from filename.
		fileName := path.Base(filepath.ToSlash(filePath))
		fileExt := path.Ext(fileName)
		fileBase := fileName[0 : len(fileName)-len(fileExt)]

		// Print out input fileName.
		ansi.Print(fmt.Sprintf("%0"+strconv.FormatInt(int64(math.Log10(float64(argsLen)))+1, 10)+"d", i+1) + "/" + strconv.FormatInt(int64(argsLen), 10) + "\x1b[33;1m " + fileName + "\x1b[0m")
		// Open the file.
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			ansi.Println("  \x1b[31;1m", err, "\x1b[0m")
			continue
		}

		// Fill in the metadata struct.
		m := meta{}
		err = json.Unmarshal(file, &m)
		if err != nil {
			ansi.Println("  \x1b[31;1m", err, "\x1b[0m")
			continue
		}

		// Store project name, season and episode number.
		n := m.Project
		s := m.Season
		e := m.Episode

		if n == "" {
			ansi.Println("   \x1b[31;1mJSON: project is null\x1b[0m")
			continue
		}
		if s == 0 {
			ansi.Println("   \x1b[31;1mJSON: season is null\x1b[0m")
			continue
		}
		if e == 0 {
			ansi.Println("   \x1b[31;1mJSON: episode is null\x1b[0m")
			continue
		}

		// Translit the project name.
		n, err = translit.Translit(n)
		if err != nil {
			ansi.Println("  \x1b[31;1m", err, "\x1b[0m")
			continue
		}

		// Make the first letter of the project name uppercase.
		r := []rune(n)
		r[0] = unicode.ToUpper(r[0])
		n = string(r)

		// Create the new filename.
		newFileName := fmt.Sprintf("%v_s%02de%02d_%v%v", n, s, e, fileBase, fileExt)

		// Print out the new filename.
		ansi.Println(" > \x1b[32;1m" + newFileName + "\x1b[0m")

		// Indent the json file.
		b, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			ansi.Println("  \x1b[31;1m", err, "\x1b[0m")
			continue
		}

		// Print out the json file.
		// ansi.Println(string(b))

		// Write the new json file.
		err = ioutil.WriteFile(newFileName, b, 0775)
		if err != nil {
			ansi.Println("  \x1b[31;1m", err, "\x1b[0m")
			continue
		}
	}
}
