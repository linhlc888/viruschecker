package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var form *template.Template

func main() {
	form = template.Must(template.ParseFiles("./files/upload.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			form.Execute(w, nil)
		} else {
			w.WriteHeader(405)
			w.Write([]byte("Method not allowed"))
		}
	})
	http.Handle("/api/v1/scan", ScanHandler(upload))
	fmt.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
