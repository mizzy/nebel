package nebel

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var watchDirs = []string{"posts", "layouts", "static"}

// Serve generates the site, serves "public" over HTTP, and regenerates
// whenever a file under posts, layouts or static changes.
func Serve(addr string) error {
	if err := Generate(); err != nil {
		return err
	}

	go watch(500 * time.Millisecond)

	log.Printf("Serving public/ at http://%s/", addr)
	return http.ListenAndServe(addr, http.FileServer(http.Dir("public")))
}

// ponytail: polls mtimes instead of fsnotify; switch if the site grows enough for walking to be slow.
func watch(interval time.Duration) {
	last := snapshot(watchDirs)
	for range time.Tick(interval) {
		cur := snapshot(watchDirs)
		if cur == last {
			continue
		}
		last = cur
		start := time.Now()
		if err := Generate(); err != nil {
			log.Printf("generate failed: %v", err)
			continue
		}
		log.Printf("regenerated in %v", time.Since(start).Round(time.Millisecond))
	}
}

// snapshot returns a string that changes when any file under dirs is
// added, removed, or modified.
func snapshot(dirs []string) string {
	var b strings.Builder
	for _, dir := range dirs {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			fmt.Fprintf(&b, "%s %d %d\n", path, info.Size(), info.ModTime().UnixNano())
			return nil
		})
	}
	return b.String()
}
