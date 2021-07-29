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

			pubdir_ends := pubdir[len(pubdir) - 1]

			fmt.Printf("Dir is %s\nDir ends with: %s\n", pubdir, string(pubdir_ends))
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

			topmost_dir_arr := strings.Split(pubdir, "/")
			topmost_dir := topmost_dir_arr[len(topmost_dir_arr) - 1]
			dot_dot := strings.TrimSuffix(pubdir, topmost_dir + "/")
			
			dirstruct := struct {
				Dirname string
				Filenames []string
				Dotdot string
			}{	
				Dirname: pubdir,
				Filenames: files_arr,
				Dotdot: dot_dot}
			

			tmpl, err := template.ParseFiles("templates/dir.html")

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Render Error")
				return

			}

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

