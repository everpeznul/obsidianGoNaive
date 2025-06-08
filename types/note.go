package types

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type note struct {
	title       string
	text        string
}

func parseNote(filePath string) note {

	file := filepath.Base(filePath)
	fileName := strings.TrimSuffix(file, filepath.Ext(file))

	title := getTitle(fileName)
	text := getText(filePath)

	newNote := note{title, text}

	connections := fmt.Sprintf("[fouder:: [[%s]]]\n[ancestor:: [[%s]]]\n[father:: [[%s]]]", newNote.findFounder(//))

	return note{title, text}
}

func getTitle(fileName string) string {

	parts := strings.Split(fileName, ".")
	name := strings.Split(parts[len(parts)-1], "_")

	runes := []rune(name[0])
	runes[0] = unicode.ToUpper(runes[0])
	name[0] = string(runes)

	title := strings.Join(name, " ")
	return title
}

func getText(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	body := string(content)
	re := regexp.MustCompile(`(?m)^# .*$`)
	parts := re.Split(body, -1)

	if len(parts) > 0 {
		text := parts[len(parts)-1]

		return text
	}
	return body
}

func (n *note) find(files []string, fileName string) string {

	for _, tempFile := range files {

		tempFileName := strings.TrimSuffix(tempFile, filepath.Ext(tempFile))
		if strings.EqualFold(tempFileName, fileName) {
			return tempFileName
		}
	}
	return ""
}

func (n *note) findFather(files []string, fileName string) string {

	fileNameParts := strings.Split(n.title, ".")
	fatherName := strings.Join(fileNameParts[:len(fileNameParts)-1], ".")
	return n.find(files, fatherName)
}

func (n *note) findAncestor(files []string, fileName string) string {

	fileNameParts := strings.Split(n.title, ".")
	fatherName := strings.Join(fileNameParts[:len(fileNameParts)-2], ".")
	return n.find(files, fatherName)
}

func (n *note) findFounder(files []string, fileName string) string {

	fileNameParts := strings.Split(n.title, ".")
	fatherName := fileNameParts[0]
	return n.find(files, fatherName)
}
