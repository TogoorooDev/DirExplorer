package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/pelletier/go-toml/v2"
)

// Config
/*var directory string
var path string
var port string*/
var config config_format

// Session
var session *sessions.Session
var secret = []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

// func static_svg(w http.ResponseWriter, r *http.Request)

func static_svg(w http.ResponseWriter, r *http.Request) {

	file_name := strings.Replace(r.URL.Path, "/sysreserved-static/icons-svg", "", 1)

	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeFile(w, r, "static/icons-svg/"+file_name)
	//w.WriteHeader(200)
}

func main() {

	config_file_handle, err := os.Open("config.toml")
	if err != nil {
		fmt.Printf("Cannot read config.toml: %s", err.Error())
		return
	}
	config_file_data_raw := make([]byte, 10000)
	count, err := config_file_handle.Read(config_file_data_raw)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	config_file_data := config_file_data_raw[:count]

	err = toml.Unmarshal(config_file_data, &config)

	if err != nil {
		fmt.Printf("Cannot parse configuration data: %s\n", err.Error())
		return
	}

	if config.Server.Dir == "" {
		panic("FATAL: Directory not specified in config.toml")
	}

	//fmt.Printf("Directory: %s\nPath %s\nPort %s\n", directory, path, port)

	session = sessions.New(secret)
	session.Lifetime = 3 * time.Hour

	handler := http.NewServeMux()

	handler.HandleFunc("/", render)
	handler.HandleFunc("/sysreserved-static/icons-svg/", static_svg)

	err = http.ListenAndServe(":"+config.Server.Port, session.Enable(handler))

	if err != nil {
		fmt.Println(err.Error())
	}
}
