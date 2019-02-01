package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("gitignore downloader", "downloads and create gitignore file")

	iflist := parser.Flag("l", "list", &argparse.Options{Required: false, Help: "Lists all available gitignore files"})
	searchlist := parser.List("s", "search", &argparse.Options{Required: false, Help: "Search for a list of gitignore files in repo"})
	downloadlst := parser.List("c", "create", &argparse.Options{Required: false, Help: "Downloads and merge gitignore files in given list"})
	saveFlag := parser.Flag("S", "save", &argparse.Options{Required: false, Help: "Saves results to .gitignore in current dir"})

	err := parser.Parse(os.Args)

	if *iflist == true {
		err := list(os.Stdout)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(*searchlist) > 0 {
		_, err := search(*searchlist, os.Stdout)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	if len(*downloadlst) > 0 {
		var w *os.File
		if *saveFlag == true {
			file, err := os.OpenFile(".gitignore", os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			w = file
		} else {
			w = os.Stdout
		}

		err := findAndDownload(*downloadlst, w)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	fmt.Print(parser.Usage(nil))

	if err != nil {
		fmt.Print(parser.Usage(err))
	}
}

type FileList struct {
	Name         string
	Type         string
	Path         string
	Download_url string
}

type EntryIndex struct {
	Name string
	URL  string
}

func buildIndex() ([]EntryIndex, error) {
	endpoints := []string{
		"https://api.github.com/repos/github/gitignore/contents",
		"https://api.github.com/repos/github/gitignore/contents/Global?ref=master",
	}

	var mainIndex []EntryIndex
	for _, endpoint := range endpoints {
		fmt.Println("Downloading current index ", endpoint)
		defer fmt.Println("Downloaded current index", endpoint)

		index, err := buildSingleIndex(endpoint)
		if err != nil {
			return nil, err
		}

		mainIndex = append(mainIndex, index...)
	}
	return mainIndex, nil
}

func buildSingleIndex(endpoint string) ([]EntryIndex, error) {
	var index []EntryIndex

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	var lst []FileList
	err = json.NewDecoder(resp.Body).Decode(&lst)

	if err != nil {
		return nil, err
	}

	rx := regexp.MustCompile(`^(?P<name>[a-zA-Z0-9]+)\.gitignore`)
	for _, entry := range lst {
		if entry.Type != "file" {
			continue
		}
		if !rx.MatchString(entry.Name) {
			continue
		}
		// rx.MatchString(entry.Name)
		names := rx.FindStringSubmatch(entry.Name)

		e := EntryIndex{
			Name: names[1],
			URL:  entry.Download_url,
		}
		index = append(index, e)
	}

	return index, nil
}

func download(entry EntryIndex) (string, error) {
	resp, err := http.Get(entry.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "", err
	}
	return string(bodyBytes), nil

}

func search(strs []string, iowriter io.Writer) ([]EntryIndex, error) {
	index, err := buildIndex()
	if err != nil {
		return nil, nil
	}
	var result []EntryIndex

	for _, query := range strs {
		for _, e := range index {
			contain := strings.Contains(strings.ToLower(e.Name), strings.ToLower(query))
			if contain {
				result = append(result, e)
			}
		}
	}

	for _, e := range result {
		fmt.Fprintln(iowriter, e.Name)
	}

	return result, nil
}

func findAndDownload(strs []string, iowriter io.Writer) error {
	index, err := buildIndex()
	if err != nil {
		return err
	}

	var toDownload []EntryIndex
	for _, query := range strs {
		for _, entry := range index {
			if strings.ToLower(query) == strings.ToLower(entry.Name) {
				toDownload = append(toDownload, entry)
			}
		}
	}

	fmt.Fprintln(iowriter, "# Found these files:")

	for _, entry := range toDownload {
		fmt.Fprintln(iowriter, "# "+entry.Name+".gitignore")
	}

	var resultStr []string
	for _, entry := range toDownload {
		partial, err := download(entry)
		if err != nil {
			return err
		}
		resultStr = append(resultStr, "\n"+"# "+entry.URL+"\n")
		resultStr = append(resultStr, partial)
	}

	fmt.Fprint(iowriter, strings.Join(resultStr, "\n"))
	return nil
}

func list(iowriter io.Writer) error {
	index, err := buildIndex()
	if err != nil {
		return err
	}
	for _, e := range index {
		fmt.Fprintln(iowriter, fmt.Sprintf("%s\n", e.Name))
	}
	return nil
}
