package operate

import (
	"ndb/common"
	"strings"
	"strconv"
	"regexp"
)

func Locate(node *common.Node, query string, isCreate bool, action func(node *common.Node)) {
	if query == "" || node == nil {
		return
	}
	
	queryKey := query
	subQuery := query
	
	if strings.Contains(query, "->") {
		queryKey = strings.TrimSpace(query[:strings.Index(query, "->")])
		subQuery = strings.TrimSpace(query[strings.Index(query, "->") + 2 : len(query)])
	}
	
	if subQuery != queryKey || strings.HasPrefix(queryKey, ":") {
		if strings.HasPrefix(queryKey, ":") {
			// 路径模糊查询
			exp := queryKey[1:]
			children := node.GetChileren()
			for _, child := range children {
				key := child.GetName()
				if  CheckValue(key, exp) {
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
				newNode := new(common.Node)
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
	if len(exp) > 2 && strings.HasPrefix(exp, "/")  && strings.HasSuffix(exp, "/"){
		exp = exp[1 : len(exp) - 1]
		reg := regexp.MustCompile(exp)
		if reg.MatchString(value) {
			return true
		}
	}
	
	//number region match
	if len(exp) > 3 && strings.HasPrefix(exp, "[")  && strings.HasSuffix(exp, "]"){
		exp = exp[1 : len(exp) - 1]
		numbers := strings.Split(exp, ",")
		if numbers != nil && len(numbers) == 2 {
			val, valErr := strconv.Atoi(value)
			max, maxErr := strconv.Atoi(numbers[0])
			min, minErr := strconv.Atoi(numbers[1])
			if valErr != nil || maxErr != nil || minErr != nil {
				return false
			}
			if val <= max && val >=min {
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
		if strings.HasSuffix(value, exp[: len(exp) - 1]) {
			return true
		}
	} else {
		if value != "" && value == exp {
			return true
		}
	}
	
	return false
}


func CovertValueMap(updateValue string) (map[string]string) {
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