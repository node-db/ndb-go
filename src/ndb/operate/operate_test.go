package operate

import (
	"testing"
	"ndb/common"
)

func LoadTestData() *common.Node {
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

func TestSelect(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		query := "root->parent->child->name:lily"
		children, _ := Select(node, query)
		if len(children) == 1 {
			person := children[0]
			if person.GetValueString("age") != "17" {
				t.Fatalf("select test fail : %s", query)
			}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "root->parent->child->name:m$ && sex:^m"
		children, _ = Select(node, query)
		if len(children) != 2 {
			t.Fatalf("select test fail : %s", query)
		}

	}
}

func TestDelete(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		path := "root->parent->child->name:lily"
		value := "[age]"
		node, _ := Delete(node, path, value)
		children, _ := Select(node, "root->parent->child->name:lily")
		if len(children) == 1 {
			person := children[0]
			if person.GetValueString("age") != "" {
				t.Fatalf("delete test fail : %s, %s", path, value)
			}
		} else {
			t.Fatalf("delete test fail : %s, %s", path, value)
		}
		
		path = "root->parent->child->name:lily"
		value = "block"
		node, _ = Delete(node, path, value)
		children, _ = Select(node, "root->parent->child->name:lily")
		if len(children) > 0 {
			t.Fatalf("delete test fail : %v, %v", path, value)
		}
	}
}

func TestUpdate(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		path := "root->parent->child->name:lily"
		value := "age=33, phone=13343351822"
		node, _ := Update(node, path, value)
		children, _ := Select(node, "root->parent->child->name:lily")
		if len(children) == 1 {
			person := children[0]
			if person.GetValueString("age") != "33" || person.GetValueString("phone") != "13343351822" {
				t.Fatalf("update test fail : %s, %s", path, value)
			} 
		} else {
			t.Fatalf("update test fail : %s, %s", path, value)
		}
	}
}

func TestInsert(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		path := "root->parent->house"
		value := "phone=82988679, address=Foshan"
		node, _ := Insert(node, path, value)
		houses, _ := Select(node, "root->parent->house")
		if len(houses) == 1 {
			house := houses[0]
			if house.GetValueString("phone") != "82988679" || house.GetValueString("address") != "Foshan" {
				t.Fatalf("insert test fail : %s, %s", path, value)
			} 
		} else {
			t.Fatalf("insert test fail : %s, %s", path, value)
		}
	}
}
