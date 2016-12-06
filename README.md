# jobs
This repository contains helpful functions and types for specifying autograder jobs.

# Golang Graph Format for Tasks

A Task is anything scheduled to be performed by a job.
We categorized Tasks into three categories: PreAction, Action, and PostAction. Consider a generalization of this model, where the job is represented as a graph of tasks, where two tasks are connected if they are dependent on each other.

    PreAction -> Action -> PostAction

Now, we add an invisible root action Root, just for the sake of the demostration. Now, if we connect Root to PreAction, performing a topoligical sort on Root will always return a parallelizable list of tasks which are not dependent on each level.

    Root -> PreAction -> Action -> PostAction.

Then, in case we need to add initialization or error handling or whatever, then we can insert these anywhere in the graph and perform them in parallel.


Now, each node in the graph holds a slice of tasks registered under that particular task category. For example, there could be two PreActions registered. When the Root resolves, it gets all of it's neighbors with a call to `Neighbors()`. Then, Each of the neighbors executes in parallel. A Goroutine will launch for each of the neighbors, and then each of the tasks in the slice will execute in it's own Goroutine. Of course, we need some way to pass the results of a prior task into the children tasks. For that, we pass along a Context object.

Thus, we're going to have a Tasks directory that defines all of the TaskCategories.

    type TaskCategory int
    const (
        Root TaskCategory = iota
        PreAction
        Action
        PostAction
    )

Then, we add a node to the graph for the Root.

    dependencies Graph = new Graph()
    dependencies.addNode()
