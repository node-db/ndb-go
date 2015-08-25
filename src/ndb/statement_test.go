package main

import (
	"fmt"
	"ndb/common"
	"ndb/operate"
)

func main() {
	node, err := common.Read("example.ndb")
	if err == nil {
		children := operate.Select(node, "root->parent->child->name:lily")
		if len(children) == 1 {
			person := children[0]
			fmt.Println(person.GetValues())
		}
		children = operate.Select(node, "root->parent->child->name:m$ && sex:^m")
		fmt.Println(len(children))
	}

}
