package operate

import (
	"ndb/data"
)

func Update(node *data.Node, path string, updateValue string) (*data.Node, bool) {
	
	updateValueMap := CovertValueMap(updateValue)
	found := false
	
	Locate(node, path, false, func (node *data.Node) {
		for key, value := range updateValueMap {
			node.SetValue(key, []string{value})
		}
		found = true
	})
	
	return node, found
}





