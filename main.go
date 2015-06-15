package main

import (
	"flag"
	"fmt"
	_ "github.com/k0kubun/pp"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

var portNumber = flag.String("port", "8080", "port number.")
var basicAuthUser = flag.String("user", "", "basic auth user name")
var basicAuthPass = flag.String("pass", "", "basic auth user pass")

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

func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	if *basicAuthUser == "" || *basicAuthPass == "" {
		return true
	}

	username, password, ok := r.BasicAuth()
	// log.Println(username, password, ok)
	if ok == false {
		return false
	}
	return username == *basicAuthUser && password == *basicAuthPass
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	if checkAuth(w, r) == false {
		w.Header().Set("WWW-Authenticate", `Basic realm="Atto"`)
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}

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
	flag.Parse()
	if *basicAuthUser != "" && *basicAuthPass != "" {
		log.Println("basic auth: " + *basicAuthUser)
	}
	log.Println("listen:" + *portNumber)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/images/", imagesHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":"+*portNumber, nil)
}
