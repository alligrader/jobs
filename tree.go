package jobs

import (
	"context"
	"log"
	"sync"
)

func (tree *Tree) Walk() {
	var done sync.WaitGroup
	done.Add(1) // we have to visit the Root first.
	log.Println("Beginning to walk the tree.")

	go func() {
		// make the parent context
		ctx := context.Background()
		tree.Root.visitAllChildren(ctx, done)

		done.Done()
	}()

	done.Wait()
}

func (node *TreeNode) Walk(ctx context.Context, done sync.WaitGroup) {
	// tell your parent that you're done
	defer done.Done()

	// visit the contents
	childCtx, err := node.contents.Visit(ctx)
	if err != nil {
		// if there was an error, report it
		log.Fatal(err)
	}

	// visit each child.
	node.visitAllChildren(childCtx, done)
}

func (node *TreeNode) visitAllChildren(ctx context.Context, wg sync.WaitGroup) {
	wg.Add(len(node.children))
	for _, child := range node.children {
		go func() {
			log.Printf("Walking child %v\n", child)
			tree := node.container.Contents
			tree[child].Walk(ctx, wg)
		}()
	}
}

// NewTree creates a new, empty tree.
func NewTree() *Tree {
	t := &Tree{}
	t.Root = &TreeNode{
		contents:  nil,
		children:  []string{},
		container: t,
	}
	t.Contents = map[string]*TreeNode{
		"Root": t.Root,
	}

	return t
}

func (tree *Tree) NewNode(name string, parent string, contents Visitable) {
	if _, ok := tree.Contents[parent]; !ok {
		panic("No such parent.")
	}
	newChild := &TreeNode{
		contents:  contents,
		children:  []string{},
		container: tree,
	}
	//tree.Contents[parent].children = append(tree.Contents[parent].children, newChild)
	//tree.Contents[parent] = append(tree.Contents[parent], newChild) //appends newNode to the tree .... i imagine this is wanted?
	tree.Contents[name] = newChild
	tree.Contents[parent].children = append(tree.Contents[parent].children, name) //adds the name of the newnode to the children array of the parent
}

func (t *Tree) Execute() {
	t.Walk()
}

func (f funcVisitable) Visit(ctx context.Context) (context.Context, error) {
	return f(ctx)
}
