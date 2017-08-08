package main

import (
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/uiacode/web/moustache"
)

var uploadTemplate = template.Must(template.ParseFiles("templates/upload.html"))
var errorTemplate = template.Must(template.ParseFiles("templates/error.html"))
var editTemplate = template.Must(template.ParseFiles("templates/edit.html"))

func edit(w http.ResponseWriter, r *http.Request) {
	editTemplate.Execute(w, r.FormValue("id"))
}

func img(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("image-" + r.FormValue("id"))
	check(err)
	m, _, err := image.Decode(f)
	check(err)

	x, _ := strconv.Atoi(r.FormValue("x"))
	log.Printf("x=%d", x)
	y, _ := strconv.Atoi(r.FormValue("y"))
	log.Printf("y=%d", y)
	s, _ := strconv.Atoi(r.FormValue("s"))
	d, _ := strconv.Atoi(r.FormValue("d"))
	m = moustache.Moustache(m, x, y, s, d)

	w.Header().Set("Content-type", "image/jpeg")
	jpeg.Encode(w, m, nil) // Default JPEG options
}

func errorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				w.WriteHeader(500)
				errorTemplate.Execute(w, e)
				log.Println(e)
			}
		}()
		fn(w, r)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	//r.Header.Add("Content-Type", "multipart/form-data")

	fmt.Fprint(w, "Her skal det bli FAQ for IS-105.")
}

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, "image-"+r.FormValue("id"))
}

func upload(w http.ResponseWriter, r *http.Request) {
	//r.Header.Add("Content-Type", "multipart/form-data")
	if r.Method != "POST" {
		uploadTemplate.Execute(w, nil)
		return
	}
	f, _, err := r.FormFile("image")
	check(err)
	defer f.Close()

	t, err := ioutil.TempFile(".", "image-")
	check(err)
	defer t.Close()

	_, err = io.Copy(t, f)
	check(err)
	http.Redirect(w, r, "/edit?id="+t.Name()[6:], 302)
}

func main() {
	http.HandleFunc("/", errorHandler(upload))
	http.HandleFunc("/view", errorHandler(view))
	http.HandleFunc("/edit", errorHandler(edit))
	http.HandleFunc("/img", errorHandler(img))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.ListenAndServe(":8008", nil)
}
