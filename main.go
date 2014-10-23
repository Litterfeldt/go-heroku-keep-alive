package main

import (
	"fmt"
	curl "github.com/andelf/go-curl"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	go wakeLoop()
	http.HandleFunc("/", wakeSelf)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("Error: %v", err)
	}
}

func wakeSelf(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "I'm Awake")
}

func wakeLoop() {
	for {
		fmt.Println("Wake up! Polling urls...")
		fmt.Println()
		wakeMachines()
		time.Sleep(30 * time.Second)
	}
}

func wakeMachines() {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], "WAKE_") {
			fmt.Println("Waking ", pair[1])
			success := wake(pair[1])
			if success {
				fmt.Println("Successfully contacted ", pair[1])
			} else {
				fmt.Println("Could not wake ", pair[1])
			}
			fmt.Println()
		}
	}
}

func wake(url string) bool {
	easy := curl.EasyInit()
	defer easy.Cleanup()

	easy.Setopt(curl.OPT_URL, url)

	fooTest := func(buf []byte, userdata interface{}) bool {
		if len(buf) <= 0 {
			return false
		}
		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return false
	}
	return true
}
