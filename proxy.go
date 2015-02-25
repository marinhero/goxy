/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Tue Feb 24 19:09:50 2015 Marin Alcaraz
 */

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
)

type blackList struct {
	records []string
}

func (bl *blackList) populateBlackList() {
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
	for k, each := range rawCSVdata {
		fmt.Printf("%s\n", each[0])
		bl.records[k] = each[0]
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	target := req.URL.String()
	//Compare with blacklist
	if target == "http://sourcemaking.com/sites/all/themes/sm7/images/logo.png" {
		io.Copy(w, nil)
	} else {

		client := &http.Client{}

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
	var hostBlackList blackList

	hostBlackList.populateBlackList()
	//http.HandleFunc("/", handler)
	//fmt.Println("[!]Local service binded on :8080/")
	//err := http.ListenAndServe(":8080", nil)
	//if err != nil {
	//log.Fatal(err)
	//}
}
