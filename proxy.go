/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Tue Feb 24 14:32:24 2015 Marin Alcaraz
 */

package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	target := req.URL.String()
	fmt.Printf("Target [%s]\n", target)
	for key := range req.Header {
		w.Header().Set(key, req.Header.Get(key))
	}
	//for key := range req.Header {
	//for _, v := range req.Header[key] {
	//fmt.Printf("[%s]=>%s\n", key, v)
	//w.Header().Add(key, v)
	//}
	//}
	//proxyRequest := goreq.Request{Uri: target}
	//for key := range req.Header {
	//proxyRequest.AddHeader(key, req.Header.Get(key))
	//fmt.Printf("%s=>%s\n", key, req.Header.Get(key))
	//}
	client := &http.Client{}

	proxyRequest, _ := http.NewRequest(req.Method, target, nil)

	for key := range req.Header {
		for _, v := range req.Header[key] {
			proxyRequest.Header.Add(key, v)
		}
	}

	proxyResponse, _ := client.Do(proxyRequest)

	for key := range proxyResponse.Header {
		for _, v := range proxyResponse.Header[key] {
			w.Header().Add(key, v)
		}
	}
	w.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	proxyResponse.Write(w)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("[!]Local service binded on :8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
