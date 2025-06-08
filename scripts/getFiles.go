package scripts

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func isNoteType(path string) bool {

	ext := filepath.Ext(path)
	return ext == ".md"

}

func collectFiles(dir string) ([]string, error) {

	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && isNoteType(path) {
			files.append(files, path)
		}

		return nil
	})

	return files, err
}

func splitFilesByGroup(files []files) map[string][]string {

	var groups map[string][]string = make(map[string][]string)

	for _, file := range files {

		ancestor := getAncestor(filepath.Base(file))
		groups[ancestor] = append(groups[ancestor], file)

	}

	return groups
}

func getAncestor(fileName string) string {

	nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	parts := strings.Split(nameWithoutExt, ".")

	if len(parts) > 0 && parts[0] != "" {
		return strings.ToLower(strings.TrimSpace(parts[0]))
	}

	return "misc"
}

// func splitGroups
