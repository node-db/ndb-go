package main

import (
	"ndb/common"
	"ndb/operate"
	"fmt"
)

func execute(node *common.Node, query string) []*common.Node{
	return nil
}

func read(filename string) (*common.Node, error) {
	return common.Read(filename)
}

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
		
		node = operate.Update(node, "root->parent->child->name:lily", "age=33, phone=13343351822")
		children = operate.Select(node, "root->parent->child->name:lily")
		if len(children) == 1 {
			person := children[0]
			fmt.Println(person.GetValues())
		}
		
		node = operate.Delete(node, "root->parent->child->name:lily", "[age, phone]")
		children = operate.Select(node, "root->parent->child->name:lily")
		if len(children) == 1 {
			person := children[0]
			fmt.Println(person.GetValues())
		}
		
		node = operate.Delete(node, "root->parent->child->name:lily", "block")
		children = operate.Select(node, "root->parent->child->name:lily")
		if len(children) == 0 {
			fmt.Println("Not found : Lily")
		}
	}

}