package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func getDistroName() string {
	if runtime.GOOS != "linux" {
		return runtime.GOOS
	}

	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "Unknown Linux"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[13:], `"`)
		}
	}

	return "Unknown Linux"
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		osName := getDistroName()
		fmt.Fprintf(w, "Hello from %s!\n", osName)
	})

	fmt.Println("Server is running on http://localhost:5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}

