package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/uploader.html")
	if err != nil {
		log.Fatal("tamplate error")
	}
	t.Execute(w, struct{}{})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
	file, fileHeader, err := r.FormFile("csv")
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "ファイルを開くことができません")
		return
	}

	filename := fileHeader.Filename
	ext := filepath.Ext(filename)

	w.Header().Add("Content-Type", "text/csv")
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s-shiftjis%s", strings.TrimSuffix(filename, ext), ext))
	writer := transform.NewWriter(w, japanese.ShiftJIS.NewEncoder())

	tee := io.TeeReader(file, writer)
	s := bufio.NewScanner(tee)
	for s.Scan() {
	}
	if err := s.Err(); err != nil {
		log.Println(err)
	}
	log.Println("done")
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", handler)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":"+port, nil)
}
