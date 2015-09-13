package ndb

import (
	"strings"
)

func Execute(node *Node, query string) (interface{}, bool) {
	command := query

	if strings.Contains(query, ":") {
		command = strings.TrimSpace(query[0:strings.Index(query, ":")])
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
				result, found := Select(node, path)

				if command == "one" {
					if found {
						return result[0], true
					} else {
						return nil, false
					}
				} else if command == "exist" {
					if found {
						return nil, true
					} else {
						return nil, false
					}
				}

				return result, found
			} else if command == "update" {
				return Update(node, path, value)
			} else if command == "delete" {
				return Delete(node, path, value)
			} else if command == "insert" {
				return Insert(node, path, value)
			} else {
				panic("unknow operate : " + command)
			}
		}
	} else {
		panic("unknow query : " + query)
	}

	return nil, false
}