package main

import (
	"io"
	"os"
	"fmt"
	"io/ioutil"
	"strings"
	"net/http"
	"html/template"
)

func main(){
	// Config

	// The directory to serve
	const directory = "example/"

	// The route to start the directory from. This must end with /
	//const route = "/example/"

	// The port to run server on. You may want to make this different from 
	const port = ":8080"

	handler := http.NewServeMux()

	handler.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
		dir := strings.Replace(r.URL.Path, "/", "", 1)

		if dir == "" { dir = "/" }
		dir = directory + dir

		fileinfo, err := os.Stat(dir)

		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, fmt.Sprintf("File does not exist"))
			return
		}

		if fileinfo.IsDir(){
			//io.WriteString(w, fmt.Sprintf("%s is a directory", dir))

			pubdir := strings.TrimPrefix(dir, directory)
		
			files, err := ioutil.ReadDir(dir)
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "text/plain")
				io.WriteString(w, err.Error())
			}

			var files_arr []string

			for _, file := range files {
				var filename string
				if file.IsDir() {
					filename = file.Name() + "/"
				} else {
					filename = file.Name()
				}
				files_arr = append(files_arr, filename)
			}

			
			dirstruct := struct {
				Dirname string
				Filenames []string
			}{	
				Dirname: pubdir,
				Filenames: files_arr}

			tmpl, _ := template.ParseFiles("dir.html")
			tmpl.Execute(w, dirstruct)

			//w.Header().Set("Content-Type", "text/plain")
			//io.WriteString(w, out)

		}else{
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

		

	})

	/*handler.HandleFunc(, func (w http.ResponseWriter, r *http.Request){
		
	})*/

	err := http.ListenAndServe(port, handler)

	if err != nil {
		fmt.Println(err.Error())
	}
}

