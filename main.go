package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func getOSName() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "Unknown OS"
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			return strings.Trim(line[13:], "\"")
		}
	}
	return "Unknown OS"
}

func getCPUUsage() float64 {
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) == 0 {
		return 0.0
	}
	return percent[0]
}

func getMemoryUsage() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0.0
	}
	return v.UsedPercent
}

func handler(w http.ResponseWriter, r *http.Request) {
	osName := getOSName()
	cpuUsage := getCPUUsage()
	memUsage := getMemoryUsage()
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>System Info</title>
	<style>
		body {
			background-color: #f4f4f4;
			font-family: sans-serif;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
		}
		.container {
			background: white;
			padding: 30px 40px;
			border-radius: 12px;
			box-shadow: 0 4px 12px rgba(0,0,0,0.1);
			text-align: center;
		}
		h1 {
			margin-bottom: 20px;
		}
		p {
			margin: 8px 0;
			font-size: 1.1rem;
		}
	</style>
</head>
<body>
	<div class="container">
		<h1>📟 System Info</h1>
		<p><strong>OS:</strong> %s</p>
		<p><strong>CPU Usage:</strong> %.2f%%</p>
		<p><strong>Memory Usage:</strong> %.2f%%</p>
		<p><strong>Datetime:</strong> %s</p>
	</div>
</body>
</html>
`, osName, cpuUsage, memUsage, currentTime)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running on http://localhost:5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		panic(err)
	}
}

