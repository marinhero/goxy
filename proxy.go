/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Wed Feb 25 12:00:25 2015 Marin Alcaraz
 */

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
		fmt.Printf("%s against %s\n", target, val)
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

		client := &http.Client{}

		//What is wrong with the POSTS requests?
		req.ParseForm()

		data := req.Form.Encode()
		bufferedData := bytes.NewBufferString(data)

		proxyRequest, err := http.NewRequest(req.Method,
			target, bufferedData)
		check(err)
		proxyRequest.Form = req.Form
		proxyRequest.ParseForm()
		proxyResponse, err := client.Do(proxyRequest)
		check(err)
		defer proxyRequest.Body.Close()

		w.Header().Set("Content-Type",
			proxyResponse.Header.Get("Content-Type"))
		io.Copy(w, proxyResponse.Body)

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
