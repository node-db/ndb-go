package operate

import (
	"ndb/data"
)

func Select(node *data.Node, path string) ([]*data.Node, bool) {
	
	result := []*data.Node{}
	found := false
	
	Locate(node, path, false, func (node *data.Node) {
		result = append(result, node)
		found = true
	})
	
	return result, found
}





