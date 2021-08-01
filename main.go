package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golangcollege/sessions"
)

// Globals:

// Config
var config config_format

// Session
var session *sessions.Session
var secret = []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

// Cacheing
var filenames_cache cache_struct

// Mutex
var config_lock sync.Mutex

func static_svg(w http.ResponseWriter, r *http.Request) {

	file_name := strings.Replace(r.URL.Path, "/sysreserved-static/icons-svg", "", 1)

	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeFile(w, r, "static/icons-svg/"+file_name)
	//w.WriteHeader(200)
}

func main() {
	fmt.Println("Starting initialization process")
	fmt.Println("Reading config file")

	parse_config()

	go watch_config()

	//fmt.Printf("Directory: %s\nPath %s\nPort %s\n", directory, path, port)

	session = sessions.New(secret)
	session.Lifetime = 3 * time.Hour

	handler := http.NewServeMux()

	handler.HandleFunc("/", render)
	handler.HandleFunc("/sysreserved-static/icons-svg/", static_svg)

	err := http.ListenAndServe(":"+config.Server.Port, session.Enable(handler))

	if err != nil {
		fmt.Println(err.Error())
	}
}
