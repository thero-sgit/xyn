package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

var modes = []string{"ask", "agent"}

type mode struct {
	Icon string
	Name string
	i    int
}

func (m *mode) Toggle() {
	m.i = (m.i + 1) % len(modes)
	m.Name = modes[m.i]

	switch m.Name {
	case "ask":
		m.Icon = "?"
	case "agent":
		m.Icon = "W"
	}
}

var Config configuration

type configuration struct {
	WorkingDir       string
	SanitizedWd      string
	GroqApiKey       string
	Model 			 string
	Mode             mode
}

func getWd() string {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide a directory path. Usage: xyn <path>")
	}
	path := os.Args[1]

	if strings.Trim(path, " ") == "." {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not open directory: %s", err.Error())
		}

		path = wd
	}

	cleanPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Could not open directory: %s", err.Error())
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		log.Fatalf("Could not open directory: %s", err.Error())
	}

	if !info.IsDir() {
		log.Fatalf("%s is not a directory", cleanPath)
	}

	return cleanPath
}

func workingDirSanitized() string {
	path := Config.WorkingDir
	splitPath := strings.Split(path, "/")

	if len(splitPath) < 2 {
		return path
	}

	return strings.Join(
		[]string{
			"...",
			splitPath[len(splitPath)-2],
			splitPath[len(splitPath)-1],
		},
		"/",
	)
}

func apiKey() string {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatalf("GROQ_API_KEY is missing. Get your api key from Groq (https://console.groq.com/keys)")
	}

	return apiKey
}

func Init() {
	Config = configuration {
		WorkingDir: getWd(),
		GroqApiKey: apiKey(),
		Model:		"openai/gpt-oss-120b",
	}

	Config.SanitizedWd = workingDirSanitized()
	Config.Mode = mode{i: 0}
	Config.Mode.Toggle()
}