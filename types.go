package main

type dir_struct struct {
	Dirname    string
	Filenames  []fileinfo_internal
	Dotdot     string
	Thumbnails bool
	Path       string
}

type config_format struct {
	Dir  string
	Port string
}

type fileinfo_internal struct {
	Name string
	Ext  string
	Dir  bool
}
