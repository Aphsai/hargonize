package main

import (
	"net/http"
	"fmt"
	"reflect"
	"io"
	"io/ioutil"
	"os"
	"path"
	"flag"
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

func main() {
	directory := os.Getenv("HOME") + "/.hargonize"
	err := os.Chdir(directory)
	pUrl := flag.String("url", "", "URL to be processed")
	flag.Parse()
	url := *pUrl
	if url == "" {
		fmt.Fprintf(os.Stderr, "Error: empty URL\n")
		return
	}

	filename := path.Base(url)
	//fmt.Println("Checking if " + filename + " exists...")
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		err := download(url, filename)
		if err != nil {
			panic(err)
		}
		fmt.Println(filename + " created!")
	} else {
		//Compare the two files, and if different, output updated
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		download(url, filename)
		updated_file, err := ioutil.ReadFile(filename)
		fmt.Println("Comparing previous version of " + filename)
		if reflect.DeepEqual(file, updated_file) {
			fmt.Println(filename + " same")
		} else {
			fmt.Println(filename + " updated")
		}
	}
}


