package ndb

import (
	"fmt"
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

func ReadAsList(filename string) ([]string, error) {
	if filename == "" {
		return nil, errors.New("Filename is NULL")
	}
	if !strings.Contains(filename, "\\") && !strings.Contains(filename, "/") {
		filename = GetCurrPath() + "/" + filename
	}
	
	fin, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fin.Close()
	bfRd := bufio.NewReader(fin)
	
	content := []string{}
	
	for {
		line, err := bfRd.ReadBytes('\n')
		content = append(content, string(line))

		if err != nil {
			if err == io.EOF {
				break
			}
		}
	}
	
	return content, nil
}

func Read(filename string) (*Node, error) {
	content, err := ReadAsList(filename)
	if err == nil {
		node, _ := ParseStringToNode(0, content, nil)
		return node, nil
	} else {
		return nil, err
	}
}

func ParseStringToNode(linenum int, contents []string, parent *Node) (*Node, int) {
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
				node.SetName(strings.TrimSpace(line[:strings.LastIndex(line, "{")]))
				nodeChild, _line := ParseStringToNode(i+1, contents, node)
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

func WriteFile(filename string, node *Node, indentFlag string) error {
	if filename == "" {
		return errors.New("Filename is NULL")
	}
	
	fout, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fout.Close()
	
	nodeStr := Print(node, indentFlag)
	fout.WriteString(nodeStr)
	
	return nil
}

func Print(node *Node, indentFlag string) string {
	return ParseNodeToString(0, node, indentFlag)
}

func ParseNodeToString(indent int, node *Node, indentFlag string) string {
	nodeStr := ""
	
	name := node.GetName()
	values := node.GetValues()
	children := node.GetChileren()
	
	nextIndent := indent + 1
	
	if name != "" {
		nodeStr = strings.Repeat(indentFlag, indent) + name + "{\n"
	} else {
		nextIndent = indent
	}
	
	for key, value := range values {
		for _, valueItem := range value {
			nodeStr = nodeStr + strings.Repeat(indentFlag, nextIndent) + key + ":" + valueItem + "\n"
		}
	}
	
	for _, child := range children {
		nodeStr = nodeStr + ParseNodeToString(nextIndent, child, indentFlag)
	}
	
	if name != "" {
		nodeStr = nodeStr + strings.Repeat(indentFlag, indent) + "}\n"
	}
	
	return nodeStr
}

func Redirect(target string, node interface{}) {
	outputData := new(Node)
	switch node.(type) {
		case []*Node:
			nodeList := node.([]*Node)
			for _, nodeItem := range nodeList {
				nodeItem.SetName("result")
				outputData.AddChild(nodeItem)
			}
			break
		case *Node:
			outputData = node.(*Node)
			break
	}
	
	if strings.ToLower(target) == "print" {
		nodeStr := Print(outputData, "\t")
		fmt.Println(nodeStr)
	} else {
		WriteFile(target, outputData, "\t")
	}
}
