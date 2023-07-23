package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type DocKey struct {
	RawText    string
	Key        string
	LineNumber int
	Spaces     int
}

func findTreeWithValidPath(roots []TreeNode, path []string) (TreeNode, error) {
	for _, root := range roots {
		if root.isValidPath(path) {
			return root, nil
		}
	}
	return TreeNode{}, errors.New("no tree contains the path")
}

func helpText(args []string) {
	fmt.Fprintln(os.Stderr, "     __  ______   _____ __  __ ")
	fmt.Fprintln(os.Stderr, "    /  |/  /   | / ___// / / /")
	fmt.Fprintln(os.Stderr, "   / /|_/ / /| | \\__ \\/ /_/ /")
	fmt.Fprintln(os.Stderr, "  / /  / / ___ |___/ / __  /  -- combine code and yaml")
	fmt.Fprintln(os.Stderr, " /_/  /_/_/  |_/____/_/ /_/  ")
	fmt.Fprintln(os.Stderr, " Version 0.1")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Usage: mash [code|yaml] <file> [into|over] <yaml_file> at <path.seperated.by.dots>")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "\t[code|yaml]        whether <file> should be inserted as code (with a | for multiline) or as yaml")
	fmt.Fprintln(os.Stderr, "\t<file>             the name of the code or yaml to insert into another file")
	fmt.Fprintln(os.Stderr, "\t[into|over]        into will produce a new file, over will overwrite")
	fmt.Fprintln(os.Stderr, "\t<yaml_file>        the yaml file to insert into. JSON-looking yaml files are not supported")
	fmt.Fprintln(os.Stderr, "\tat                 The word at. The design is very human")
	fmt.Fprintln(os.Stderr, "\t<path>             The sequence of keys in <yaml_file> representing the location to insert")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "Got %d args: %v\n", len(args), args)
	os.Exit(1)
}

func checkMode(mode string) {
	if mode != "into" && mode != "over" {
		fmt.Fprintf(os.Stderr, "Mode must be one of 'into' or 'over': '%s' is invalid\n", mode)
		os.Exit(1)
	}
}

func checkFiletype(filetype string) {
	if filetype != "code" && filetype != "yaml" {
		fmt.Fprintf(os.Stderr, "Filetype must be one of 'code' or 'yaml': '%s' is invalid\n", filetype)
		os.Exit(1)
	}
}

func main() {

	args := os.Args[1:]
	if len(args) != 6 {
		helpText(args)
	}

	filetype := args[0]
	codeFile := args[1]
	mode := args[2]
	yamlFile := args[3]
	rawPath := args[5]

	checkMode(mode)
	checkFiletype(filetype)

	// open the yaml file
	yamlFileHandle, err := os.Open(yamlFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer yamlFileHandle.Close()

	// create a regex to find keys
	keys_regex, _ := regexp.Compile(`^\s*\S+:\n|^\s*\S+:`)
	var keys []DocKey
	line := 1

	// scan the entire file for keys
	scanner := bufio.NewScanner(yamlFileHandle)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "\t") {
			fmt.Fprintln(os.Stderr, "Cannot mash tab-indented files. Convert to spaces and try again")
			os.Exit(1)
		}
		line_key := string(keys_regex.Find([]byte(scanner.Text())))
		if line_key != "" {
			keys = append(keys, DocKey{RawText: line_key, Key: stripSpacesAndQuotes(line_key), LineNumber: line, Spaces: countIndent(line_key)})
		}
		line++
	}

	// exit if no keys found
	if len(keys) == 0 {
		fmt.Fprintln(os.Stderr, "No keys in input yaml file")
		os.Exit(1)
	}

	// build trees out of keys in the file
	var roots []TreeNode
	treeKeyLists := splitKeyListIntoTrees(keys)
	for _, keys := range treeKeyLists {
		roots = append(roots, buildTree(keys))
	}

	// see which (if any) trees have the desired path
	// returns first tree with path
	path := strings.Split(rawPath, ".")
	root, err := findTreeWithValidPath(roots, path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// else {
	// 	fmt.Printf("Path contained in tree rooted at '%s'\n", root.Value.Key)
	// 	fmt.Println(root.isValidPath(path))
	// }
	startNode := root.getPathStartNode(path)
	//fmt.Printf("Node starts on line %d\n", startNode.LineNumber)

	// insert code/yaml into the file
	_, err = yamlFileHandle.Seek(0, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	codeFileHandle, err := os.Open(codeFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer codeFileHandle.Close()

	var outFile *os.File
	if true {
		outFile = os.Stdout
	} else {
		outFile, err = os.Create("output.yml")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer outFile.Close()
	}

	line = 1
	codeScanner := bufio.NewScanner(codeFileHandle)
	scanner = bufio.NewScanner(yamlFileHandle)
	for scanner.Scan() {
		if line == startNode.LineNumber {
			outFile.WriteString(scanner.Text())
			// for code, use | or > to make it multiline
			// yaml can be inserted without
			if filetype == "code" {
				outFile.WriteString(" |\n")
			} else {
				outFile.WriteString("\n")
			}
			for codeScanner.Scan() {
				outFile.WriteString(padString(codeScanner.Text()+"\n", startNode.Spaces+2))
			}
		} else {
			outFile.WriteString(scanner.Text())
			outFile.WriteString("\n")
		}
		line++
	}
}
