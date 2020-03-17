package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	optDir  string
	optBind string
)

func main() {
	flag.StringVar(&optDir, "dir", ".", "data directory")
	flag.StringVar(&optBind, "bind", ":9910", "address and port to bind")
	flag.Parse()

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = rw.Write([]byte("method not allowed"))
			return
		}
		relPath := filepath.Join(optDir, req.URL.Path[1:])
		dirPath := filepath.Dir(relPath)
		var err error
		if err = os.MkdirAll(dirPath, 0777); err != nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		log.Printf("RelPath: %s", relPath)
		log.Printf("DirPath: %s", dirPath)
		var file *os.File
		if file, err = os.OpenFile(relPath, os.O_RDWR|os.O_CREATE, 0777); err != nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
		defer file.Close()
		io.Copy(file, req.Body)
	})
	log.Fatal(http.ListenAndServe(optBind, nil))
}
