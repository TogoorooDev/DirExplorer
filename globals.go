package main

import (
	"sync"

	"github.com/golangcollege/sessions"
)

// Globals:

// Config
var config config_format

// Session
var session *sessions.Session
var secret = []byte("u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4")

// Cacheing
var filenames_cache cache_struct

// Mutex
var config_lock sync.Mutex

// Type extention arrays (fileinfo_internal.Type)
var image_ext = []string{"png", "jpg", "jpeg", "webp", "gif", "svg", "avif", "apng", "bmp", "ico", "jp2", "j2k", "jpf", "jpm", "jpg2", "j2c", "jpc", "jpx", "mj2"}
var video_ext = []string{"3gp", "3gpp", "mpeg", "mp4", "webm", "mkv", "ogv", "mov"}
var audio_ext = []string{"aac", "mp3", "flac", "wav", "alac", "aiff", "dsd", "pcm", "l16", "au"}
var text_ext = []string{"txt", "log", "doc", "docx", "pdf", "odt", "rtf"}
var binary_ext = []string{"exe", "appimage", "app"}
