package operate

import (
	"ndb/data"
)

func Insert(node *data.Node, path string, insertValue string) (*data.Node, bool) {
	
	insertValueMap := CovertValueMap(insertValue)
	found := false
	
	Locate(node, path, true, func (node *data.Node) {
		for key, value := range insertValueMap {
			node.SetValue(key, []string{value})
		}
		found = true
	})
	
	return node, found
}





