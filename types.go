package jobs

import (
	"context"
)

type (
	Job interface {
		Execute()
	}

	Tree struct {
		Contents map[string]*TreeNode
		Root     *TreeNode
	}

	TreeNode struct {
		contents  Visitable
		children  []string //the index into the Tree.Contents
		container *Tree    //a pointer to the Tree, so you can reach into the contents
	}

	Visitable interface {
		Visit(context.Context) (context.Context, error)
	}

	funcVisitable func(context.Context) (context.Context, error)
)
