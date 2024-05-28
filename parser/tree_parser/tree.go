package tree_parser

import (
	"fmt"
	"log"
	"slices"
)

type TreeNode struct {
	Label   string
	IsClose bool
	Value   Value

	children    []*TreeNode
	hasChildren bool
}

func NewTreeNode(label string, isClose bool, values ...*Token) *TreeNode {
	if len(values) > 0 {
		values = slices.Insert(values, 0, NewDelimiterToken("\n\t"))
		values = append(values, NewDelimiterToken("\n\n"))
	}
	n := &TreeNode{
		Label:    label,
		IsClose:  isClose,
		Value:    values,
		children: make([]*TreeNode, 0),
	}
	n.ParseChildren()
	return n
}

func (n *TreeNode) IsLeaf() bool {
	return !n.hasChildren
}

func (n *TreeNode) String() string {
	ret := fmt.Sprintf("[%s]", n.Label)
	for _, c := range n.children {
		ret += c.String()
	}
	if n.IsClose {
		ret += fmt.Sprintf("[/%s]", n.Label)
	}
	return ret
}

func (n *TreeNode) ParseChildren() {
	closedKeys := make(map[string]bool)
	for _, t := range n.Value {
		if t.IsCloseKey() {
			closedKeys[t.content[1:]] = true
		}
	}
	n.parseChildren(closedKeys)
}

func (n *TreeNode) parseChildren(closedKeys map[string]bool) {
	for _, t := range n.Value {
		if t.IsDelimiter() || t.IsIgnore() {
			continue
		}
		if t.IsKey() {
			n.hasChildren = true
			break
		} else {
			n.hasChildren = false
			return
		}
	}

	nodes := make([]*TreeNode, 0)
	var node *TreeNode
	keyDepth := 0 // 相同key的深度
	for _, t := range n.Value {
		switch t.tp {
		case TokenKey:
			if node == nil {
				node = &TreeNode{
					Label:   t.content,
					IsClose: closedKeys[t.content],
					Value:   make(Value, 0),
				}
				if node.IsClose {
					keyDepth++
				}
			} else {
				if node.IsClose {
					if t.IsCloseKeyBy(node.Label) {
						keyDepth--
						if keyDepth == 0 {
							node.parseChildren(closedKeys)
							nodes = append(nodes, node)
							node = nil
						}
					} else {
						node.Value = append(node.Value, t)
					}
				} else {
					node.parseChildren(closedKeys)
					nodes = append(nodes, node)
					node = &TreeNode{
						Label:   t.content,
						IsClose: closedKeys[t.content],
						Value:   make(Value, 0),
					}
					if node.IsClose {
						keyDepth++
					}
				}
			}
		default:
			if node != nil {
				node.Value = append(node.Value, t)
			}
		}
	}
	if node != nil {
		node.parseChildren(closedKeys)
		nodes = append(nodes, node)
	}
	n.children = nodes
}

func (n *TreeNode) Render() string {
	if n.Label == "root" {
		c := ""
		for _, child := range n.children {
			c += child.Render()
		}
		return c
	}

	c := fmt.Sprintf("[%s]", n.Label)
	if n.hasChildren {
		c += "\n\t"
		for _, child := range n.children {
			c += child.Render()
		}
	} else {
		for _, v := range n.Value {
			c += v.Render()
		}
	}
	if n.IsClose {
		c += fmt.Sprintf("[/%s]\n\n", n.Label)
	}
	return c
}

func (n *TreeNode) GetChildren(label string) []*TreeNode {
	ret := make([]*TreeNode, 0)
	for _, child := range n.children {
		if child.Label == label {
			ret = append(ret, child)
		}
	}
	return ret
}

func (n *TreeNode) GetFirstChild(label string) *TreeNode {
	for _, child := range n.children {
		if child.Label == label {
			return child
		}
	}
	return nil
}

func (n *TreeNode) AddChild(node *TreeNode) {
	if !n.hasChildren {
		if len(n.Value) > 0 {
			log.Printf("add node err, parent %s has value already.", n.Label)
			return
		}
		n.hasChildren = true
	}
	n.children = append(n.children, node)
}

func (n *TreeNode) SetChildren(nodes []*TreeNode) {
	if !n.hasChildren {
		if len(n.Value) > 0 {
			log.Printf("add node err, parent %s has value already.", n.Label)
			return
		}
		n.hasChildren = true
	}
	n.children = nodes
}
