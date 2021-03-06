package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func favicon(w http.ResponseWriter, r *http.Request) {
	//io.Stat(config.Server.Favicon)
	http.ServeFile(w, r, config.Server.Favicon)
}

/*func search(w http.RepsonseWriter, r *http.Request) {
	
        }*/

func render(w http.ResponseWriter, r *http.Request) {
	//slashless_path := config.Path[1:]
	dir := strings.Replace(r.URL.Path, "/", "", 1)
	//dir = strings.TrimPrefix(dir, slashless_path)
	pubdir := strings.TrimPrefix(dir, config.Server.Dir)

	//fmt.Printf("dir: %s\npubdir %s\npath %s\n", dir, pubdir, config.Path)

	if dir == "" {
		dir = "/"
	}
	dir = filepath.Join(config.Server.Dir, dir)

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

		fileinfo_arr, err := load_dir(dir)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, err.Error())
		}

		topmost_dir_arr := strings.Split(pubdir, "/")
		topmost_dir := topmost_dir_arr[len(topmost_dir_arr)-1]
		dot_dot := strings.TrimSuffix(pubdir, topmost_dir+"/")
		
		if r.Method == "POST" {
			r.ParseForm()
			form_type := r.FormValue("type")
			if form_type == "thumbnailstat" {
				fmt.Println("Thumbnail updating!1!")
				fmt.Println(r.FormValue("thumbnails"))
				thumbnail_bool := r.FormValue("thumbnails") == "true"
				session.Put(r, "thumbnail", thumbnail_bool)

			}else if form_type == "search" {
				
			}

			//fmt.Printf("thumbnail_bool: %s\n", thumbnail_bool)
			

			/* dirstruct := dir_struct{
				Dirname:    pubdir,
				Filenames:  fileinfo_arr,
				Dotdot:     dot_dot,
				Thumbnails: thumbnail_bool,
				Path:       "/",
				Header:     config.Header,
				Footer:     config.Footer,
			}

			tmpl, err := template.ParseFiles("templates/dir.html")

			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Render Error")
				return

			}

			tmpl.Execute(w, dirstruct) */
		} else { // Request is GET

			
		}
		
		thumbnail_bool := session.GetBool(r, "thumbnail")
		
		
		var dirstruct dir_struct
		dirstruct = dir_struct{
			Dirname:    pubdir,
			Filenames:  fileinfo_arr,
			Dotdot:     dot_dot,
			Thumbnails: thumbnail_bool,
			Path:       "/",
			Header:     config.Header,
			Footer:     config.Footer,
		}

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
