package ndb

import (
	"testing"
	"ndb/common"
)


func LoadTestData() *common.Node{
	child1 := new(common.Node)
	child1.SetName("child")
	child1.SetValue("name", []string{"jim"})
	child1.SetValue("age", []string{"20"})
	child1.SetValue("sex", []string{"male"})
	
	child2 := new(common.Node)
	child2.SetName("child")
	child2.SetValue("name", []string{"lily"})
	child2.SetValue("age", []string{"17"})
	child2.SetValue("sex", []string{"female"})
	
	child3 := new(common.Node)
	child3.SetName("child")
	child3.SetValue("name", []string{"tom"})
	child3.SetValue("age", []string{"28"})
	child3.SetValue("sex", []string{"male"})
	
	child4 := new(common.Node)
	child4.SetName("nephew")
	child4.SetValue("name", []string{"lucy"})
	child4.SetValue("age", []string{"12"})
	child4.SetValue("sex", []string{"female"})
	
	parent := new(common.Node)
	parent.SetName("parent")
	parent.SetValue("name", []string{"green"})
	parent.AddChild(child1)
	parent.AddChild(child2)
	parent.AddChild(child3)
	parent.AddChild(child4)
	
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
		
		query := "select:root->parent->child->name:/.*m/"
		result, found := Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
        	if ok && len(children) == 2{
        		if children[0].GetValueString("name") != "jim" || 
        			children[1].GetValueString("name") != "tom" {
        			t.Fatalf("select test fail : %s", query)
        		}
        	}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "select:root->parent->child->age:[15,25]"
		result, found = Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
        	if ok && len(children) == 2 {
        		if children[0].GetValueString("name") != "jim" || 
        			children[1].GetValueString("name") != "lily" {
        			t.Fatalf("select test fail : %s", query)
        		}
        	} else {
        		t.Fatalf("select test fail : %s", query)
        	}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "select:root->parent->child->sex:^fe"
		result, found = Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
        	if ok && len(children) == 1 {
        		if children[0].GetValueString("name") != "lily" {
        			t.Fatalf("select test fail : %s", query)
        		}
        	} else {
        		t.Fatalf("select test fail : %s", query)
        	}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "select:root->parent->child->name:m$"
		result, found = Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
        	if ok && len(children) == 2 {
        		if children[0].GetValueString("name") != "jim" || 
        			children[1].GetValueString("name") != "tom" {
        			t.Fatalf("select test fail : %s", query)
        		}
        	} else {
        		t.Fatalf("select test fail : %s", query)
        	}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "select:root->parent->child->sex:male && age:[15,25]"
		result, found = Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
        	if ok && len(children) == 1 {
        		if children[0].GetValueString("name") != "jim" {
        			t.Fatalf("select test fail : %s", query)
        		}
        	} else {
        		t.Fatalf("select test fail : %s", query)
        	}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "select:root->parent->child"
		result, found = Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
        	if ok && len(children) == 3 {
        		if children[0].GetValueString("name") != "jim" ||
        			children[1].GetValueString("name") != "lily" ||
        			children[2].GetValueString("name") != "tom" {
        			t.Fatalf("select test fail : %s", query)
        		}
        	} else {
        		t.Fatalf("select test fail : %s", query)
        	}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "select:root->parent->:/child|nephew/->sex:female"
		result, found = Execute(node, query)
		if found {
			children, ok := result.([]*common.Node)
        	if ok && len(children) == 2 {
        		if children[0].GetValueString("name") != "lily" || 
        			children[1].GetValueString("name") != "lucy" {
        			t.Fatalf("select test fail : %s", query)
        		}
        	} else {
        		t.Fatalf("select test fail : %s", query)
        	}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
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

