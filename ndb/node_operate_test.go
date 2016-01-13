package ndb

import (
	"testing"
)

func TestOperateSelect(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		query := "root->parent->child->name:lily"
		children, _, _ := Select(node, query)
		if len(children) == 1 {
			person := children[0]
			if person.GetValueString("age") != "17" {
				t.Fatalf("select test fail : %s", query)
			}
		} else {
			t.Fatalf("select test fail : %s", query)
		}
		
		query = "root->parent->child->name:m$ && sex:^m"
		children, _, _ = Select(node, query)
		if len(children) != 2 {
			t.Fatalf("select test fail : %s", query)
		}
	} else {
		t.Fatalf("Node is nil")
	}
}

func TestOperateDelete(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		path := "root->parent->child->name:lily"
		value := "[age]"
		node, _, _ := Delete(node, path, value)
		children, _, _ := Select(node, "root->parent->child->name:lily")
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
		node, _, _ = Delete(node, path, value)
		children, _, _ = Select(node, "root->parent->child->name:lily")
		if len(children) > 0 {
			t.Fatalf("delete test fail : %v, %v", path, value)
		}
	} else {
		t.Fatalf("Node is nil")
	}
}

func TestOperateUpdate(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		path := "root->parent->child->name:lily"
		value := "age=33, phone=13343351822"
		node, _, _ := Update(node, path, value)
		children, _, _ := Select(node, "root->parent->child->name:lily")
		if len(children) == 1 {
			person := children[0]
			if person.GetValueString("age") != "33" || person.GetValueString("phone") != "13343351822" {
				t.Fatalf("update test fail : %s, %s", path, value)
			} 
		} else {
			t.Fatalf("update test fail : %s, %s", path, value)
		}
	} else {
		t.Fatalf("Node is nil")
	}
}

func TestOperateInsert(t *testing.T) {
	node := LoadTestData()
	
	if node != nil {
		path := "root->parent->house"
		value := "phone=82988679, address=Foshan"
		node, _, _ := Insert(node, path, value)
		houses, _, _ := Select(node, "root->parent->house")
		if len(houses) == 1 {
			house := houses[0]
			if house.GetValueString("phone") != "82988679" || house.GetValueString("address") != "Foshan" {
				t.Fatalf("insert test fail : %s, %s", path, value)
			} 
		} else {
			t.Fatalf("insert test fail : %s, %s", path, value)
		}
	} else {
		t.Fatalf("Node is nil")
	}
}
