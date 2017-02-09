package jobs

import (
	//"fmt"
	"context"
	"testing"
	//"time"
	//"log"
)

func TestPrint(t *testing.T) {
	t.Log("woot")
}

//WHY is this map sometimes printing out the nodes in different orders?
func TestNewTree(t *testing.T) {
	myTree := NewTree()
	myTree.NewNode("A", "Root", nil)
	myTree.NewNode("B", "A", nil)
	myTree.NewNode("C", "A", nil)
	myTree.NewNode("D", "C", nil)
	//myTree.NewNode("BB", "B", nil)
	for k := range myTree.Contents {
		t.Log(k)
	}
}

//is in tree_tester package, so this should not play with concurrentTree package

type (
	Walkable interface {
		NewNode(name, parent string, contents Visitable)
		Walk()
	}

	mockTreeNode struct {
		to   chan string
		data string
	}
)

func (m *mockTreeNode) Visit(cont context.Context) (context.Context, error) {
	//m.to is supposed to be a receive only type? but we are sending to it?
	m.to <- m.data
	//need to return a context? and an error?
	return cont, nil
}

func newMockTreeNode(to chan string, data string) *mockTreeNode {
	return &mockTreeNode{to, data}
}

func TestWalker(t *testing.T) {
	const (
		grandfather = "Grandfather"
		father      = "Father"
		mother      = "Mother"
		son         = "son"
		daughter    = "daughter"
	)

	var (
		tree            = NewTree()
		output          = make(chan string, 5)
		grandfatherNode = newMockTreeNode(output, "Grandfather")
		fatherNode      = newMockTreeNode(output, "Father")
		motherNode      = newMockTreeNode(output, "Mother")
		sonNode         = newMockTreeNode(output, "Son")
		daughterNode    = newMockTreeNode(output, "Daughter")
	)

	tree.NewNode(grandfather, "Root", grandfatherNode)
	tree.NewNode(father, grandfather, fatherNode)
	tree.NewNode(mother, grandfather, motherNode)
	tree.NewNode(son, mother, sonNode)
	tree.NewNode(daughter, mother, daughterNode)
	//	close(output)

	tree.Walk()
	close(output)

	checkbox := map[string]bool{
		grandfather: false,
		father:      false,
		mother:      false,
		son:         false,
		daughter:    false,
	}
	count := 0
	expectedCount := len(checkbox)

	for node := range output {
		count++
		checkbox[node] = true
		switch node {
		case grandfather:
			if checkbox[mother] || checkbox[father] ||
				checkbox[son] || checkbox[daughter] {
				t.Fail()
			}
		case father:
			if !checkbox[grandfather] {
				t.Fail()
			}
		case mother:
			if !checkbox[grandfather] {
				t.Fail()
			}
			if checkbox[son] || checkbox[daughter] {
				t.Fail()
			}
		case son:
			if !(checkbox[mother] && checkbox[grandfather]) {
				t.Fail()
			}

		case daughter:
			if !(checkbox[mother] && checkbox[grandfather]) {
				t.Fail()
			}
		}
	}

	if count != expectedCount {
		t.Errorf("Count was expected to be %v but was observed to be %v\n", expectedCount, count)
	}
}
