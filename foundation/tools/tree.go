package tools

import (
	"fmt"
)

// Tree supports storage of cascading errors
// All errors that are thrown in the Data API process
// will be added to the tree, pruned for successes and failures and returned
type Tree struct {
	Data  string
	Nodes []*Tree
}

// Add inserts a parent and its child to tree
// Eg: 1->2->3 + parent:2 + child:4
// Will result in:
// 1->2->3
//     ->4
func (t *Tree) Add(parent string, child string) {
	childNode := Tree{Data: child}
	otherTree := Tree{Data: parent, Nodes: []*Tree{&childNode}}
	t.AddTree(otherTree)
}

// AddTree is the same as Add but can add an entire tree instead of just 2 nodes
// Eg: 1->2->3 + 2->3->4
// Will result in:
// 1->2->3
//     ->3->4
func (t *Tree) AddTree(other Tree) {
	if t.Data == other.Data {
		t.Nodes = append(t.Nodes, other.Nodes...)
		return
	}

	for _, node := range t.Nodes {
		if node.Data == other.Data {
			node.Nodes = append(node.Nodes, other.Nodes...)
		}
	}
}

// Depth will return the longest node chain from a current node down to the leaf
// Eg. 1->2->3->4
// Depth(1) = 3
// Depth(2) = 2
// Depth(3) = 1
// Depth(4) = 0
func (t Tree) Depth() (int, error) {

	// Initialize variable to return
	depth := 0

	// If the current node has no children, terminate recursive loop
	if len(t.Nodes) == 0 {
		return 0, nil
	}

	// Iterate through each child and get the depth
	// Increment child depth to get the depth of the current node
	var err error
	for _, node := range t.Nodes {
		depth, err = node.Depth()
		if err != nil {
			return 0, err
		}
		depth += 1
	}

	// Return depth
	return depth, nil

}

// GetFailures will prune all failures from the report tree returned from Evaluate
func (t Tree) GetFailures() ([]string, error) {

	var out []string

	if len(t.Nodes) == 0 {
		return out, nil
	}

	for _, node := range t.Nodes {
		if node.Passes() {
			continue
		}

		if len(node.Nodes) != 0 && len(node.Nodes[0].Nodes) == 0 {
			out = append(out, fmt.Sprintf("%v: %v", t.Data, node.Data))
			out = append(out, fmt.Sprintf("%v: %v", node.Data, node.Nodes[0].Data))
			continue
		}

		failures, err := node.GetFailures()
		if err != nil {
			return nil, err
		}
		if len(failures) != 0 {
			out = append(out, fmt.Sprintf("%v: %v", t.Data, node.Data))
			out = append(out, failures...)
		}
	}

	return out, nil
}

// GetSuccesses will prune all successes from the report tree returned from Evaluate
func (t Tree) GetSuccesses() ([]string, error) {

	var out []string

	if len(t.Nodes) == 0 {
		return out, nil
	}

	var successes []string
	for _, node := range t.Nodes {
		if !node.Passes() {
			continue
		}

		depth, err := node.Depth()
		if err != nil {
			return nil, err
		}

		if depth > 1 {
			successes, err = node.GetSuccesses()
			if err != nil {
				return nil, err
			}
			out = append(out, fmt.Sprintf("%v: %v", t.Data, node.Data))
			out = append(out, successes...)
		}

	}

	return out, nil
}

// Passes returns true if every Node does not contain an error
func (t Tree) Passes() bool {
	passedNodes := 0
	if t.Data == "[Pass()]" {
		return true
	}
	for _, node := range t.Nodes {
		if node.Data == "[Pass()]" {
			passedNodes++
			continue
		}
		if node.Passes() {
			passedNodes++
		}
	}

	return passedNodes == len(t.Nodes) && passedNodes != 0
}
