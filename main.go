package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	// Get the directory to scan for YAML files.
	dir := os.Getenv("DIR")
	if dir == "" {
		// If the environment variable is not set, use the k8s directory.
		dir = "k8s"
	}

	// Define the pattern to search for in the YAML files.
	pattern := regexp.MustCompile(`\$\{(\w+)\}`)

	// Scan the directory for YAML files.
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".yaml" {
			return nil
		}

		// Read the contents of the YAML file.
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Replace the pattern in the YAML file with the environment variable value.
		newData := pattern.ReplaceAllStringFunc(string(data), func(match string) string {
			key := match[2 : len(match)-1]
			value := os.Getenv(key)
			if value == "" {
				// If the environment variable is not set, use the original value.
				return match
			}
			return value
		})

		// Write the new contents to the YAML file.
		err = os.WriteFile(path, []byte(newData), info.Mode())
		if err != nil {
			return err
		}

		fmt.Printf("Replaced values in %s\n", path)

		return nil
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
