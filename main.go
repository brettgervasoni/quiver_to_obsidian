package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var AllPaths []string
var MarkedForDeletion []string
var ResourceDirs []string
var Count int

func main() {
	Count = 0

	fmt.Println("* processing")
	err := filepath.Walk("./source", fix)
	if err != nil {
		log.Fatalln(err)
	}

	// rename 'resources' dirs to 'attachments'
	fmt.Println("* renaming resource directories")
	for _, filePath := range ResourceDirs {
		err = rename(filePath)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// this assumes you do not have any duplicate note titles (including in sub directories), comment this out if you dont want it
	fmt.Println("* cleaning up duplicates based on filenames and length of paths")
	deleted := 0
	for _, filePath := range MarkedForDeletion {
		fmt.Println(" -", filePath)
		_ = os.Remove(filePath) // do not care if this errors (paths may have been renamed)
		deleted++
	}

	fmt.Printf("* notes converted: %d, duplicates deleted: %d, actual note count: %d\n", Count, deleted, (Count - deleted))
}

func fix(filePath string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// update image links in md files
	if filepath.Ext(filePath) == ".md" {
		fmt.Println(" +", filePath)
		contents := updateMD(filePath)

		err = ioutil.WriteFile(filePath, []byte(contents), 644)

		if err != nil {
			return err
		}

		checkForDup(filePath)
		AllPaths = append(AllPaths, filePath)
		Count++
	}

	// log for renaming later, dont do here, as filepath.Walk() will be using sub paths
	if info.IsDir() {
		if filepath.Base(filePath) == "resources" {
			ResourceDirs = append(ResourceDirs, filePath)
		}
	}

	return nil
}

func updateMD(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	contents := string(file)

	// update resource insertions, images
	re := regexp.MustCompile(`\!\[.*\]\(resources\/(.*\.[a-zA-Z]+).*\)`)
	matches := re.FindAllStringSubmatch(contents, -1)
	for _, m := range matches {
		fmt.Println("   ", m[0], "->", fmt.Sprintf("![[attachments/%s]]", m[1]))
		contents = strings.Replace(contents, m[0], fmt.Sprintf("![[attachments/%s]]", m[1]), -1)
	}

	// update resource insertions, documents, and fix names back to the original doc name
	re = regexp.MustCompile(`\[(.*)\]\(resources\/(.*)\)`)
	matches = re.FindAllStringSubmatch(contents, -1)
	for _, m := range matches {
		fmt.Println("uikgjghjhgfj")
		err = os.Rename(fmt.Sprintf("%s/resources/%s", filepath.Dir(filePath), m[2]), fmt.Sprintf("%s/resources/%s", filepath.Dir(filePath), m[1]))
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("   ", m[0], "->", fmt.Sprintf("[%s](attachments/%s)", m[1], m[1]))
		contents = strings.Replace(contents, m[0], fmt.Sprintf("[%s](attachments/%s)", m[1], m[1]), -1)
	}

	// update internal links to other pages
	re = regexp.MustCompile(`\[(.*)\]\(quiver\:.*\)`)
	matches = re.FindAllStringSubmatch(contents, -1)
	for _, m := range matches {
		fmt.Println("   ", m[0], "->", fmt.Sprintf("[[%s]]", m[1]))
		contents = strings.Replace(contents, m[0], fmt.Sprintf("[[%s]]", m[1]), -1)
	}

	return contents
}

// rename resources dir to attachments, just to be nice and inline with Obsidian terminology
func rename(filePath string) error {
	err := os.Rename(path.Join(filepath.Dir(filePath), "resources"), path.Join(filepath.Dir(filePath), "attachments"))
	if err != nil {
		return err
	}

	return nil
}

func checkForDup(filePath string) {
	for _, path := range AllPaths {
		existingFileName := filepath.Base(path)
		currentFileName := filepath.Base(filePath)

		if existingFileName == currentFileName {
			// determine whether to delete this instance based on the length of the directory
			// if the dir length is longer, we'll delete it, because it means the note has been sorted and is not lurking in a parent directory as a duplicate
			if len(path) > len(filePath) {
				MarkedForDeletion = append(MarkedForDeletion, filePath)
			} else {
				MarkedForDeletion = append(MarkedForDeletion, path)
			}
		}
	}
}
