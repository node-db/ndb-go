package main

import (
	"ndb"
	"ndb/common"
	"fmt"
)

func main() {
	node, err := ndb.Read("example.ndb")
	if err == nil {
		query := "select:root->parent->child->name:/.*m/"
		result, found := ndb.Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
			if ok {
				for _, child := range children {
					fmt.Println(child.GetValues())
				}
			}
		}
	}
}

