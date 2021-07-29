package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/pelletier/go-toml/v2"
)

type dir_struct struct {
	Dirname    string
	Filenames  []fileinfo_internal
	Dotdot     string
	Thumbnails bool
	Path       string
}

type config_format struct {
	Dir  string
	Path string
	Port string
}

type fileinfo_internal struct {
	Name string
	Ext  string
	Dir  bool
}

// Config
/*var directory string
var path string
var port string*/
var config config_format

// Session
var session *sessions.Session
var secret = []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

func static(w http.ResponseWriter, r *http.Request) {

	file_name := strings.Replace(r.URL.Path, "/sysreserved-static/", "", 1)

	file, err := os.Open("static/" + file_name)

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(404)
		return
	}

	file_buffer, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(500)
		return
	}

	content_type := http.DetectContentType(file_buffer)

	w.Header().Set("Content-Type", content_type)
	http.ServeFile(w, r, "static/"+file_name)

	//fmt.Printf("File: %s\nr.URL.Path: %s\n", file, r.URL.Path)
	//w.WriteHeader(200)
}

func render(w http.ResponseWriter, r *http.Request) {
	slashless_path := config.Path[1:]
	dir := strings.Replace(r.URL.Path, "/", "", 1)
	dir = strings.TrimPrefix(dir, slashless_path)
	pubdir := strings.TrimPrefix(dir, config.Dir)

	//fmt.Printf("dir: %s\npubdir %s\npath %s\n", dir, pubdir, config.Path)

	if dir == "" {
		dir = "/"
	}
	dir = config.Dir + dir

	fileinfo, err := os.Stat(dir)

	if err != nil {
		//404
		fmt.Println(err.Error())
		fmt.Println("Rendering 404 now...")
		w.WriteHeader(http.StatusNotFound)

		filestruct := struct {
			Filename string
		}{
			Filename: pubdir}

		tmpl, _ := template.ParseFiles("templates/404.html")
		tmpl.Execute(w, filestruct)

		return
	}

	if fileinfo.IsDir() {
		//io.WriteString(w, fmt.Sprintf("%s is a directory", dir))

		var pubdir_ends byte
		var pubdir_starts byte

		if string(pubdir) == "" {
			pubdir_ends = '/'
			pubdir_starts = '/'
		} else {
			pubdir_ends = pubdir[len(pubdir)-1]
			pubdir_starts = pubdir[0]
		}

		//fmt.Printf("Dir is %s\nDir ends with: %s\nDir starts with: %s\n", pubdir, string(pubdir_ends), string(pubdir_starts))

		if pubdir_ends != '/' {
			var redir string
			if pubdir_starts == '/' {
				redir = string(pubdir) + "/"
			} else {
				redir = "/" + string(pubdir) + "/"
			}
			fmt.Printf("Redirecting to %s", redir)
			http.Redirect(w, r, redir, http.StatusPermanentRedirect)
			return
		}

		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, err.Error())
		}

		var fileinfo_arr []fileinfo_internal
		var fileinfo_dir_arr []fileinfo_internal
		var fileinfo_file_arr []fileinfo_internal

		for _, fileinfo := range files {
			var ext string
			if strings.Contains(fileinfo.Name(), ".") {
				exts := strings.Split(fileinfo.Name(), ".")
				ext = exts[len(exts)-1]
			} else {
				ext = ""
			}

			fileinfo_to_append := fileinfo_internal{
				Name: fileinfo.Name(),
				Ext:  ext,
				Dir:  fileinfo.IsDir()}
			//fileinfo_arr = append(fileinfo_arr, fileinfo_to_append)
			if fileinfo_to_append.Dir {
				fileinfo_dir_arr = append(fileinfo_dir_arr, fileinfo_to_append)
			} else {
				fileinfo_file_arr = append(fileinfo_file_arr, fileinfo_to_append)
			}
		}

		fileinfo_arr = append(fileinfo_dir_arr, fileinfo_file_arr...)

		topmost_dir_arr := strings.Split(pubdir, "/")
		topmost_dir := topmost_dir_arr[len(topmost_dir_arr)-1]
		dot_dot := strings.TrimSuffix(pubdir, topmost_dir+"/")

		if r.Method == "POST" {
			r.ParseForm()
			thumbnail_string := r.FormValue("thumbnails")

			session.Put(r, "thumbnail", thumbnail_string)

			//fmt.Printf("thumbnail_bool: %s\n", thumbnail_bool)
			var thumbnail_bool bool
			if thumbnail_string == "false" {
				thumbnail_bool = false
			} else {
				thumbnail_bool = true
			}

			dirstruct := dir_struct{
				Dirname:    pubdir,
				Filenames:  fileinfo_arr,
				Dotdot:     dot_dot,
				Thumbnails: thumbnail_bool,
				Path:       config.Path}

			tmpl, err := template.ParseFiles("templates/dir.html")

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Render Error")
				return

			}

			tmpl.Execute(w, dirstruct)
		} else { // Request is GET

			thumbnail_string := session.GetString(r, "thumbnail")

			if thumbnail_string == "" {
				thumbnail_string = "true"
			}

			var thumbnail_bool bool

			if thumbnail_string == "true" {
				thumbnail_bool = true
			} else {
				thumbnail_bool = false
			}
			var dirstruct dir_struct
			dirstruct = dir_struct{
				Dirname:    pubdir,
				Filenames:  fileinfo_arr,
				Dotdot:     dot_dot,
				Thumbnails: thumbnail_bool,
				Path:       config.Path}

			tmpl, err := template.ParseFiles("templates/dir.html")

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Render Error")
				return

			}

			tmpl.Execute(w, dirstruct)
		}

		//w.Header().Set("Content-Type", "text/plain")
		//io.WriteString(w, out)

	} else { // Requested content is a file
		/*io.WriteString(w, fmt.Sprintf("%s is a file", dir))*/
		file, err := os.Open(dir)

		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, fmt.Sprintf("Error opening file"))
		}

		file_buffer, err := ioutil.ReadAll(file)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, fmt.Sprintf("Error opening file"))
		}

		content_type := http.DetectContentType(file_buffer)

		w.Header().Set("Content-Type", content_type)
		http.ServeFile(w, r, dir)

	}

	//w.Header().Set("Content-Type", "text/plain")
	//io.WriteString(w, fmt.Sprintf("Welcome to %s\n", dir))
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

	if config.Dir == "" {
		panic("FATAL: Directory not specified in config.toml")
	}

	if config.Path == "" {
		config.Path = "/"
	}

	//fmt.Printf("Directory: %s\nPath %s\nPort %s\n", directory, path, port)

	session = sessions.New(secret)
	session.Lifetime = 3 * time.Hour

	handler := http.NewServeMux()

	handler.HandleFunc(config.Path, render)
	handler.HandleFunc("/sysreserved-static/", static)

	err = http.ListenAndServe(":"+config.Port, session.Enable(handler))

	if err != nil {
		fmt.Println(err.Error())
	}
}
