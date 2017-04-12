package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true

	},
}

type myStruct struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func main() {
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})*/
	var messageCounter uint64
	var tNow time.Time
	var tCurrent time.Duration

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)

		go func(conn *websocket.Conn) {
			tNow = time.Now()
			for {
				mType, msg, err := conn.ReadMessage()
				if err != nil {
					conn.Close()
					return
				}
				messageCounter++
				fmt.Println(mType, string(msg))
				fmt.Println("Total number of requests received: ", messageCounter, "In time: ", time.Since(tNow))
				/*
					err = conn.WriteMessage(mType, msg)
					if err != nil {
						tCurrent = time.Since(tNow)
						fmt.Println("Total number of requests received: ", messageCounter, "In time: ", tCurrent)
					}*/
			}
			tCurrent = time.Since(tNow)
			fmt.Println("Total number of requests received: ", messageCounter, "In time: ", tCurrent)
		}(conn)
		fmt.Println("Total number of requests received: ", messageCounter, "In time: ", tCurrent)
		go func(conn *websocket.Conn) {
			ch := time.Tick(5 * time.Second)

			for range ch {
				conn.WriteJSON(myStruct{
					"UsErNaMe", "FiRsTnAmE", "LaStNaMe",
				})
			}
		}(conn)

	})
	//handle server interrupts
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Stopping Server ...")
		tCurrent = time.Since(tNow)
		fmt.Println("Total number of requests received: ", messageCounter, "In time: ", tCurrent)
		os.Exit(1)
	}()
	log.Fatal(http.ListenAndServe(":9000", nil))
	//	http.ListenAndServe(":9000", nil)
}
