package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

const asciiStart = 32
const charHeight = 8

// ================= LOAD BANNER =================
func loadBanner(filename string) (map[rune][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	banner := make(map[rune][]string)

	var block []string
	char := rune(asciiStart)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			if len(block) == charHeight {
				banner[char] = block
				char++
			}
			block = []string{}
		} else {
			block = append(block, line)
		}
	}

	return banner, nil
}

// ================= GENERATE ASCII =================
func generateASCII(input string, banner map[rune][]string) string {
	lines := strings.Split(input, "\n")
	var result string

	for _, textLine := range lines {

		output := make([]string, charHeight)

		for _, char := range textLine {
			asciiChar, ok := banner[char]
			if !ok {
				asciiChar = banner[' ']
			}

			for i := 0; i < charHeight; i++ {
				output[i] += asciiChar[i]
			}
		}

		for _, line := range output {
			result += line + "\n"
		}
	}

	return result
}

// ================= HANDLERS =================

// GET /
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// POST /ascii-art
func asciiHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	bannerName := r.FormValue("banner")

	if text == "" || bannerName == "" {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}

	bannerFile := bannerName + ".txt"

	banner, err := loadBanner(bannerFile)
	if err != nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	result := generateASCII(text, banner)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Result string
	}{
		Result: result,
	}

	tmpl.Execute(w, data)
}

// ================= MAIN =================
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiHandler)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
