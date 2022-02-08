package file

import (
	"log"
	"os"
)

// CreateProject initializes the project directory and returns the path to the project
// TODO make function more modular to obtain different html files
func CreateProject(projectName string) string {
	// current workin directory
	path := currentDirectory()

	// define project path
	projectPath := path + "/" + projectName

	// create base directory
	err := os.MkdirAll(projectPath, 0777)
	check(err)

	// main inedx file
	_, err = os.Create(projectPath + "/" + "index.html")
	check(err)
	// project path
	return projectPath
}

// currentDirectory get the current working directory
func currentDirectory() string {
	path, err := os.Getwd()
	check(err)
	return path
}

// createCSS create a css directory in the current path
func createCSS(path string) {
	// create css directory
	err := os.MkdirAll(path+"/"+"css", 0777)
	check(err)
}

// createJS create a JS directory in the current path
func createJS(path string) {
	err := os.MkdirAll(path+"/"+"js", 0777)
	check(err)
}

// createIMG create a image directory in the current path
func createIMG(path string) {
	err := os.MkdirAll(path+"/"+"imgs", 0777)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
