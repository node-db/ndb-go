package operate

import (
	"ndb/common"
)

func Update(node *common.Node, path string, updateValue string) *common.Node {
	
	updateValueMap := CovertValueMap(updateValue)
	
	Locate(node, path, false, func (node *common.Node) {
		for key, value := range updateValueMap {
			node.SetValue(key, []string{value})
		}
	})
	
	return node
}





