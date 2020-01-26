package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/dutchcoders/go-clamd"
)

//ScanHandler handle uploading file
type ScanHandler struct {
	uploadHandler func(http.ResponseWriter, *http.Request) ([]string, error)
	clamd         *clamd.Clamd
}

type Result struct {
	Filename    string `json:"file_name"`
	Result      string `json:"result"`
	Description string `json:"description"`
}

func (s ScanHandler) scanVirus(filePath string) Result {
	uploadedFileName := strings.Join(strings.Split(filePath, "_")[1:], "")
	f, _ := os.Open(filePath)
	resCh, err := s.clamd.ScanStream(f, make(chan bool))
	if err != nil {
		fmt.Printf("%v\n", err)
		return Result{
			Filename:    uploadedFileName,
			Result:      "Internal server error. Cannot scan uploaded file.",
			Description: "Maybe file size > 25M.",
		}
	}
	res, ok := <-resCh
	if ok {
		return Result{
			Filename:    uploadedFileName,
			Result:      res.Status,
			Description: res.Description,
		}
	}
	return Result{
		Filename: uploadedFileName,
		Result:   "Internal server error. Cannot scan uploaded file.",
	}
}

func (s ScanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files, err := s.uploadHandler(w, r)
	if err == nil {
		var wg sync.WaitGroup
		resultFiles := []Result{}
		for _, f := range files {
			wg.Add(1)
			go func(file string) {
				res := s.scanVirus(file)
				defer os.Remove(file)
				resultFiles = append(resultFiles, res)
				wg.Done()
			}(f)
		}
		wg.Wait()
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		m := make(map[string]interface{})
		m["results"] = resultFiles
		json, _ := json.Marshal(m)
		w.Write(json)
		return
	}
	log.Printf("An error occured: %v", err)
	clientErr, ok := err.(ClientError)
	if !ok {
		w.WriteHeader(500)
		return
	}
	body, err := clientErr.ResponseBody()
	if err != nil {
		log.Printf("An error occured: %v", err)
		w.WriteHeader(500)
		return
	}
	status, headers := clientErr.ResponseHeaders()
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
	w.Write(body)
}
