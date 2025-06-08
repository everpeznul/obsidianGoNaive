package scripts

func update() {

	files, _ := collectFiles("./")
	groups := splitFilesByGroup(files)
}
