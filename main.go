package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	_ "github.com/k0kubun/pp"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// pp.Print(r)
	
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	imagedir := path.Join(dir, "images")
	if err := os.Mkdir(imagedir, 0755); err != nil && !os.IsExist(err) {
		fmt.Fprintln(w, err)
		return
	}
	file, _, err := r.FormFile("imagedata")
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()
	basename := strconv.FormatInt(time.Now().UnixNano(), 10) + ".jpg"
	imagefile := path.Join(imagedir, basename)
	out, err := os.Create(imagefile)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// pp.Print(header)
	fmt.Fprintf(w, "http://%s/images/%s", r.Host, basename)
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	// pp.Print(r)
	imagefile := path.Join(dir, r.URL.Path)
	// pp.Print(imagefile)
	http.ServeFile(w, r, imagefile)
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/images/", imagesHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}
