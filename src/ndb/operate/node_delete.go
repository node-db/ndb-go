package operate

import (
	"strings"
	"ndb/common"
)

func Delete(node *common.Node, path string, deleteValue string) (*common.Node, bool) {
	
	columns := []string{}
	clear := false
	
	found := false
	
	if deleteValue != "" {
		if strings.HasPrefix(deleteValue, "[") && strings.HasSuffix(deleteValue, "]") {
			columns = strings.Split(deleteValue[1:len(deleteValue)-1], ",")
		} else if deleteValue == "block" {
			clear = true
		}
	}
	
	Locate(node, path, false, func (node *common.Node) {
		if clear {
			node.ClearValue()
		} else {
			for _, column := range columns {
				node.DeleteValue(strings.TrimSpace(column))
			}
		}
		
		found = true
	})
	
	return node, found
}





