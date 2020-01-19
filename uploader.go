package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func upload(w http.ResponseWriter, r *http.Request) ([]string, error) {
	if r.Method != http.MethodPost {
		return nil, NewHTTPError(nil, 405, "Method not allowed", "Allow only POST method")
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart") {
		return nil, NewHTTPError(nil, 400, "Upload error", "Content-Type is not multipart/formdata")
	}
	err := r.ParseMultipartForm(16 << 20)
	if err != nil {
		return nil, NewHTTPError(err, 400, "File too large", "Max file size is 16M")
	}
	data := r.MultipartForm
	files := data.File["files"]
	if len(files) == 0 {
		return nil, NewHTTPError(err, 400, "Upload error", "No file to upload")
	}

	if len(files) > 3 {
		return nil, NewHTTPError(err, 400, "Upload error", "Upload too many files. Max: 3")
	}

	res := []string{}
	for _, fh := range files {

		filePath, err := copy(fh)
		if err != nil {
			return nil, NewHTTPError(err, 400, "File error", "Cannot upload file")
		}
		res = append(res, filePath)
	}
	return res, nil
}

func copy(fh *multipart.FileHeader) (string, error) {
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	filePath := fmt.Sprintf("./tmp/%s_%s", RandomString(4), fh.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, f)
	if err != nil {
		return "", err
	}

	return filePath, nil

}
