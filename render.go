package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

func index_file(file os.FileInfo, out chan fileinfo_internal) {
	var ext string
	if strings.Contains(file.Name(), ".") {
		exts := strings.Split(file.Name(), ".")
		ext = exts[len(exts)-1]
	} else {
		ext = ""
	}

	fileinfo_to_append := fileinfo_internal{
		Name: file.Name(),
		Ext:  ext,
		Dir:  file.IsDir()}

	out <- fileinfo_to_append
}

func render(w http.ResponseWriter, r *http.Request) {
	//slashless_path := config.Path[1:]
	dir := strings.Replace(r.URL.Path, "/", "", 1)
	//dir = strings.TrimPrefix(dir, slashless_path)
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

		file_chan := make(chan fileinfo_internal)

		count := len(files)

		for _, fileinfo := range files {
			go index_file(fileinfo, file_chan)
		}

		for count > 0 {
			fileinfo_arr = append(fileinfo_arr, <-file_chan)
			count--
		}

		//Split into two arrays
		for _, fileinfo := range fileinfo_arr {
			if fileinfo.Dir {
				fileinfo_dir_arr = append(fileinfo_dir_arr, fileinfo)
			} else {
				fileinfo_file_arr = append(fileinfo_file_arr, fileinfo)
			}
		}

		sort.Slice(fileinfo_dir_arr, func(i, j int) bool {
			return fileinfo_dir_arr[i].Name < fileinfo_dir_arr[j].Name
		})

		sort.Slice(fileinfo_file_arr, func(i, j int) bool {
			return fileinfo_file_arr[i].Name < fileinfo_file_arr[j].Name
		})

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
				Path:       "/"}

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
				Path:       "/"}

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
