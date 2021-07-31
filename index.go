package main

import (
	// "encoding/json"
	// "errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	//"github.com/fsnotify/fsnotify"
)

//func init_notify_

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
		Dir:  file.IsDir(),
	}

	out <- fileinfo_to_append
}

func cache() error {
	var dirs []string
	var dir_cache cached_dir
	err := filepath.Walk(config.Server.Dir, func(pathX string, infoX os.FileInfo, errX error) error {
		if errX != nil {
			return errX
		}

		path := strings.TrimPrefix(pathX, config.Server.Dir)

		// if path == "" {
		// path = "/"
		// }

		if infoX.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})
	if err != nil {
		return err
	}
	fmt.Println(dirs)

	for _, dir := range dirs {

	}

	return nil
}

func index_dir(dir string) ([]fileinfo_internal, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
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

	return fileinfo_arr, nil

}

func load_dir(dir string) ([]fileinfo_internal, error) {
	cache()
	return index_dir(dir)
}
