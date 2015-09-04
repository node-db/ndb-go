package operate

import (
	"testing"
	"ndb/data"
)

func LoadTestData() *data.Node {
	
	NewChild := func (node string, name string, age string, sex string) *data.Node {
		child := new(data.Node)
		
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
	
	parent := new(data.Node)
	parent.SetName("parent")
	parent.SetValue("name", []string{"green"})
	parent.AddChildren([]*data.Node{child1, child2, child3, child4})
	
	root := new(data.Node)
	root.SetName("root")
	root.AddChild(parent)
	
	node := new(data.Node)
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
