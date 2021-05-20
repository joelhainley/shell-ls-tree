package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const ITEM_GLYPH = "\u251C\u2500"
const LAST_ITEM_GLYPH = "\u2514\u2500"
const LINE_CONT_GLYPH = "\u2502"

func process(path string, maxDepth int, curDepth int, isLastItemInParent bool, prefix[] string) {
	files, err := ioutil.ReadDir(path)

	check(err)

	itemGlyph := ITEM_GLYPH
	isLastItem := false

	newPrefix := prefix

	for i, file := range files {
		// not dealing with symlinks at this time
		if file.Mode()&os.ModeSymlink != 0 {
			continue
		}

		if i == len(files) - 1 {
			isLastItem = true
			itemGlyph = LAST_ITEM_GLYPH
		}

		prefixSlug := getPrefixSlug(newPrefix)

		if file.IsDir() {
			// TOOD figure out better mechanism for creating paths
			newPath := fmt.Sprintf("%s/%s", path, file.Name())

			format := "%s%s %s\n"

			var nextPrefix[] string

			// TODO logic could be cleaner
			if isLastItemInParent && isLastItem {
				nextPrefix = append(newPrefix, "")
			} else if isLastItem {
				nextPrefix = append(newPrefix, "")
			} else {
				nextPrefix = append(newPrefix, LINE_CONT_GLYPH)
			}

			processNextDir := true

			_, err = ioutil.ReadDir(newPath)

			filename := file.Name()

			if err != nil && os.IsPermission(err) {
				var sb strings.Builder
				sb.WriteString(filename)
				sb.WriteString( " (directory - permission denied)")
				filename = sb.String()
				processNextDir = false
			}

			fmt.Printf(format, prefixSlug, itemGlyph, filename)

			if processNextDir {
				process(newPath, maxDepth, curDepth+1, isLastItem, nextPrefix)
			}
		} else {
			fmt.Printf("%s%s %s\n", prefixSlug, itemGlyph, file.Name())
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getPrefixSlug(prefix[] string) string {
	var sb strings.Builder

	for i := 0; i < len(prefix); i++ {
		sb.WriteString(prefix[i])
		sb.WriteString("   ")
	}

	return sb.String()
}


func main() {
	var prefix[] string

	args := os.Args[1:] // omit the application org
	searchPath := args[0]
	maxDepth := -1
	curDepth := 0
	lastItemInParent := false

	process(searchPath, maxDepth, curDepth, lastItemInParent, prefix)
}
