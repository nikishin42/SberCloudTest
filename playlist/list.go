package playlist

import "time"

type List struct {
	start *Node
	end   *Node
}

func (l *List) Append(newNode *Node) {
	newNode.prev = l.end
	l.end.next = newNode
}

type Node struct {
	prev *Node
	next *Node
	Song
}

type Song struct {
	Name     string
	Duration time.Duration
}
