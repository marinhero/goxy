/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Tue Feb 24 19:28:36 2015 Marin Alcaraz
 */

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type blackList struct {
	records []string
}

var hostBlackList blackList

func (bl *blackList) populateBlackList() {
	bl.records = make([]string, 1)
	csvfile, err := os.Open("blacklist.csv")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below
	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, each := range rawCSVdata {
		bl.records = append(bl.records, each[0])
	}
}

func (bl *blackList) isInBlackList(target string) bool {
	for _, val := range bl.records {
		if val == target {
			return true
		}
	}
	return false
}

func handler(w http.ResponseWriter, req *http.Request) {
	target := req.URL.String()

	//Todo: pattern matching!
	if hostBlackList.isInBlackList(target) {
		io.Copy(w, nil)
	} else {

		client := &http.Client{}

		//What is wrong with the POSTS requests?
		req.ParseForm()

		data := req.Form.Encode()
		bufferedData := bytes.NewBufferString(data)

		proxyRequest, _ := http.NewRequest(req.Method, target, bufferedData)
		proxyRequest.Form = req.Form
		proxyRequest.ParseForm()
		proxyResponse, _ := client.Do(proxyRequest)

		defer proxyRequest.Body.Close()

		w.Header().Set("Content-Type", proxyResponse.Header.Get("Content-Type"))
		io.Copy(w, proxyResponse.Body)

	}

}

func main() {

	hostBlackList.populateBlackList()
	http.HandleFunc("/", handler)
	fmt.Println("[!]Local service binded on :8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
