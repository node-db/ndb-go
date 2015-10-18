package ndb

import (
	"testing"
	"fmt"
	"os"
)

func TestReadFile(t *testing.T) {
	
}

func TestParseStringToNode(t *testing.T) {
	
}

func TestWriteFile(t *testing.T) {
	filename := "/var/test/example-1.ndb"
	
	node := LoadTestData()
	err := WriteFile(filename, node)
	if err != nil {
		t.Fatalf("Write node fail")
	}
	node, err = ReadFile(filename)
	if err != nil {
		t.Fatalf("Read node fail")
	}
	query := "root->parent->child->name:lily"
	if node != nil {
		children, _, _ := Select(node, query)
		if len(children) == 1 {
			person := children[0]
			if person.GetValueString("age") != "17" {
				t.Fatalf("Write data is error : %s", query)
			}
		} else {
			t.Fatalf("Write data is error : %s", query)
		}
	} else {
		t.Fatalf("Write data is error : %s", query)
	}
	
	err = os.Remove(filename)
	if err != nil {
		t.Fatalf("Delete file fail")
	}
}

func TestParseNodeToString(t *testing.T) {
	node := LoadTestData()
	if node != nil {
		nodeStr := ParseNodeToString(node)
		fmt.Println(nodeStr)
	} else {
		t.Fatalf("Node is nil")
	}
}
