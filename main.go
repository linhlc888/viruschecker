package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/dutchcoders/go-clamd"
)

var form *template.Template

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		printVersion()
		return
	}

	form = template.Must(template.ParseFiles("./files/upload.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			form.Execute(w, nil)
		} else {
			w.WriteHeader(405)
			w.Write([]byte("Method not allowed"))
		}
	})
	s := fmt.Sprintf("tcp://%s:%d", *clamdHost, *clamdPort)
	fmt.Println("ClamAV server: " + s)
	scanHander := ScanHandler{
		uploadHandler: upload,
		clamd:         clamd.NewClamd(s),
	}
	http.Handle("/api/v1/scan", scanHander)
	fmt.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}
