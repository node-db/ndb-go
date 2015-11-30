package ndb

import (
	"testing"
	"os"
	"fmt"
	"strconv"
)

func LoadTestData() *Node {
	
	dataSource := "local"
	
	if dataSource == "local" {

		NewChild := func(node string, name string, age string, sex string) *Node {
			child := new(Node)

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

		parent := new(Node)
		parent.SetName("parent")
		parent.SetValue("name", []string{"green"})
		parent.AddChildren([]*Node{child1, child2, child3, child4})

		root := new(Node)
		root.SetName("root")
		root.AddChild(parent)

		node := new(Node)
		node.AddChild(root)

		return node
	} else {
		node, _ := Read(dataSource)
		return node
	}
}

func ValueAssert(node *Node, field string, expect string, query string) {
	if node.GetValue(field)[0] != expect {
		fmt.Printf("expect %s but %s, %s\n", expect, node.GetValue(field)[0], query)
	}
}

func NullAssert(node *Node, field string, query string) {
	length := len(node.GetValue(field))
	if length > 0 {
		fmt.Printf("expect null but not, %s\n", query)
	}
}

func LengthAssert(nodes []*Node, expect int, query string) {
	if len(nodes) != expect {
		fmt.Printf("expect %d but %d, %s\n", expect, len(nodes), query)
	}
}

func TestStart(t *testing.T) {
	pid := os.Getpid()
	fmt.Println("PID : " + strconv.Itoa(pid))
}

func TestExits(t *testing.T) {
	node := LoadTestData()

	if node != nil {
		query := "exist:root->parent->child->name:jim"
		result, found, _ := Execute(node, query)
		if result != nil && found == false {
			t.Fatalf("exits test fail : %s", query)
		}

		query = "exist:root->parent->child->sex:male && name:m$"
		result, found, _ = Execute(node, query)
		if found == false {
			t.Fatalf("exits test fail : %s", query)
		}

		query = "exist:root->parent->child->sex:female && name:m$"
		result, found, _ = Execute(node, query)
		if found == true {
			t.Fatalf("exits test fail : %s", query)
		}
	}
}

func TestOne(t *testing.T) {
	node := LoadTestData()

	if node != nil {
		query := "one:root->parent->child->sex:male"
		result, _, _ := Execute(node, query)
		child, ok := result.(*Node)
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

		SelectAssert := func(query string, expect []string) {
			result, found, _ := Execute(node, query)
			if found {
				children, ok := result.([]*Node)
				if ok && len(children) == len(expect) {
					for i := 0; i < len(expect); i++ {
						child := children[i]
						if child.GetValueString("name") != expect[i] {
							t.Fatalf("select test fail : %s", query)
							break
						}
					}
				} else {
					t.Fatalf("select test fail len(children) %d != %d : %s", len(children), len(expect), query)
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
		result, found, _ := Execute(node, query)
		if found {
			updateResult, _, _ := Execute(result.(*Node), "one:root->parent->child->name:jim")
			child, ok := updateResult.(*Node)
			if ok {
				ValueAssert(child, "name", "jim", query)
				ValueAssert(child, "age", "21", query)
				ValueAssert(child, "address", "China", query)
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
		result, found, _ := Execute(node, query)
		if found {
			deleteResult, _, _ := Execute(result.(*Node), "one:root->parent->child->name:jim")
			child, ok := deleteResult.(*Node)
			if ok {
				ValueAssert(child, "name", "jim", query)
				NullAssert(child, "sex", query)
				NullAssert(child, "age", query)
			} else {
				t.Fatalf("delete test fail : %s", query)
			}
		} else {
			t.Fatalf("delete test fail : %s", query)
		}

		query = "delete:root->parent->child->name:jim !! block"
		result, found, _ = Execute(node, query)
		if found {
			deleteResult, _, _ := Execute(result.(*Node), "select:root->parent->child->name:jim")
			children, _ := deleteResult.([]*Node)
			LengthAssert(children, 0, query)
		} else {
			t.Fatalf("delete test fail : %s", query)
		}
	}
}

func TestInsert(t *testing.T) {
	node := LoadTestData()

	if node != nil {
		query := "insert:root->parent->child !! name=bill, sex=male, age=31"
		result, _, _ := Execute(node, query)

		insertResult, _, _ := Execute(result.(*Node), "one:root->parent->child->name:bill")
		child, ok := insertResult.(*Node)
		if ok {
			ValueAssert(child, "name", "bill", query)
			ValueAssert(child, "sex", "male", query)
			ValueAssert(child, "age", "31", query)
		} else {
			t.Fatalf("insert test fail : %s", query)
		}
	}
}


/*
func TestRedirect(t *testing.T) {

	node := LoadTestData()
	if node != nil {
		query := "select:root->parent->:/child|nephew/->sex:female >> select.ndb"
		Execute(node, query)
		tempNode, _ := Read("select.ndb")
		result, found, err := Execute(tempNode, "select:result->sex:female")
		if found && err == nil {
			children := result.([]*Node)

			LengthAssert(children, 2, query)
			ValueAssert(children[0], "name", "lucy", query)
			ValueAssert(children[1], "name", "lily", query)
		}

		query = "insert:root->parent->child !! name=bill, sex=male, age=31 >> insert.ndb"
		Execute(node, query)
		tempNode, _ = Read("insert.ndb")
		result, found, err = Execute(tempNode, "select:root->parent->child->name:bill")
		if found && err == nil {
			children := result.([]*Node)

			LengthAssert(children, 1, query)
			ValueAssert(children[0], "name", "bill", query)
			ValueAssert(children[0], "sex", "male", query)
			ValueAssert(children[0], "age", "31", query)
		}

		query = "update:root->parent->child->name:jim !! age=21, address=China >> update.ndb"
		Execute(node, query)
		tempNode, _ = Read("update.ndb")
		result, found, err = Execute(tempNode, "select:root->parent->child->name:jim")
		if found && err == nil {
			children := result.([]*Node)
			LengthAssert(children, 1, query)
			ValueAssert(children[0], "name", "jim", query)
			ValueAssert(children[0], "address", "China", query)
			ValueAssert(children[0], "age", "21", query)
		}

		tempFiles := []string{"select.ndb", "insert.ndb", "update.ndb"}
		for _, tempFile := range tempFiles {
			os.Remove(tempFile)
		}
	}
}
*/

func TestScript(t *testing.T) {
	node := LoadTestData()
	if node != nil {
		query := "script:d:/example.script"
		result, _, _ := Execute(node, query)
		tempNode, ok := result.(*Node)
		
		if ok {
			selectResult, found, err := Execute(tempNode, "select:root->parent->child->name:bill")
			if found && err == nil {
				children := selectResult.([]*Node)
				LengthAssert(children, 1, query)
				ValueAssert(children[0], "name", "bill", query)
				ValueAssert(children[0], "sex", "male", query)
				ValueAssert(children[0], "age", "31", query)
			}
			
			selectResult, found, err = Execute(tempNode, "select:root->parent->child->name:lily")
			if found && err == nil {
				children := selectResult.([]*Node)
				LengthAssert(children, 1, query)
				ValueAssert(children[0], "name", "lily", query)
				ValueAssert(children[0], "sex", "China", query)
				ValueAssert(children[0], "age", "21", query)
			}
			
			selectResult, found, err = Execute(tempNode, "select:root->parent->child->name:jim")
			if found && err == nil {
				children := selectResult.([]*Node)
				LengthAssert(children, 1, query)
				ValueAssert(children[0], "name", "jim", query)
				NullAssert(children[0], "sex", query)
				NullAssert(children[0], "age", query)
			}
		}
	}
}