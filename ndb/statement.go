package ndb

import (
	"errors"
	"strings"
)

func Execute(node *Node, query string) (interface{}, bool, error) {

	if node == nil {
		return nil, false, errors.New("Node is NULL")
	}

	var result interface{} = nil
	var found bool = false
	var err error = nil

	command := query

	if strings.Contains(query, ":") {
		command = strings.TrimSpace(query[0:strings.Index(query, ":")])
		query = strings.TrimSpace(query[strings.Index(query, ":")+1 : len(query)])
	}

	queryItems := strings.Split(query, "!!")

	if queryItems != nil && len(queryItems) > 0 {
		path := strings.TrimSpace(queryItems[0])

		value := ""
		if len(queryItems) > 1 {
			value = strings.TrimSpace(queryItems[1])
		}

		redirect := ""
		if strings.Contains(path, ">>") {
			pathItems := strings.Split(path, ">>")
			if len(pathItems) == 2 {
				path = strings.TrimSpace(pathItems[0])
				redirect = strings.TrimSpace(pathItems[1])
			}
		} else if strings.Contains(value, ">>") {
			valueItems := strings.Split(value, ">>")
			if len(valueItems) == 2 {
				value = strings.TrimSpace(valueItems[0])
				redirect = strings.TrimSpace(valueItems[1])
			}
		}

		if command != "" {
			command = strings.ToLower(command)
			if command == "select" || command == "one" || command == "exist" {
				selectResult, selectFound, selectErr := Select(node, path)

				found = selectFound
				err = selectErr

				if command == "select" {
					result = selectResult
				} else if command == "one" {
					if found {
						result = selectResult[0]
					}
				} else if command == "exist" {
					result = nil
				}
			} else if command == "update" {
				result, found, err = Update(node, path, value)
			} else if command == "delete" {
				result, found, err = Delete(node, path, value)
			} else if command == "insert" {
				result, found, err = Insert(node, path, value)
			} else if command == "script" {

			} else {
				err = errors.New("unknow operate : " + command)
			}
		}

		if redirect != "" {
			Redirect(redirect, result)
		}
	} else {
		err = errors.New("unknow query : " + query)
	}

	return result, found, err
}
