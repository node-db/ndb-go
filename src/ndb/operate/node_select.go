package operate

import (
	"ndb/common"
)

func Select(node *common.Node, path string) []*common.Node {
	
	result := []*common.Node{}
	
	Locate(node, path, false, func (node *common.Node) {
		result = append(result, node)
	})
	
	return result
}





