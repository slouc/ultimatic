package main

import (
    "net/http"
    "io"
    "io/ioutil"
    "github.com/gorilla/websocket"
    
)

var upgrader = websocket.Upgrader{
    ReadBufferSize: 4096, 
    WriteBufferSize: 4096,
    CheckOrigin: func (r *http.Request) bool {
        return true
    },
}

func handler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    c := &connection{send: make(chan []byte, 256), ws: ws}
    s.register <- c
    go c.writer()
    c.reader()               
}

func serveHome(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type","text/html; charset=utf-8")
    io.WriteString(w, readFile("/root/ultimattt/index.html"))        
}

type connection struct {
    ws *websocket.Conn
    send chan []byte
}

type server struct {
    broadcast chan []byte
    register chan *connection
    connections map[*connection]bool
}

var s = server {
    broadcast: make(chan []byte),
    register: make(chan *connection),
    connections: make(map[*connection]bool),
}

func (s *server) run() {
    for {
        select {
            case c:= <-s.register:
                s.connections[c] = true
            case m:= <-s.broadcast:
                for c:= range s.connections {
                    select {
                        case c.send <- m:
                        default:
                            delete(s.connections, c)
                            close(c.send)
                    }
                }
        }
   }
}

func (c *connection) reader() {
    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            break
        }
        s.broadcast <- message
    }
    c.ws.Close()
}

func (c *connection) writer() {
    for message := range c.send {
        err := c.ws.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            break
        }
    }
    c.ws.Close()
}

func readFile(filename string) string {
    content, _ := ioutil.ReadFile(filename) //handle err
    //lines := strings.Split(string(content), "\n")
    return string(content)
}

func main() { 
    http.HandleFunc("/", serveHome)
    http.HandleFunc("/ws", handler)
    go s.run()
    http.ListenAndServe(":80", nil)    
}
