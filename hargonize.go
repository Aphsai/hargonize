package main

import (
	"net/http"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"flag"
	"bytes"
)

func download(url, filename string) (err error) {
	//fmt.Println("Downloading ", url, " to ", filename)
	// Get response from URL
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	//Closes body of response when everything is done with
	defer resp.Body.Close()
	//Creates a file with the name of the url
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	//Closes out stream of os when everything is done with
	defer f.Close()
	//Copies body to file
	_, err = io.Copy(f, resp.Body)
	return
}

func compareExistingURLs(url string, filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// Create file if it does not exist
		err := download(url, filename)
		if err != nil {
			panic(err)
		}
		fmt.Println(filename + " created")
	} else {
		// Compare the two files, and if different, output updated
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		download(url, filename)
		updated_file, err := ioutil.ReadFile(filename)
//fmt.Println("Comparing previous version of " + filename)
		if bytes.Equal(file, updated_file) {
			fmt.Println(filename + " same")
		} else {
			fmt.Println(filename + " updated")
		}
	}

}

func handleFile(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("...error in reading file!")
	}
	urls := strings.Split(string(file), "\n")
	for _, url := range urls {
		if url != "" {
			compareExistingURLs(url, path.Base(url))
		}
	}
}

func main() {
	// Set directory to $HOME/.hargonize
	directory := os.Getenv("HOME") + "/.hargonize"
	err := os.Chdir(directory)
	if err != nil {
		panic(err)
	}
	// Handle flags
	pUrl := flag.String("url", "", "URL to be processed")
	pFile := flag.String("",  "", "File that contains urls")
	flag.Parse()
	url := *pUrl
	filename := *pFile
	if url != "" {
		fmt.Println("URL flag is not empty")
		compareExistingURLs(url, path.Base(url))
	} else if filename != "" {
		handleFile(filename)
	} else {
		handleFile("urls")
	}
}
