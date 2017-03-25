package router

import "fmt"

// node的类别
const (
	normal = iota
	pathReg
	pathVar
	matchAll
)

// GetRoot get a root node
// this node start with '/'
func GetRoot() *Node {
	return &Node{class: normal}
}

// ReserveType 转化储存的实体类型
// 一面后面修改麻烦
type ReserveType RestHandler

// conflictRTDealFunc 注册时冲突处理函数type
type conflictRTDealFunc func(old, new ReserveType) ReserveType

var (
	crtdf conflictRTDealFunc = func(old, new ReserveType) ReserveType {
		old.Joint(new)
		return old
	}
)

// Node 装载url每一个由`/`隔开的分段的实体
// 递归结构
type Node struct {
	class    int
	value    nodeV
	handler  RestHandler
	children *childrens
}

// InsertChild 插入一个子node到一个node中
// 递归插入
// 直到url到最后
func (n *Node) InsertChild(URL string, rt ReserveType) {
	n.insert(URL, rt)
}

func (n *Node) insert(URL string, rt ReserveType) {
	prefix, other := splitURL(URL)
	if prefix == "" {
		if n.handler == nil {
			n.handler = rt
		}
		n.handler = crtdf(n.handler, rt)
		return
	}
	if n.children == nil {
		n.children = new(childrens)
	}
	node := n.children.getNode(prefix, nil)
	if node == nil {
		node = creatNode(prefix)
	}
	node.insert(other, rt)
	n.children.insert(node)

}

// creatNode create a new Node
func creatNode(passageURL string) *Node {
	return &Node{
		class:    getClass(passageURL),
		value:    createNodeValue(passageURL),
		children: new(childrens),
	}
}

// GetRT get reserve type
// if reserve type value is nil , return error
func (n *Node) GetRT(url string, paramMap map[string]string) (ReserveType, error) {
	isMatch := false
	node := n.getNode(url, paramMap, &isMatch)
	if !isMatch {
		return nil, NotMatch{"this url not rt"}
	}
	return node.handler, nil
}

// getNode 递归获取 匹配的node
// 用户获取路由Node
// 与下方的GetNode不同，GetNode用于获取注册Node
func (n *Node) getNode(url string, paramMap map[string]string, isMatch *bool) (node *Node) {

	// 设置node的值并返回
	var setNode = func(n *Node) {
		*isMatch = true
		node = n
	}

	prefix, other := splitURL(url)
	if prefix == "" {
		setNode(n)
		return
	}
	childrens := n.children
	defer func() {
		if childrens.suffixMatch != nil && !*isMatch {
			addValueToPathParam(paramMap, suffixMatchSymbol, url)
			setNode(childrens.suffixMatch)
		}
	}()
	node = childrens.getNode(prefix, paramMap)
	if node != nil {
		if other == "" {
			setNode(node)
		} else {
			node = node.getNode(other, paramMap, isMatch)
		}
	}
	return
}

// GetNode 根据url 获取一个注册Node
func (n *Node) GetNode(url string) *Node {
	if url == "" {
		return n
	}
	_, err := n.GetRT(url, nil)
	if err != nil {
		n.InsertChild(url, nil)
	}
	_true := true
	return n.getNode(url, nil, &_true)
}

// String 打印内容
func (n *Node) String() string {
	return fmt.Sprintf("{ class:%d,value:%v,handler:%v,children:%v }", n.class, n.value, n.handler, n.children)
}

// childrens 用于封装该node的所有子node
// 不同的类型使用不同的储存方式
// 以提高性能
type childrens struct {
	normal   map[string]*Node
	regexp   []*Node
	allMatch *Node
	// 尾部全匹配以栈的形式随方法存入
	// 当最后没有匹配是，将获取栈中的倒数第一个元素放回
	suffixMatch *Node
}

// getNode 根据以`/`分段的url获取子Node
// passageURL 一段URL
func (c *childrens) getNode(passageURL string, paramMap map[string]string) (node *Node) {
	ok := false
	var v *Node
	if c.normal != nil {
		v, ok = c.normal[passageURL]
	}
	if ok {
		node = v
	} else if c.regexp != nil {
		for i := 0; i < len(c.regexp); i++ {
			nowNode := c.regexp[i]
			value := nowNode.value
			if match := value.isMatch(passageURL); match {
				node = nowNode
				value.paramValue(passageURL, paramMap)
				break
			}
		}
	} else if c.allMatch != nil {
		node = c.allMatch
		c.allMatch.value.paramValue(passageURL, paramMap)
	}
	return
}

func (c *childrens) insert(node *Node) {
	switch node.class {
	case normal:
		if c.normal == nil {
			c.normal = make(map[string]*Node)
		}
		c.normal[node.value.getOrigin()] = node
	case pathReg:
		if c.regexp == nil {
			c.regexp = make([]*Node, 0, 3)
		}
		c.regexp = append(c.regexp, node)
	case pathVar:
		c.allMatch = node
	case matchAll:
		c.suffixMatch = node
	}
}

func (c *childrens) String() string {
	return fmt.Sprintf("{ normal:%v,regexp:%v,allMatch:%v,suffixMatch:%v }", c.normal, c.regexp, c.allMatch, c.suffixMatch)
}
