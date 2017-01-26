package tree_tester

import (
	"context"
	"testing"
)

type (
	Visitable interface {
		Visit(context.Context) (context.Context, error)
	}

	Walkable interface {
		NewNode(name, parent string, contents Visitable)
		Walk()
	}

	mockTreeNode struct {
		to   <-chan string
		data string
	}
)

func (m *mockTreeNode) Visit(context.Context) (context.Context, error) {
	m.to <- m.data
}

func newMockTreeNode(to chan string, data string) *mockTreeNode {
	return &mockTreeNode{to, data}
}

func TestWalker(t *testing.T, tree Walkable) {

	const (
		grandfather = "Grandfather"
		father      = "Father"
		mother      = "Mother"
		son         = "son"
		daughter    = "daughter"
	)

	var (
		output          = make(chan string)
		grandfatherNode = newMockTreeNode(output, "Grandfather")
		fatherNode      = newMockTreeNode(output, "Father")
		motherNode      = newMockTreeNode(output, "Mother")
		sonNode         = newMockTreeNode(output, "Son")
		daughter        = newMockTreeNode(output, "Daughter")
	)

	tree.NewNode(grandfather, "Root", grandfatherNode)
	tree.NewNode(father, grandfather, fatherNode)
	tree.NewNode(mother, grandfather, motherNode)
	tree.NewNode(son, mother, sonNode)
	tree.NewNode(daughter, mother, daughterNode)
	close(output)

	go tree.Walk()

	checkbox := map[string]bool{
		grandfather: false,
		father:      false,
		mother:      false,
		son:         false,
		daughter:    false,
	}

	for node := range output {
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
}
