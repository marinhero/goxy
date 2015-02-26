/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Thu Feb 26 18:33:38 2015 Marin Alcaraz
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
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

//Satisfy the Sort interface for blacklist type

func (bl blackList) Len() int {
	return len(bl)
}

func (bl blackList) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}

func (bl blackList) Less(i, j int) bool {
	return len(bl[i]) < len(bl[j])
}

func (bl *blackList) contains(target string) bool {
	for _, val := range *bl {
		//TODO Regular expressions
		if val == target {
			return true
		}
	}
	return false
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
			sort.Sort(blackList(hostBlackList))
			return
		default:
			log.Fatal(err)
		}
	}

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

		//Make the request trough our new defaultclient
		proxyResponse, err := client.Do(req)
		defer proxyResponse.Body.Close()
		check(err)

		//Start to populate the fields of the response writer header
		for key := range proxyResponse.Header {
			w.Header().Add(key, proxyResponse.Header.Get(key))
		}

		// Header returns the header map that will be sent by WriteHeader.
		// Changing the header after a call to WriteHeader (or Write) has
		// no effect.
		w.WriteHeader(proxyResponse.StatusCode)
		_, err = io.Copy(w, proxyResponse.Body)

		check(err)
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
