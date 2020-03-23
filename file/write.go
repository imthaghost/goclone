package file

import (
	"io/ioutil"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Println(err)
	}
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

// WriteStream will write a stream to desierd file
func WriteStream(filename string, stream []byte) {
	// wite string to file
	ioutil.WriteFile(filename, []byte(stream), 0777)
	// dst, err := os.Create(filepath.Join(dir, filepath.Base(file.Filename))) // dir is directory where you want to save file.
	// if err != nil {
	// 	checkErr(err)
	// 	return
	// }
	// defer dst.Close()
	// if _, err = io.Copy(dst, src); err != nil {
	// 	checkErr(err)
	// 	return
	// }

}

// CreateProject initializes the project directory and returns the path to the project
func CreateProject(projectName string) string {
	path := currentDirectory()
	// define project path
	projectPath := path + "/" + projectName
	// create root directory
	err := os.MkdirAll(projectPath, 0777)
	check(err)
	createCSS(projectPath)
	createJS(projectPath)
	createIMG(projectPath)
	f, err := os.Create(projectPath + "/" + "index.html")
	check(err)
	f.Close()
	// project path
	return projectPath
}
