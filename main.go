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
	var threadId [5]int
	var threadIndex int
	var threadMsgCount [5]uint64
	var threadMsgTime [5]time.Duration
	var threadMsgStartTime [5]time.Time

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var conn, _ = upgrader.Upgrade(w, r, nil)

		go func(conn *websocket.Conn) {
			var threadIndexForCounter = threadIndex
			threadId[threadIndexForCounter] = threadIndex
			threadIndex++
			tNow = time.Now()
			threadMsgStartTime[threadIndexForCounter] = time.Now()
			for {
				mType, msg, err := conn.ReadMessage()
				if err != nil {
					tCurrent = time.Since(tNow)
					threadMsgTime[threadIndexForCounter] = time.Since(threadMsgStartTime[threadIndexForCounter])
					fmt.Println("Before Closing connection Total number of requests received: ", messageCounter, "In time: ", tCurrent, mType)
					fmt.Println("TotalThreads:", threadIndex)
					conn.Close()
					return
				}
				messageCounter++
				threadMsgCount[threadIndexForCounter]++
				//fmt.Println(mType, string(msg))
				fmt.Println("MsgSize ", len(msg), "Message number ", messageCounter, "Time Elapsed: ", time.Since(tNow))
				//fmt.Println("Total number of requests received: ", messageCounter, "In time: ", time.Since(tNow))
				/*
					err = conn.WriteMessage(mType, msg)
					if err != nil {
						tCurrent = time.Since(tNow)
						fmt.Println("Total number of requests received: ", messageCounter, "In time: ", tCurrent)
					}*/
			}
			tCurrent = time.Since(tNow)
			//fmt.Println("Total number of requests received: ", messageCounter, "In time: ", tCurrent)
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
		fmt.Println("Total Threads ", threadIndex)

		for index := 0; index <= threadIndex; index++ {
			threadMsgTime[index] = time.Since(threadMsgStartTime[index])
			fmt.Println("Thread Index: ", index, " Message Count: ", threadMsgCount[index], "Time Elapsed: ", threadMsgTime[index])
		}
		os.Exit(1)
	}()
	log.Fatal(http.ListenAndServe(":9000", nil))
	//	http.ListenAndServe(":9000", nil)
}
