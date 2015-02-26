/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Thu Feb 26 15:48:35 2015 Marin Alcaraz
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"
)

type blackList []string

var hostBlackList blackList

func check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

}

func populateBlackList(bl *blackList, blackListFilename string) {
	//Is there any better way to do this? slice size = 1
	*bl = make(blackList, 1)
	blackListFile, err := os.Open(blackListFilename)
	check(err)

	defer blackListFile.Close()
	reader := bufio.NewReader(blackListFile)
	for {
		switch blackTarget, err := reader.ReadBytes('\n'); err {
		case nil:
			*bl = append(*bl, string(blackTarget))
		case io.EOF:
			return
		default:
			log.Fatal(err)
		}
	}

}

func (bl *blackList) contains(target string) bool {
	for _, val := range *bl {
		if val == target {
			return true
		}
	}
	return false
}

func proxyRequestHandler(w http.ResponseWriter, req *http.Request) {
	target := req.URL.String()

	//Todo: pattern matching!
	if hostBlackList.contains(target) {
		fmt.Println("[+]Filtering: ", target)
		http.NotFound(w, req)
	} else {

		client := http.DefaultClient

		//By RFC 2616 RequestURI must be empty
		req.RequestURI = ""

		//URL Scheme must be lowercase
		req.URL.Scheme = strings.Map(unicode.ToLower, req.URL.Scheme)

		proxyResponse, err := client.Do(req)
		check(err)

		w.WriteHeader(proxyResponse.StatusCode)
		for key := range proxyResponse.Header {
			w.Header().Add(key, proxyResponse.Header.Get(key))
		}
		if webContent, err := ioutil.ReadAll(proxyResponse.Body); err == nil {
			w.Write(webContent)
		}
	}

}

func main() {
	listFileName := flag.String("list", "blackList.txt",
		"New line separated host file")
	flag.Parse()

	populateBlackList(&hostBlackList, *listFileName)
	fmt.Printf("[+] Blackisted: %d hosts\n", len(hostBlackList))

	http.HandleFunc("/", proxyRequestHandler)
	fmt.Println("[+] Local service binded on :8080/")
	err := http.ListenAndServe(":8080", nil)
	check(err)
}
