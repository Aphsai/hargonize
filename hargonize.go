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
	"github.com/godbus/dbus"
	"github.com/esiqveland/notify"
)

func download(url, filename string) (err error) {
	//fmt.Println("Downloading ", url, " to ", filename)
	//Create client
	client := &http.Client{}
	// Get response from URL
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.105 Safari/537.36 Vivaldi/2.4.1488.38")
	resp, err := client.Do(request)
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
			fmt.Println("Error in reading url file! Make sure the correct protocols are present")
			panic(err)
		}
		sendNotification(filename + " was created")
	} else {
		// Compare the two files, and if different, output updated
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		download(url, filename)
		updated_file, err := ioutil.ReadFile(filename)
		if !bytes.Equal(file, updated_file) {
			sendNotification(filename + " was updated")
		}
	}
}

func sendNotification(summary string) {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Error connecting to dbus session")
		return
	}
	notification := notify.Notification {
		Summary: summary,
		ExpireTimeout: int32(1000),
	}

	_, err = notify.SendNotification(conn, notification)
	if err != nil {
		fmt.Println("Error sending notification")
		return
	}
}

func handleFile(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("...error in reading file!")
		fmt.Println("If no file is specified, please ensure the 'urls' file exists at $HOME/.hargonize")
		return
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
	if os.IsNotExist(err) {
		os.Mkdir(directory, os.ModePerm)
	}
	// Handle flags
	pUrl := flag.String("url", "", "URL to be processed")
	pFile := flag.String("f",  "urls", "File that contains urls")
	flag.Parse()
	url := *pUrl
	filename := *pFile
	if url != "" {
		fmt.Println("URL flag is not empty")
		compareExistingURLs(url, path.Base(url))
	} else {
		handleFile(filename)
	}
}
