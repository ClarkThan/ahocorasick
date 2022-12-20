// Package ahocorasick provides the golang implementation of the Aho Corasick multiple pattern match algorithm
package ahocorasick

type Matcher struct {
	root  *Node
	ready bool // indicating build is completed
}

type Node struct {
	child map[rune]*Node
	fail  *Node
	ends  []int // length of those pattern strings which ended with this node
}

// A Hit records matched pattern string's start index in the searched string,
// and length of matched pattern string.
type Hit struct {
	Start, Len int
}

// A fifoQueue aims to reduce memory allocating - circular queue
type fifoQueue struct {
	head, cnt int
	data      []*Node
}

func newQueue() *fifoQueue {
	return &fifoQueue{
		data: make([]*Node, 8),
	}
}

func (q *fifoQueue) Push(n *Node) {
	capacity := cap(q.data) // underlying array length

	if q.cnt < capacity {
		q.data[(q.head+q.cnt)%capacity] = n
		q.cnt++
	} else {
		// current capacity is not enough, allocate new length array
		dat := make([]*Node, 2*(capacity+1))
		for i := 0; i < q.cnt; i++ {
			dat[i] = q.data[(q.head+i)%capacity]
		}
		dat[q.cnt] = n
		q.cnt++
		q.data = dat
		q.head = 0
	}
}

func (q *fifoQueue) pop() *Node {
	if q.cnt == 0 {
		return nil
	}

	capacity := cap(q.data)
	ret := q.data[q.head%capacity]
	q.head = (q.head + 1) % capacity
	q.cnt--
	return ret
}

func makeNode() *Node {
	return &Node{
		child: make(map[rune]*Node),
	}
}

func NewMatcher() *Matcher {
	return &Matcher{
		root: makeNode(),
	}
}

// doAddPattern handles the trie building with the input pattern string
func (m *Matcher) doAddPattern(pat string) {
	node := m.root
	len := 0
	for _, chr := range pat {
		nn, exists := node.child[chr]
		if !exists {
			nn = makeNode()
			node.child[chr] = nn
		}
		node = nn
		len++
	}

	// if existed then skip
	for _, e := range node.ends {
		if e == len {
			return
		}
	}
	node.ends = append(node.ends, len)
}

func (m *Matcher) AddPattern(pat string) {
	m.doAddPattern(pat)
	m.ready = false
}

// Build handles the fail pointer building after trie built
func (m *Matcher) Build() {
	m.ready = true
	if m.root == nil || len(m.root.child) == 0 {
		return
	}

	q := newQueue()
	for n := m.root; n != nil; n = q.pop() {
		for c, v := range n.child {
			if n == m.root {
				v.fail = m.root // fail pointer of first char of pattern always has to be the root
			} else {
				fatherFail := n.fail
				for ; fatherFail != nil; fatherFail = fatherFail.fail {
					if k, exists := fatherFail.child[c]; exists {
						v.fail = k
						// if k.ends's element has existed in v.ends then skip
						for _, ke := range k.ends {
							exists := false
							for _, ve := range v.ends {
								if ke == ve {
									exists = true
									break
								}
							}
							if !exists {
								v.ends = append(v.ends, ke)
							}
						}
						break
					}
				}

				if fatherFail == nil { // backtrace to the root, still not find the char
					v.fail = m.root
				}
			}

			q.Push(v)
		}
	}
}

// BuildWithPatterns handles the trie and fail pointer building
func (m *Matcher) BuildWithPatterns(patterns []string) {
	for i := range patterns {
		m.doAddPattern(patterns[i])
	}
	m.Build()
}

func (m *Matcher) check() {
	if !m.ready {
		panic("you should use `Build() or BuildWithPatterns()` before searching")
	}
}

// SearchIndexed return start index in the searched string and length of the matched pattern strings
func (m *Matcher) SearchIndexed(s string) (ret []Hit) {
	m.check()
	node := m.root
	for i, c := range []rune(s) {
		for node != nil {
			n, exists := node.child[c]
			if !exists {
				node = node.fail // try to find at its fail pointer node
			} else {
				node = n
				for _, e := range n.ends {
					ret = append(ret, Hit{Start: i - e + 1, Len: e})
				}
				break // hit a char, then find next char
			}
		}

		// walk to the root, still not find
		if node == nil {
			node = m.root
		}
	}

	return
}

// Search return the matched pattern strings
func (m *Matcher) Search(s string) (ret []string) {
	m.check()
	node := m.root
	chars := []rune(s)
	for i, c := range chars {
		for node != nil {
			n, exists := node.child[c]
			if !exists {
				node = node.fail
			} else {
				node = n
				for _, e := range n.ends {
					ret = append(ret, string(chars[(i-e+1):i+1]))
				}
				break
			}
		}

		if node == nil {
			node = m.root
		}
	}

	return
}

// Match return true if does matched
func (m *Matcher) Match(s string) bool {
	m.check()
	node := m.root
	for _, c := range s {
		for node != nil {
			n, exists := node.child[c]
			if !exists {
				node = node.fail
			} else {
				node = n
				if len(n.ends) > 0 {
					return true
				}
				break
			}
		}

		if node == nil {
			node = m.root
		}
	}

	return false
}
