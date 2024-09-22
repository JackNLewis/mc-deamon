package main

import (
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func tailLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// cmd := exec.Command("tail", "-f", "./biglogfile.txt")
	cmd := exec.Command("docker", "logs", "--follow", "mc-manager")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			log.Println(err)
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, append(buf[:n], '\n')); err != nil {
			log.Println(err)
			break
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/ws", tailLog)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
