package common

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	pathSplit := strings.Split(path, "\\")
	size := len(pathSplit)
	pathSplit = strings.Split(path, pathSplit[size-1])
	currPath := strings.Replace(pathSplit[0], "\\", "/", size-1)
	return currPath
}

func Read(filename string) (*Node, error) {
	if filename == "" {
		return nil, errors.New("Filename is NULL")
	}
	if !strings.Contains(filename, "\\") && !strings.Contains(filename, "/") {
		filename = GetCurrPath() + "/" + filename
	}
	content, err := ReadFile(filename)
	if err == nil {
		node, _ := Parse(0, content, nil)
		return node, nil
	} else {
		return nil, err
	}
}

func ReadFile(filename string) ([]string, error) {
	fi, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	result := []string{}

	bfRd := bufio.NewReader(fi)
	for {
		line, err := bfRd.ReadBytes('\n')
		result = append(result, string(line))

		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
	}
	return result, nil
}

func Parse(linenum int, contents []string, parent *Node) (*Node, int) {
	if parent == nil {
		parent = new(Node)
	}

	if contents != nil {
		for i := linenum; i < len(contents); i++ {
			line := strings.TrimSpace(contents[i])
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if strings.HasSuffix(line, "{") {
				node := new(Node)
				node.name = strings.TrimSpace(line[:strings.LastIndex(line, "{")])
				nodeChild, _line := Parse(i+1, contents, node)
				parent.AddChild(nodeChild)

				i = _line
			} else if strings.HasPrefix(line, "}") {
				return parent, i
			} else {
				if strings.Index(line, ":") > 0 {
					itemName := strings.TrimSpace(line[:strings.Index(line, ":")])
					itemValue := strings.TrimSpace(line[strings.Index(line, ":")+1 : len(line)])

					valueList := parent.GetValue(itemName)
					if valueList != nil {
						valueList = []string{}
					}
					valueList = append(valueList, itemValue)
					parent.SetValue(itemName, valueList)
				}
			}
		}
	}

	return parent, len(contents)
}
