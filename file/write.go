package file

import (
	"log"
	"os"
)

// CreateProject initializes the project directory and returns the path to the project
func CreateProject(projectName string) string {
	path := currentDirectory()

	// define project path
	projectPath := path + "/" + projectName

	// create root directory
	err := os.MkdirAll(projectPath, 0777)
	check(err)

	// Create CS/JS/img directories
	createCSS(projectPath)
	createJS(projectPath)
	createIMG(projectPath)

	_, err = os.Create(projectPath + "/" + "index.html")
	check(err)
	// project path
	return projectPath
}

func currentDirectory() string {
	path, err := os.Getwd()
	check(err)
	return path
}

func createCSS(path string) {
	// create css directory
	err := os.MkdirAll(path+"/"+"css", 0777)
	check(err)
}
func createJS(path string) {
	err := os.MkdirAll(path+"/"+"js", 0777)
	check(err)
}

func createIMG(path string) {
	err := os.MkdirAll(path+"/"+"imgs", 0777)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
