package concurrentTree

import(
	"context"
)
type Tree struct{
	Contents map[string] *TreeNode
	Root *TreeNode
}

type TreeNode struct{
	contents Visitable
	children []string	//the index into the Tree.Contents
	container *Tree     //a pointer to the Tree, so you can reach into the contents
}

type Visitable interface{
	Visit(context.Context) (context.Context, error)
}