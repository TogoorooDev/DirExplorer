package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml/v2"
)

func watch_config() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error initializing  watcher. ")
	}
	defer watcher.Close()

	fmt.Println("Starting watcher")

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//fmt.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//fmt.Println("modified file: ", event.Name)
					if strings.HasSuffix(event.Name, "config.toml") {
						parse_config()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println(err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		fmt.Printf("Cannot watch config.toml. The program will still run, just without automatic config updates.")
	}
	<-done
}

func parse_config() {
	config_lock.Lock()

	config_file_handle, err := os.Open("config.toml")
	if err != nil {
		fmt.Printf("Cannot open config.toml: %s", err.Error())
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
		fmt.Println("FATAL: Directory not specified in config.toml")
	}

	if config.Caching.Filenames.Enable {
		err := cache()
		if err != nil {
			config.Caching.Filenames.Enable = false
		}
	}

	config_lock.Unlock()
}
