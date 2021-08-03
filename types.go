package main

type cache_struct struct {
	Cache map[string]cached_dir
}

type cached_dir struct {
	Files []fileinfo_internal
}

type dir_struct struct {
	Dirname    string
	Filenames  []fileinfo_internal
	Dotdot     string
	Thumbnails bool
	Path       string
	Header     header
	Footer     footer
}

type server struct {
	Dir     string
	Port    string
	Favicon string
}

type link struct {
	Href string
	Text string
}

/*type text struct {
	Text string
	Font string
	Size int
}*/

type footer struct {
	Show bool
	Text string
	Link link
}

type header struct {
	Show    bool
	Text    string
	Subtext string
}

type cache_type struct {
	Enable bool
	//In_Memory bool
}

type caching struct {
	Filenames cache_type
}

type config_format struct {
	Server  server
	Header  header
	Footer  footer
	Caching caching
}

type fileinfo_internal struct {
	Name string
	Ext  string
	Dir  bool
	Type string
}
