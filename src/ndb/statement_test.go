package ndb

import (
	"testing"
	"ndb/common"
)


func LoadTestData() *common.Node{
	
	NewChild := func (node string, name string, age string, sex string) *common.Node {
		child := new(common.Node)
		
		child.SetName(node)
		child.SetValue("name", []string{name})
		child.SetValue("age", []string{age})
		child.SetValue("sex", []string{sex})
		
		return child
	}
	
	child1 := NewChild("child", "jim", "20", "male")
	child2 := NewChild("child", "lily", "17", "female")
	child3 := NewChild("child", "tom", "28", "male")
	child4 := NewChild("nephew", "lucy", "12", "female")
	
	parent := new(common.Node)
	parent.SetName("parent")
	parent.SetValue("name", []string{"green"})
	parent.AddChildren([]*common.Node{child1, child2, child3, child4})
	
	root := new(common.Node)
	root.SetName("root")
	root.AddChild(parent)
	
	node := new(common.Node)
	node.AddChild(root)
	
	return node
}

func TestExits(t *testing.T){
	node := LoadTestData()
	
	if node != nil {
		query := "exist:root->parent->child->name:jim"
        result, found := Execute(node, query)
        if result != nil && found == false {
        	t.Fatalf("exits test fail : %s", query)
        }
        
        query = "exist:root->parent->child->sex:male && name:m$"
        result, found = Execute(node, query)
        if found == false {
        	t.Fatalf("exits test fail : %s", query)
        }
        
        query = "exist:root->parent->child->sex:female && name:m$"
        result, found = Execute(node, query)
        if found == true {
        	t.Fatalf("exits test fail : %s", query)
        }
	}
}

func TestOne(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		query := "one:root->parent->child->sex:male"
        result, _ := Execute(node, query)
        child, ok := result.(*common.Node)
        if ok {
        	if child.GetValueString("name") != "jim" || child.GetValueString("age") != "20" {
        		t.Fatalf("one test fail : %s", query)
        	}
        }
	}
}

func TestSelect(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		
		SelectAssert := func (query string, expect []string) {
			result, found := Execute(node, query)
			if found {
				children, ok := result.([]*common.Node)
	        	if ok && len(children) == len(expect) {
	        		for i := 0 ; i < len(expect) ; i++ {
	        			child := children[i]
	        			if child.GetValueString("name") != expect[i] {
	        				t.Fatalf("select test fail : %s", query)
	        				break
	        			}
	        		}
	        	}
			} else {
				t.Fatalf("select test fail : %s", query)
			}
		}
		
		query := "select:root->parent->child->name:/.*m/"
		SelectAssert(query, []string{"jim", "tom"})
		
		query = "select:root->parent->child->age:[15,25]"
		SelectAssert(query, []string{"jim", "lily"})
		
		query = "select:root->parent->child->sex:^fe"
		SelectAssert(query, []string{"lily"})
		
		query = "select:root->parent->child->name:m$"
		SelectAssert(query, []string{"jim", "tom"})

		query = "select:root->parent->child->sex:male && age:[15,25]"
		SelectAssert(query, []string{"jim"})
		
		query = "select:root->parent->child"
		SelectAssert(query, []string{"jim", "lily", "tom"})

		query = "select:root->parent->:/child|nephew/->sex:female"
		SelectAssert(query, []string{"lily", "lucy"})
	}
}

func TestUpdate(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		query := "update:root->parent->child->name:jim !! age=21, address=China"
		result, found := Execute(node, query)
		if found {
			updateResult, _ := Execute(result.(*common.Node), "one:root->parent->child->name:jim")
			child, ok := updateResult.(*common.Node)
			if ok {
				if child.GetValueString("age") != "21" || 
					child.GetValueString("address") != "China" {
						t.Fatalf("update test fail : %s", query)
				}
			} else {
				t.Fatalf("update test fail : %s", query)
			}
		} else {
			t.Fatalf("update test fail : %s", query)
		}
	}
}

func TestDelete(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		query := "delete:root->parent->child->name:jim !! [sex, age]"
		result, found := Execute(node, query)
		if found {
			deleteResult, _ := Execute(result.(*common.Node), "one:root->parent->child->name:jim")
			child, ok := deleteResult.(*common.Node)
			if ok {
				if child.GetValueString("name") != "jim" || 
					child.GetValueString("sex") != "" ||
					child.GetValueString("age") != "" {
						t.Fatalf("delete test fail : %s", query)
				}
			} else {
				t.Fatalf("delete test fail : %s", query)
			}
		} else {
			t.Fatalf("delete test fail : %s", query)
		}

		query = "delete:root->parent->child->name:jim !! block"
		result, found = Execute(node, query)
		if found {
			deleteResult, _ := Execute(result.(*common.Node), "select:root->parent->child->name:jim")
			children, _ := deleteResult.([]*common.Node)
			if len(children) > 0 {
				t.Fatalf("delete test fail : %s", query)
			}
		} else {
			t.Fatalf("delete test fail : %s", query)
		}
	}
}

func TestInsert(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		query := "insert:root->parent->child !! name=bill, sex=male, age=31"
		result, _ := Execute(node, query)
		
		insertResult, _ := Execute(result.(*common.Node), "one:root->parent->child->name:bill")
		child, ok := insertResult.(*common.Node)
		if ok {
			if child.GetValueString("name") != "bill" || 
				child.GetValueString("sex") != "male" ||
				child.GetValueString("age") != "31" {
					t.Fatalf("insert test fail : %s", query)
			}
		} else {
			t.Fatalf("insert test fail : %s", query)
		}
	}
}

