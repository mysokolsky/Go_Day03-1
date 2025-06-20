package main

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // Linux и WSL
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Start()
	if err != nil {
		log.Printf("Не удалось открыть браузер: %v", err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет! Сервер работает."))
	})

	go openBrowser("http://localhost:8888") // Открываем браузер в отдельной горутине

	log.Println("Сервер запущен на http://localhost:8888")
	
	err := http.ListenAndServe(":8888", nil)
	
	if err != nil {
		log.Fatal(err)
	}
}