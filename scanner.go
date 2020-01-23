package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

//ScanHandler handle uploading file
type ScanHandler func(http.ResponseWriter, *http.Request) ([]string, error)

type Result struct {
	Filename string `json:"file_name"`
	Result   string `json:"result"`
}

func NewResult(out string) Result {
	fileNameWithResult := strings.Join(strings.Split(out, "_")[1:], "")
	parts := strings.Split(fileNameWithResult, ": ")
	return Result{
		Filename: parts[0],
		Result:   strings.ReplaceAll(parts[1], "\n", ""),
	}
}

func ScanVirus(filePath string) (string, error) {
	cmd := exec.Command("clamscan", "--no-summary", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err != nil {
		return out.String(), err
	}
	err = cmd.Wait()
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func (f ScanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files, err := f(w, r)
	if err == nil {
		var wg sync.WaitGroup
		resultFiles := []Result{}
		for _, f := range files {
			wg.Add(1)
			go func(file string) {
				out, _ := ScanVirus(file)
				resultFiles = append(resultFiles, NewResult(out))
				os.Remove(file)
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
