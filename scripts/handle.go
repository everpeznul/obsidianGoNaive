package scripts

import (
	"obsidianGoNaive/types"
)

func handleGroups(groups map[string][]string) {

	for _, group := range groups {

		handleGroup(group)
	}
}

func handleGroup(group []string) {

	for _, file := range group {

		handleFile(file)
	}
}

func handleFile(file string) {

}
