package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/comail/colog"
	"github.com/gorilla/mux"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

type roomSettings struct {
	Host string
	Name string
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	if strings.HasSuffix(t.filename, ".js") {
		w.Header().Add("Content-Type", "application/x-javascript")
	}
	s := roomSettings{Host: r.Host, Name: name}
	t.templ.Execute(w, s)
}

var rooms = map[string]*room{}

func signaling(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	room, exist := rooms[name]

	if !exist {
		log.Println("debug: Create new room", name)
		room = newRoom()
		go room.run(func() {
			log.Println("debug: Delete room", name)
			delete(rooms, name)
		})
		rooms[name] = room
	}
	room.ServeHTTP(w, r)
}

func main() {
	colog.SetDefaultLevel(colog.LDebug)
	colog.Register()

	var addr = flag.String("addr", ":8080", "The listen port of the application.")
	flag.Parse()

	router := mux.NewRouter()

	router.HandleFunc("/WebRTCHandsOn/{name}", (&templateHandler{filename: "index.html"}).ServeHTTP)
	router.HandleFunc("/WebRTCHandsOn/{name}/webrtc.js", (&templateHandler{filename: "webrtc.js"}).ServeHTTP)
	router.HandleFunc("/WebRTCHandsOnSig/{name}", signaling)

	http.Handle("/", router)

	log.Println("info: Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
