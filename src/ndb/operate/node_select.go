package operate

import (
	"ndb/common"
)

func Select(node *common.Node, path string) ([]*common.Node, bool) {
	
	result := []*common.Node{}
	found := false
	
	Locate(node, path, false, func (node *common.Node) {
		result = append(result, node)
		found = true
	})
	
	return result, found
}





