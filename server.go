/*
** server.go
** Author: Marin Alcaraz
** Mail   <marin.alcaraz@gmail.com>
** Started on  Fri Feb 20 18:44:36 2015 Marin Alcaraz
** Last update Fri Feb 20 18:55:23 2015 Marin Alcaraz
 */

package main

import (
	"fmt"
	"net"
)

func handleConnection(connection net.Conn) {
	fmt.Println(connection)
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}
