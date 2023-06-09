package project

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Project struct {
	Commands           []string `yaml:"commands"`
	ExcludeExtensions  []string `yaml:"exclude_extensions"`
	ExcludeDirectories []string `yaml:"exclude_directories"`
}

var defaultProject = Project{
	Commands: []string{"echo 'Hello World!'"},
	ExcludeExtensions: []string{
		".exe",
	},
	ExcludeDirectories: []string{
		".git",
	},
}

// ParseProject parses the .anymon.yaml file in the current working directory
func ParseProject() (Project, error) {
	// get current working directory
	path, err := os.Getwd()
	if err != nil {
		return Project{}, err
	}

	// open config file
	file, err := os.Open(path + "/.anymon.yaml")
	defer file.Close()
	if err != nil {
		return Project{}, err
	}

	var project Project
	err = yaml.NewDecoder(file).Decode(&project)
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

// CreateProject creates a new empty config file, named .anymon.yaml
// in the current working directory
func CreateProject() error {
	file, err := os.Create(".anymon.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	err = yaml.NewEncoder(file).Encode(defaultProject)
	if err != nil {
		return err
	}

	return nil
}

func Watch(project Project) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer watcher.Close()

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = watcher.Add(wd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// add all subdirectories
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}

		return nil
	})

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			clear := exec.Command("cmd", "/c", "cls")
			clear.Stdout = os.Stdout
			clear.Run()

			for _, dir := range project.ExcludeDirectories {
				if strings.Contains(event.Name, dir) {
					continue
				}
			}

			for _, ext := range project.ExcludeExtensions {
				if event.Name[len(event.Name)-len(ext):] == ext {
					continue
				}
			}

			for _, command := range project.Commands {
				fmt.Printf("Executing: %s\n", command)

				if c, err := exec.Command("cmd", "/c", command).CombinedOutput(); err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println(string(c))
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
	}
}
