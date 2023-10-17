package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must add argument with folder path")
	}

	folder := flag.String("f", "", "Set hosting folder.")
	port := flag.Int("port", 3000, "Set port for server.")
	flag.Parse()

	if *folder == "" {
		log.Fatal("You must add flag -f with folder path.\n `-f foldername`")
	}

	core := newCore(*folder, *port)
	core.searchFolder(*folder, "")
	core.startServer()

	for {
		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			log.Fatal(err)
		}
		if strings.ToLower(string(char)) == "r" {
			core.reload()
		}
	}
}
