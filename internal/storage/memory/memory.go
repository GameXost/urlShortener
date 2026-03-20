package memory

import (
	"context"
	"sync"
	"url_shortener/internal/storage"
)

type Cache struct {
	aliasMap map[string]*Node
	urlMap   map[string]*Node
	mu       sync.RWMutex
	size     int
	capacity int
	head     *Node
	tail     *Node
}

type Node struct {
	alias string
	url   string
	prev  *Node
	next  *Node
}

func NewCache(capacity int) *Cache {
	return &Cache{
		aliasMap: make(map[string]*Node, capacity), // alias - ключ
		urlMap:   make(map[string]*Node, capacity), // url - ключ
		size:     0,
		capacity: capacity,
		head:     nil,
		tail:     nil,
	}
}

func (c *Cache) GetURL(ctx context.Context, alias string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	node, has := c.aliasMap[alias]
	if !has {
		return "", storage.ErrURLNotFound
	}
	c.moveToTop(node)
	return node.url, nil
}

func (c *Cache) SaveURL(ctx context.Context, url, alias string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, has := c.urlMap[url]
	if has {
		c.moveToTop(node)
		return node.alias, nil
	}
	if _, has = c.aliasMap[alias]; has {
		return "", storage.ErrCollision
	}

	node = &Node{
		alias: alias,
		url:   url,
		prev:  nil,
		next:  nil,
	}

	if c.size >= c.capacity {
		c.deleteBottom()
	}

	c.aliasMap[alias] = node
	c.urlMap[url] = node
	c.addToFront(node)
	return node.alias, nil
}

func (c *Cache) addToFront(node *Node) {
	node.next = c.head
	node.prev = nil

	if c.head != nil {
		c.head.prev = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
	c.size++
}

func (c *Cache) deleteBottom() {
	if c.tail == nil {
		return
	}

	node := c.tail
	if c.head == c.tail {
		c.tail = nil
		c.head = nil
	} else {
		c.tail = c.tail.prev
		c.tail.next = nil
	}
	delete(c.urlMap, node.url)
	delete(c.aliasMap, node.alias)
	c.size--
}

func (c *Cache) moveToTop(node *Node) {
	if node == c.head {
		return
	}
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = c.tail.prev
	}

	node.prev = nil
	node.next = c.head
	if c.head != nil {
		c.head.prev = node
	}
	c.head = node
}
