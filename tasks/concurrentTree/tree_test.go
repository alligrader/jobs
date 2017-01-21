package concurrentTree

import(
	//"fmt"
	"testing"
	//"log"
)

func TestPrint(t *testing.T){
	t.Log("woot")
}

//WHY is this map sometimes printing out the nodes in different orders?
func TestNewTree(t *testing.T){
	myTree := NewTree()
	myTree.NewNode("A", "Root", nil)
	myTree.NewNode("B", "A", nil)
	myTree.NewNode("C", "A", nil)
	myTree.NewNode("D", "C", nil)
	//myTree.NewNode("BB", "B", nil)
	for k := range myTree.Contents{
		t.Log(k)
	}
}