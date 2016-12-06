package main

import (
	"context"
	"log"
	"sync"

	"github.com/twmb/algoimpl/go/graph"
)

type (

	// Category defines when in the Job lifecycle a task will execute, and what it's dependent on.
	Category int

	// Task is any executable
	Task func(ctx context.Context) context.Context

	// Job is a construct that executes once per container.
	Job struct {
		g   graph.Graph
		ctx context.Context
	}

	categoryData struct {
		c     Category
		tasks []Task
	}
)

//go:generate stringer -type=Category
const (

	// Root is the true entry point into the job.
	root Category = iota

	// PreAction is a Task that occurs before the Action
	PreAction

	// Action is the primary task performed by the job
	Action

	// PostAction is the task that occurs directly after the Action
	PostAction
)

// NewJob creates and initializes a Job with the approperate dependencies
// `payload` are the default values in the context.
func NewJob(payload map[string]interface{}) *Job {
	g := graph.New(graph.Directed)
	initializeCategoryDeps(g)

	ctx := context.Background()
	for k, v := range payload {
		ctx = context.WithValue(ctx, k, v)
	}

	return &Job{
		g:   g,
		ctx: ctx,
	}
}

// AddTask adds a new Task to the graph.
func (j *Job) AddTask(c Category, t Task) {
	for _, n := range j.g.TopologicalSort() {
		data := n.Value.(categoryData)
		if data.c == c {
			data.tasks = append(data.tasks, t)
			break
		}
	}

}

func initializeCategoryDeps(g *graph.Graph) {
	var (
		rootNode       = g.MakeNode()
		preActionNode  = g.MakeNode()
		actionNode     = g.MakeNode()
		postActionNode = g.MakeNode()
	)

	rootNode.Value = newCategoryData(root)
	preActionNode.Value = newCategoryData(PreAction)
	actionNode.Value = newCategoryData(Action)
	postActionNode.Value = newCategoryData(postActionNode)

	g.AddEdge(rootNode, preActionNode)
	g.AddEdge(preActionNode, actionNode)
	g.AddEdge(actionNode, postActionNode)
}

func newCategoryData(c Category) categoryData {
	return categoryData{
		c:     c,
		tasks: make([]Task, 0),
	}
}

// Need to figure out how to pass around an immutable context...
// And how to abstact away this wg nonsense into a decorator
// Exec runs the job
/*
func (j *Job) Exec() {
	for _, n := range j.g.TopologicalSort() {
		data := n.Value.(categoryData)
		// check to make sure that the first node is actually Root.
		if n.c != root {
			log.Fatal("Uh oh this should be Root....")
		}

	    var wg sync.WaitGroup
		for _, runnable := range n.Neighbors() {

            wg.Add(1)
            go func() {
                defer wg.Done()
                newData := runnable.Value.(categoryData)
                for _, t := range newData.tasks {
                    wg.Add(1)
                    go func() {
                        defer wg.Done()
                        t()
                    }
                }
            }
		}
	}
}
*/
