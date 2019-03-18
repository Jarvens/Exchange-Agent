// date: 2019-03-14
package balance

type Node struct {
	nodeKey   string
	spotValue uint32
}

type nodesArray []Node

func (p nodesArray) Len() int {
	return len(p)
}

func (p nodesArray) Less() {

}
