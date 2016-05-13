package main

import (
	"flag"
	"go/build"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"strings"
)

var (
	addr      = flag.String("addr", "0.0.0.0:8080", "http service address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
	homeTempl *template.Template
)

func defaultAssetPath() string {
	p, err := build.Default.Import("github.com/garyburd/gary.burd.info/go-websocket-chat", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

var globalMapRemote map[string]int
var hubsMap map[string]*hub
func homeHandler(c http.ResponseWriter, req *http.Request) {
	if(globalMapRemote[strings.Split(req.RemoteAddr,":")[0]] == 0){
		globalMapRemote[strings.Split(req.RemoteAddr,":")[0]] = 1
		hubsMap[strings.Split(req.RemoteAddr,":")[0]] = newHub()
		go hubsMap[strings.Split(req.RemoteAddr,":")[0]].run()
	}
	homeTempl.Execute(c, req.Host)
}

func main() {
	globalMapRemote = make(map[string]int)
	hubsMap = make(map[string]*hub)
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "home.html")))
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", wsHandler{hubsMap: hubsMap})
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
