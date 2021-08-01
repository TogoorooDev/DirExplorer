package main

import (
	"encoding/json"
	"errors"
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
	var dir_cache cache_struct
	err := filepath.Walk(config.Server.Dir, func(pathX string, infoX os.FileInfo, errX error) error {
		if errX != nil {
			return errX
		}

		path := strings.TrimPrefix(pathX, config.Server.Dir)

		if path == "" {
			path = "/"
		}

		if infoX.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	//var files_arr []fileinfo_internal

	for _, dir := range dirs {
		var container cached_dir

		true_dir := filepath.Join(config.Server.Dir, dir)
		files, err := ioutil.ReadDir(true_dir)

		if err != nil {
			return err
		}

		files_chan := make(chan fileinfo_internal)
		var count int

		for _, file := range files {
			go index_file(file, files_chan)
			count++
		}

		for count > 0 {
			//fmt.Println(count)
			container.Files = append(container.Files, <-files_chan)
			count--
		}

		if this_cache, ok := dir_cache.Cache[dir]; ok {
			this_cache = container
			dir_cache.Cache[dir] = this_cache
		}

		//dir_cache.Cache[dir].Files = append(dir_cache.Cache[dir].Files, container)
	}

	filenames_cache = dir_cache

	out, err := json.Marshal(dir_cache)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("index.json", out, 0644)
	if err != nil {
		return err
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
	if dir == "" {
		dir = "/"
	}
	if config.Caching.Filenames.Enable {
		ret, ok := filenames_cache.Cache[dir]
		if ok {
			return ret.Files, nil
		} else {
			return nil, errors.New("Directory not found")
		}
	} else {
		return index_dir(dir)
	}
}
