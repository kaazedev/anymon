package project

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"os"
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

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return

			}

			fmt.Println("Event:", event)

			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("File modified")
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				fmt.Println("File created")
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Error:", err)
		}
	}
}
