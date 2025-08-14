package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func browserAuth(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		// err = exec.Command("xdg-open", url).Start()
		err = exec.Command("cmd.exe", "/c", "start", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatalf("could not open browser: %v", err)
	}

	fmt.Println("Please log into Spotify by visiting this URL in your browser:")
	fmt.Println(url)
}
