package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatalln(err)
		}

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			// Print the msg to the console

			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/upload", ws)

	fmt.Println("Starting server on Port 5050..")

	http.ListenAndServe(":5050", nil)
}

func ws(w http.ResponseWriter, r *http.Request) {
	// TODO: Remove this checkOrigin thinge later.
	// Added it for local testing.
	fmt.Println("KSS: Handle /upload")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Upgrader Fatal error: ", err)
	}
	renderer(conn)
}

func renderer(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Fatalln("Read message error: ", err)
		}
		if msgType == 2 {
			fmt.Println("Yes got the bindata\n")
			fmt.Println("Data: ", msg)
			f := writeBinFile(msg)
			readBinaryFile(f)

		} else {
			fmt.Println("got the some other data\n")
		}
	}
}

func readBinaryFile(a string) {
	f, err := os.Open(a)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	var pi float64
	err = binary.Read(f, binary.LittleEndian, &pi)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(pi)
}

func writeBinFile(b []byte) string {
	type FileDetails struct {
		LastModified int64
		Name         string
		Size         int64
		Type         string
	}

	// exp := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	// fmt.Println("writing bytes: %X", b)
	i := bytes.Split(b, []byte("\r\n\r\n"))
	// fmt.Println("i[0] : ", i[0])
	// fmt.Println("expe    bytes: %X", exp)
	// fmt.Println("i[1] : ", i[1])

	hdr := bytes.Split(i[0], []byte("!"))
	// fmt.Println("ENCODING: ", hdr[0])
	// fmt.Println("HDR: ", hdr[1])
	// jr := json.RawMessage(hdr[1])
	// fmt.Println("RAW JSON: ", jr)
	// fmt.Println("RAW JSON: ", string(jr))

	fldetails := &FileDetails{}
	_ = json.Unmarshal(hdr[1], fldetails)
	a := fldetails.Name

	f, err := os.Create(a)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	err = binary.Write(f, binary.LittleEndian, i[1])
	if err != nil {
		log.Fatalln(err)
	}
	return a
}
