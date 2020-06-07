package file


import (
	"testing"
	"os"
)

func TestCreateProject(t *testing.T) {
	// try creating necessary directories for project setup
	projectPath := CreateProject()
	// get current working directory
	cwd, err := os.Getwd()

	if cwd != projectPath {
		t.Error("project path is not current working directory")
	}
}



