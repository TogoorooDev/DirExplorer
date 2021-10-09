package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golangcollege/sessions"
)

func contains(container []string, containee string) bool {
	for _, a := range container {
		if a == containee {
			return true
		}
	}
	return false
}

func static_svg(w http.ResponseWriter, r *http.Request) {

	file_name := strings.Replace(r.URL.Path, "/sysreserved-static/svg", "", 1)

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
	handler.HandleFunc("/sysreserved-static/svg/", static_svg)
	//handler.HandleFunc("/sysreserved-static/cache", cache_handler)
	handler.HandleFunc("/favicon.ico", favicon)

	err := http.ListenAndServe(":"+config.Server.Port, session.Enable(handler))

	if err != nil {
		fmt.Println(err.Error())
	}
}
