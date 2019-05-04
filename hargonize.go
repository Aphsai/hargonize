package main

import (
	"net/http"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"os"
	"path"
	"strings"
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

func compareExistingURLs(url string, filename string) {
	fmt.Println("Checking if " + filename + " exists...")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		//Create file if it does not exist
		err := download(url, filename)
		if err != nil {
			panic(err)
		}
		fmt.Println(filename + " created")
	} else {
		//Compare the two files, and if different, output updated
		file, err := ioutil.ReadFile(filename)
		if err != nil {
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
//
//func handleURLFlag(url, filename string) (err error) {
//}
//
func handleDefault(urls []string) (err error) {
	fmt.Println(len(urls))
	for _, url := range urls {
		if url != "" {
			compareExistingURLs(url, path.Base(url))
		}
	}
	return
}
//
//func handleFileFlag() (err error) {
//}


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
	file := *pFile
	if url != "" {
		fmt.Println("URL flag is not empty")
		//filename := path.Base(url)
		//handleURLFlag(url, filename)
//		file, err := ioutil.ReadFile("urls")
//		urls := strings.Split(string(file), "\n")
//		fmt.Println(urls[0])
//		if err != nil {
//			panic(err)
//		}
//		return
	} else if file != "" {
		fmt.Println("File flag is not empty")
	} else {
		file, err := ioutil.ReadFile("urls")
		if err != nil {
			panic(err)
		}
		urls := strings.Split(string(file), "\n")
		handleDefault(urls)
	}

	//fmt.Println("Checking if " + filename + " exists...")
//	_, err = os.Stat(filename)
//	if os.IsNotExist(err) {
//		err := download(url, filename)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(filename + " created!")
//	} else {
//		//Compare the two files, and if different, output updated
//		file, err := ioutil.ReadFile(filename)
//		if err != nil {
//			panic(err)
//		}
//		download(url, filename)
//		updated_file, err := ioutil.ReadFile(filename)
//		fmt.Println("Comparing previous version of " + filename)
//		if reflect.DeepEqual(file, updated_file) {
//			fmt.Println(filename + " same")
//		} else {
//			fmt.Println(filename + " updated")
//		}
//	}
}


