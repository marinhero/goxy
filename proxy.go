/*
** proxy.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Mon Feb 23 18:32:37 2015 Marin Alcaraz
 */

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/franela/goreq"
)

func handler(w http.ResponseWriter, req *http.Request) {
	target := req.URL.String()
	for key := range req.Header {
		w.Header().Set(key, req.Header.Get(key))
		fmt.Println(w.Header().Get(key))
	}
	//for key := range req.Header {
	//for _, v := range req.Header[key] {
	//fmt.Printf("[%s]=>%s\n", key, v)
	//w.Header().Add(key, v)
	//}
	//}
	res, _ := goreq.Request{Uri: target}.Do()
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	content, _ := res.Body.ToString()
	w.Write([]byte(content))
	res.Body.Close()
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("[!]Local service binded on :8080/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
