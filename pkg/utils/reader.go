package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func ReadFile(target string) ([]byte, error) {
	content, err := os.ReadFile(target)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}
	return content, nil
}

func ReadFromURL(target string) ([]byte, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %v", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read URL response: %v", err)
	}
	return content, nil
}
