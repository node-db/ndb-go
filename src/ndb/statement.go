package ndb

import (
	"strings"
	"ndb/common"
	"ndb/operate"
)

func Execute(node *common.Node, query string) interface{} {
	command := query
	
	if strings.Contains(query, ":") {
		command = strings.TrimSpace(query[0 : strings.Index(query, ":")])
		query = strings.TrimSpace(query[strings.Index(query, ":") + 1 : len(query)])
	}
	
	queryItems := strings.Split(query, "!!")
	
	if queryItems != nil && len(queryItems) > 0 {
		path := strings.TrimSpace(queryItems[0])
		
		value := ""
		if len(queryItems) > 1 {
			value = strings.TrimSpace(queryItems[1])
		}
		
		if command != "" {
			command = strings.ToLower(command)
			if command == "select" || command == "one" || command == "exist" {
				result := operate.Select(node, path)
				
				if command == "one" {
					if result != nil && len(result) > 0 {
						return result[0]
					} else {
						return nil
					}
				} else if command == "exist"{
					if  result != nil && len(result) > 0 {
						return true 
					} else {
						return false
					}
				}
				
				return result
			} else if command == "update" {
				return operate.Update(node, path, value)
			} else if command == "delete" {
				return operate.Delete(node, path, value)
			} else if command == "insert" {
				return operate.Insert(node, path, value)
			}
		}
	}
    
    return nil
}

func Read(filename string) (*common.Node, error) {
	return common.Read(filename)
}