package operate

import (
	"ndb/common"
)

func Insert(node *common.Node, path string, insertValue string) (*common.Node, bool) {
	
	insertValueMap := CovertValueMap(insertValue)
	found := false
	
	Locate(node, path, true, func (node *common.Node) {
		for key, value := range insertValueMap {
			node.SetValue(key, []string{value})
		}
		found = true
	})
	
	return node, found
}





