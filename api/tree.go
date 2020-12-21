package api

import (
	"errors"
)

type (
	Tree struct {
		root *Node
	}

	Node struct {
		left  *Node
		value int
		right *Node
	}
)

func (tree *Tree) Insert(value int) error {
	if tree.root == nil {
		tree.root = &Node{value: value, left: nil, right: nil}
	} else {
		err := tree.root.insert(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (node *Node) insert(value int) error {
	if node == nil {
		return errors.New("node is nil")
	}

	if node.value == value {
		return errors.New("this node value already exists")
	}

	if node.value > value {
		if node.left == nil {
			node.left = &Node{value: value}
			return nil
		}
		return node.left.insert(value)
	}

	if node.value < value {
		if node.right == nil {
			node.right = &Node{value: value}
			return nil
		}
		return node.right.insert(value)
	}

	return nil
}

func (tree *Tree) Search(value int) bool {
	if tree.root == nil {
		return false
	}

	return tree.root.search(value)
}

func (node *Node) search(value int) bool {
	if node == nil {
		return false
	}

	switch {
	case value == node.value:
		return true
	case value < node.value:
		return node.left.search(value)
	default:
		return node.right.search(value)
	}
}

func (node *Node) findMax(parent *Node) (*Node, *Node) {
	if node == nil {
		return nil, parent
	}
	if node.right == nil {
		return node, parent
	}
	return node.right.findMax(node)
}

func (node *Node) replaceNode(parent, replacement *Node) error {
	if node == nil {
		return errors.New("replaceNode() not allowed on a nil node")
	}

	if node == parent.left {
		parent.left = replacement
		return nil
	}

	parent.right = replacement
	return nil
}

func (tree *Tree) Delete(value int) error {
	if tree.root == nil {
		return errors.New("can not delete from an empty tree")
	}

	parent := &Node{right: tree.root}
	err := tree.root.delete(value, parent)
	if err != nil {
		return err
	}

	if parent.right == nil {
		tree.root = nil
	}

	return nil
}

func (node *Node) delete(value int, parent *Node) error {
	if node == nil {
		return errors.New("value to be deleted does not exist in the tree")
	}

	switch {
	case value < node.value:
		return node.left.delete(value, node)
	case value > node.value:
		return node.right.delete(value, node)
	default:
		if node.left == nil && node.right == nil {
			return node.replaceNode(parent, nil)
		}
		if node.left == nil {
			return node.replaceNode(parent, node.right)
		}
		if node.right == nil {
			return node.replaceNode(parent, node.left)
		}
		replacement, replParent := node.left.findMax(node)
		node.value = replacement.value
		return replacement.delete(replacement.value, replParent)
	}
}
