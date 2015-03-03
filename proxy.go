/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Tue Mar 03 12:58:17 2015 Marin Alcaraz
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type blackList []*regexp.Regexp

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
	return len(bl[i].String()) < len(bl[j].String())
}

func (bl *blackList) contains(target string) bool {
	for _, val := range *bl {
		fmt.Printf("%s.MatchString(%s)", val, target)
		if val.MatchString(target) == true {
			return true
		}
	}
	return false
}

func populateBlackList(bl *blackList, blackListFilename string) {
	*bl = make(blackList, 0)
	blackListFile, err := os.Open(blackListFilename)
	check(err)

	defer blackListFile.Close()
	reader := bufio.NewReader(blackListFile)
	for {
		switch blackTarget, err := reader.ReadBytes('\n'); err {
		case nil:
			stringTarget := strings.Trim(string(blackTarget), "\n")
			*bl = append(*bl,
				regexp.MustCompile(`\w*://\w*\.*`+stringTarget))
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

	if hostBlackList.contains(target) {
		fmt.Println("[+] Filtering: ", target)
		http.NotFound(w, req)
	} else {
		client := http.DefaultClient

		//By RFC 2616 RequestURI must be empty
		req.RequestURI = ""

		//URL Scheme must be lowercase
		req.URL.Scheme = strings.Map(unicode.ToLower, req.URL.Scheme)

		//Set the cookies

		if client.Jar == nil {
			client.Jar, _ = cookiejar.New(nil)
			client.Jar.SetCookies(req.URL, req.Cookies())
		}

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
	fmt.Println(hostBlackList)
	http.HandleFunc("/", proxyRequestHandler)
	fmt.Println("[+] Local service binded on :8080/")
	err := http.ListenAndServe(":8080", nil)
	check(err)
}
