package main

import "log"

func checkError(e error) {
	if e != nil {
		log.Println(e)
	}
}
