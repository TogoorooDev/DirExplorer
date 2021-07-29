package main

import (
	"io"
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"strings"
	"net/http"
	"html/template"

	"github.com/golangcollege/sessions"
)

type fileinfo_internal struct{
	Name string
	Ext string
	Dir bool
	
}

// Config
var directory = "example/"
var port = ":8080"

// Session
var session *sessions.Session
var secret = []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

func render(w http.ResponseWriter, r *http.Request){
	dir := strings.Replace(r.URL.Path, "/", "", 1)
	pubdir := strings.TrimPrefix(dir, directory)

	if dir == "" { dir = "/" }
	dir = directory + dir

	fileinfo, err := os.Stat(dir)

	if err != nil {
		//404
		fmt.Println(err.Error())
		fmt.Println("Rendering 404 now...")
		w.WriteHeader(http.StatusNotFound)

		filestruct := struct {
			Filename string
		} {
			Filename: pubdir}
		
		tmpl, _ := template.ParseFiles("templates/404.html")
		tmpl.Execute(w, filestruct)
		
		return
	}

	if fileinfo.IsDir(){
		//io.WriteString(w, fmt.Sprintf("%s is a directory", dir))

		/*var pubdir_ends byte

		if pubdir == "/"{
			pubdir_ends = '/'
		}else {
			pubdir_ends = pubdir[len(pubdir) - 1]
		}

		fmt.Printf("Dir is %s\nDir ends with: %s\n", pubdir, string(pubdir_ends))*/
		/*
		if string(pubdir_ends) != "/" {
			redir := string(pubdir + "/")
			fmt.Printf("Redirecting to %s", redir)
			http.Redirect(w, r, redir, http.StatusPermanentRedirect)
			return
		}
		*/
	
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, err.Error())
		}
		
		var fileinfo_arr []fileinfo_internal
		for _, fileinfo := range files {
			var ext string
			if strings.Contains(fileinfo.Name(), ".") {
				exts := strings.Split(fileinfo.Name(), ".")
				ext = exts[len(exts) - 1]
			}else {
				ext = ""
			}
		
			fileinfo_to_append := fileinfo_internal {
				Name: fileinfo.Name(),
				Ext: ext,
				Dir: fileinfo.IsDir()}
			fileinfo_arr = append(fileinfo_arr, fileinfo_to_append)
		}

		topmost_dir_arr := strings.Split(pubdir, "/")
		topmost_dir := topmost_dir_arr[len(topmost_dir_arr) - 1]
		dot_dot := strings.TrimSuffix(pubdir, topmost_dir + "/")

		if r.Method == "POST" {
			r.ParseForm()
			thumbnail_string := r.FormValue("thumbnails")

			session.Put(r, "thumbnail", thumbnail_string)
			
			//fmt.Printf("thumbnail_bool: %s\n", thumbnail_bool)
			var thumbnail_bool bool
			if thumbnail_string == "false" {
				thumbnail_bool = false
			}else {
				thumbnail_bool = true
			}
			
			dirstruct := struct {
				Dirname string
				Filenames []fileinfo_internal
				Dotdot string
				Thumbnails bool
			}{	
				Dirname: pubdir,
				Filenames: fileinfo_arr,
				Dotdot: dot_dot,
				Thumbnails: thumbnail_bool}

			tmpl, err := template.ParseFiles("templates/dir.html")

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Render Error")
				return

			}

			tmpl.Execute(w, dirstruct)		
		}else { // Request is GET

			thumbnail_string := session.GetString(r, "thumbnail")

			if thumbnail_string == "" {
				thumbnail_string = "true"
			}

			var thumbnail_bool bool
			
			if thumbnail_string == "true" {
				thumbnail_bool = true
			}else {
				thumbnail_bool = false
			}
		
			dirstruct := struct {
				Dirname string
				Filenames []fileinfo_internal
				Dotdot string
				Thumbnails bool
			}{	
				Dirname: pubdir,
				Filenames: fileinfo_arr,
				Dotdot: dot_dot,
				Thumbnails: thumbnail_bool}
			

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

	}else{ // Requested content is a file
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

func main(){
	session = sessions.New(secret)
	session.Lifetime = 3 * time.Hour
	
	handler := http.NewServeMux()

	handler.HandleFunc("/", render)

	err := http.ListenAndServe(port, session.Enable(handler))

	if err != nil {
		fmt.Println(err.Error())
	}
}

