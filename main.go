package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"sync"
	"os"
	"github.com/chat/trace"
)

type templateHandler struct {
	once		sync.Once
	filename	string
	templ		*template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Do calls the function if and only if Do is called for the first time for this instance of Once
	t.once.Do(func () {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil{
		log.Fatal("ListenAndServe:", err)
	}

}
