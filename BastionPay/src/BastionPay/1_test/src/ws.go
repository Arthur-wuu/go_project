package src

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"time"

	"github.com/gorilla/websocket"
)

type Person struct {
	Name string
	Age  int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	indexFile, err := os.Open("html/index.html")
	checkErr(err)

	index, err := ioutil.ReadAll(indexFile)
	checkErr(err)

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			if string(msg) == "ping" {
				fmt.Println("ping")
				time.Sleep(time.Second * 2)
				err = conn.WriteMessage(msgType, []byte("pong"))
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				conn.Close()
				fmt.Println(string(msg))
				return
			}
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(index))
	})

	http.ListenAndServe(":3000", nil)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}