package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var Config configuration

type configuration struct {
	WorkingDir       string
	SanitizedWd      string
	GroqClientConfig openai.ClientConfig
	Model 			 string
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
		// log.Fatalf("GROQ_API_KEY is missing. Get your api key from Groq (https://console.groq.com/keys)")
		apiKey = "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"
	}

	return apiKey
}

func getGroqClientConfig() openai.ClientConfig {
	config := openai.DefaultConfig(apiKey())
	config.BaseURL = "https://api.groq.com/openai/v1"

	return config
}

func Init() {
	Config = configuration {
		WorkingDir:       getWd(),
		GroqClientConfig: getGroqClientConfig(),
		Model:			  "openai/gpt-oss-120b",
	}

	Config.SanitizedWd = workingDirSanitized()
}