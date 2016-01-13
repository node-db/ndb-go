package ndb

import (
	"regexp"
	"strconv"
	"strings"
)

func Locate(node *Node, query string, isCreate bool, action func(node *Node)) {
	if query == "" || node == nil {
		return
	}

	queryKey := query
	subQuery := query

	if strings.Contains(query, "->") {
		queryKey = strings.TrimSpace(query[:strings.Index(query, "->")])
		subQuery = strings.TrimSpace(query[strings.Index(query, "->")+2 : len(query)])
	}

	if subQuery != queryKey || strings.HasPrefix(queryKey, ":") {
		if strings.HasPrefix(queryKey, ":") {
			// 路径模糊查询
			exp := queryKey[1:]
			children := node.GetChileren()
			for _, child := range children {
				key := child.GetName()
				if CheckValue(key, exp) {
					if strings.HasPrefix(subQuery, ":") {
						Locate(node, key, isCreate, action)
					} else {
						Locate(child, subQuery, isCreate, action)
					}
				}
			}
		} else {
			children := node.FindChildByName(queryKey)
			for _, child := range children {
				Locate(child, subQuery, isCreate, action)
			}

		}
	} else {
		if strings.Contains(subQuery, ":") {

			matchItems := strings.Split(subQuery, "&&")
			matchResult := true

			for _, matchItem := range matchItems {
				items := strings.Split(strings.TrimSpace(matchItem), ":")
				if len(items) == 2 {
					key := strings.TrimSpace(items[0])
					exp := strings.TrimSpace(items[1])

					value := node.GetValueString(key)
					if !CheckValue(value, exp) {
						matchResult = false
					}
				}
			}

			if matchResult {
				action(node)
			}
		} else {
			children := node.FindChildByName(queryKey)
			if isCreate == true {
				newNode := new(Node)
				newNode.SetName(queryKey)
				action(newNode)
				node.AddChild(newNode)
			} else {
				for _, child := range children {
					action(child)
				}
			}
		}
	}

}

func CheckValue(value string, exp string) bool {
	if value == "" || exp == "" {
		return false
	}

	//regex valueression match
	if len(exp) > 2 && strings.HasPrefix(exp, "/") && strings.HasSuffix(exp, "/") {
		exp = exp[1 : len(exp)-1]
		reg := regexp.MustCompile(exp)
		if reg.MatchString(value) {
			return true
		}
	}

	//number region match
	if len(exp) > 3 && strings.HasPrefix(exp, "[") && strings.HasSuffix(exp, "]") {
		exp = exp[1 : len(exp)-1]
		numbers := strings.Split(exp, ",")
		if numbers != nil && len(numbers) == 2 {
			val, valErr := strconv.Atoi(value)
			min, maxErr := strconv.Atoi(numbers[0])
			max, minErr := strconv.Atoi(numbers[1])
			if valErr != nil || maxErr != nil || minErr != nil {
				return false
			}
			if val <= max && val >= min {
				return true
			}
		}
	}

	//startswith match
	if strings.HasPrefix(exp, "^") {
		if strings.HasPrefix(value, exp[1:]) {
			return true
		}
	}

	//endswith match
	if strings.HasSuffix(exp, "$") {
		if strings.HasSuffix(value, exp[:len(exp)-1]) {
			return true
		}
	} else {
		if value != "" && value == exp {
			return true
		}
	}

	return false
}

func CovertValueMap(updateValue string) map[string]string {
	valueMap := make(map[string]string)
	values := strings.Split(updateValue, ",")
	for _, value := range values {
		valuePair := strings.Split(value, "=")
		if len(valuePair) == 2 {
			valueMap[strings.TrimSpace(valuePair[0])] = strings.TrimSpace(valuePair[1])
		}
	}

	return valueMap
}

func Delete(node *Node, path string, deleteValue string) (*Node, bool, error) {

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

	Locate(node, path, false, func(node *Node) {
		if clear {
			node.ClearValue()
		} else {
			for _, column := range columns {
				node.DeleteValue(strings.TrimSpace(column))
			}
		}

		found = true
	})

	return node, found, nil
}

func Insert(node *Node, path string, insertValue string) (*Node, bool, error) {

	insertValueMap := CovertValueMap(insertValue)
	found := false

	Locate(node, path, true, func(node *Node) {
		for key, value := range insertValueMap {
			node.SetValue(key, []string{value})
		}
		found = true
	})

	return node, found, nil
}

func Select(node *Node, path string) ([]*Node, bool, error) {

	result := []*Node{}
	found := false

	Locate(node, path, false, func(node *Node) {
		result = append(result, node)
		found = true
	})

	return result, found, nil
}

func Update(node *Node, path string, updateValue string) (*Node, bool, error) {

	updateValueMap := CovertValueMap(updateValue)
	found := false

	Locate(node, path, false, func(node *Node) {
		for key, value := range updateValueMap {
			node.SetValue(key, []string{value})
		}
		found = true
	})

	return node, found, nil
}

func Script(node *Node, scriptFilename string) (*Node, error) {

	script, err := ReadAsList(scriptFilename)
	if err != nil {
		for _, query := range script {
			query = strings.TrimSpace(query)
			if query != "" {
				result, found, err := Execute(node, query)
				if found && err == nil {
					switch result.(type) {
					case []*Node:
						break
					case *Node:
						node = result.(*Node)
						break
					}
				}
			}
		}
	}

	return node, nil
}
