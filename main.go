package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func processFile(path string, pattern *regexp.Regexp, info os.FileInfo, ext string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	newData := pattern.ReplaceAllStringFunc(string(data), func(match string) string {
		key := match[2 : len(match)-1]
		value := os.Getenv(key)
		if value == "" {
			return match
		}
		return value
	})

	err = os.WriteFile(path, []byte(newData), info.Mode())
	if err != nil {
		return err
	}

	fmt.Printf("Replaced values in %s\n", path)

	return nil
}

func main() {
	path := os.Getenv("PATH")
	if path == "" {
		path = "k8s"
	}

	pattern := regexp.MustCompile(`\$\{(\w+)\}`)

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if info.IsDir() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".yaml" && filepath.Ext(path) != ".json" {
				return nil
			}

			return processFile(path, pattern, info, filepath.Ext(path))
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".json" {
			err = processFile(path, pattern, info, filepath.Ext(path))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
