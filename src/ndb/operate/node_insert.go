package operate

import (
	"ndb/common"
)

func Insert(node *common.Node, path string, insertValue string) *common.Node {
	
	insertValueMap := CovertValueMap(insertValue)
	
	Locate(node, path, true, func (node *common.Node) {
		for key, value := range insertValueMap {
			node.SetValue(key, []string{value})
		}
	})
	
	return node
}





