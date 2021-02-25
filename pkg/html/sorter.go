package html

// LinkRestructure grabs all html files in project directory
// reorganizes each file with local links (css js images)
func LinkRestructure(projectDir string) {
	// Redirect JS/CSS/Img tags to the correct place :)
	arrange(projectDir)
}
