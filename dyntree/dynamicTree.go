package dyntree

import (
	. "github.com/teomat/mater/aabb"
	. "github.com/teomat/mater/vect"
	"log"
	"math"
)

type DynamicTreeNode struct {
	aabb         AABB
	child1       int
	child2       int
	leafCount    int
	parentOrNext int
	userData     interface{}
}

func (n *DynamicTreeNode) isLeaf() bool {
	return n.child1 == nullNode
}

func (n *DynamicTreeNode) AABB() AABB {
	return n.aabb
}

const nullNode = -1

type DynamicTree struct {
	_freeList       int
	_insertionCount int
	_nodeCapacity   int
	_nodeCount      int
	_path           int
	_root           int
	_nodes          []DynamicTreeNode
	_stackCount     int
	_stack          [255]int
	AABBExtension   float64
	AABBMultiplier  float64
}

func NewDynamicTree() *DynamicTree {
	dt := new(DynamicTree)
	dt._root = nullNode
	dt._nodeCapacity = 16
	dt._nodes = make([]DynamicTreeNode, dt._nodeCapacity)

	for i := 0; i < dt._nodeCapacity-1; i++ {
		dt._nodes[i].parentOrNext = i + 1
	}
	dt._nodes[dt._nodeCapacity-1].parentOrNext = nullNode

	dt.AABBExtension = 0.1
	dt.AABBMultiplier = 2.0

	return dt
}

func (dt *DynamicTree) AddProxy(aabb AABB, userData interface{}) int {
	proxyId := dt.allocateNode()

	r := Vect{dt.AABBExtension, dt.AABBExtension}
	dt._nodes[proxyId].aabb.Lower = Sub(aabb.Lower, r)
	dt._nodes[proxyId].aabb.Upper = Add(aabb.Upper, r)
	dt._nodes[proxyId].userData = userData
	dt._nodes[proxyId].leafCount = 1

	dt.insertLeaf(proxyId)

	return proxyId
}

func (dt *DynamicTree) RemoveProxy(proxyId int) {
	if proxyId < 0 || proxyId > dt._nodeCapacity {
		log.Printf("Assertion Error: Expected: 0 <= value < %v, got: %v", dt._nodeCapacity, proxyId)
	}
	if isLeaf := dt._nodes[proxyId].isLeaf(); !isLeaf {
		log.Printf("Assertion Error: Expected: value == true, got: %v", isLeaf)
	}

	dt.removeLeaf(proxyId)
	dt.freeNode(proxyId)
}

func (dt *DynamicTree) MoveProxy(proxyId int, aabb AABB, displacement Vect) bool {
	if proxyId < 0 || proxyId > dt._nodeCapacity {
		log.Printf("Assertion Error: Expected: 0 <= value < %v, got: %v", dt._nodeCapacity, proxyId)
	}
	if isLeaf := dt._nodes[proxyId].isLeaf(); !isLeaf {
		log.Printf("Assertion Error: Expected: value == true, got: %v", isLeaf)
	}

	if dt._nodes[proxyId].aabb.Contains(aabb) {
		return false
	}

	dt.removeLeaf(proxyId)

	var b AABB = aabb
	r := Vect{dt.AABBExtension, dt.AABBExtension}
	b.Lower = Sub(b.Lower, r)
	b.Upper = Add(b.Upper, r)

	d := Mult(displacement, dt.AABBMultiplier)

	if d.X < 0.0 {
		b.Lower.X += d.X
	} else {
		b.Upper.X += d.X
	}
	if d.Y < 0.0 {
		b.Lower.Y += d.Y
	} else {
		b.Upper.Y += d.Y
	}

	dt._nodes[proxyId].aabb = b

	dt.insertLeaf(proxyId)
	return true
}

func (dt *DynamicTree) Rebalance(iterations int) {
	if dt._root == nullNode {
		return
	}

	for i := 0; i < iterations; i++ {
		node := dt._root

		var bit uint = 0
		for !dt._nodes[node].isLeaf() {
			var selector int = (dt._path >> bit) & 1

			if selector == 0 {
				node = dt._nodes[node].child1
			} else {
				node = dt._nodes[node].child2
			}

			// Keep bit between 0 and 31 because _path has 32 bits
			// bit = (bit + 1) % 31
			bit = (bit + 1) & 0x1F
		}
		dt._path++

		dt.removeLeaf(node)
		dt.insertLeaf(node)
	}
}

func (dt *DynamicTree) GetUserData(proxyId int) interface{} {
	if proxyId < 0 || proxyId > dt._nodeCapacity {
		log.Printf("Assertion Error: Expected: 0 <= value < %v, got: %v", dt._nodeCapacity, proxyId)
	}
	return dt._nodes[proxyId].userData
}

func (dt *DynamicTree) GetFatAABB(proxyId int) AABB {
	if proxyId < 0 || proxyId > dt._nodeCapacity {
		log.Printf("Assertion Error: Expected: 0 <= value < %v, got: %v", dt._nodeCapacity, proxyId)
	}
	return dt._nodes[proxyId].aabb
}

func (dt *DynamicTree) ComputeHeight() int {
	return dt.computeHeight(dt._root)
}

func (dt *DynamicTree) Query(callback func(int) bool, aabb AABB) {
	//clears the stack and pushes root
	dt._stack[0] = dt._root
	dt._stackCount = 1
	for dt._stackCount > 0 {
		dt._stackCount--
		nodeId := dt._stack[dt._stackCount]
		if nodeId == nullNode {
			continue
		}

		node := dt._nodes[nodeId]
		if TestOverlap(node.aabb, aabb) {
			if node.isLeaf() {
				proceed := callback(nodeId)
				if !proceed {
					return
				}
			} else {
				dt._stack[dt._stackCount] = node.child1
				dt._stackCount++
				dt._stack[dt._stackCount] = node.child2
				dt._stackCount++
			}
		}
	}
}

type RayCastCallback func(A, B Vect, maxFraction float64, nodeId int) float64

func (dt *DynamicTree) RayCast(callback RayCastCallback, p1, p2 Vect, maxFraction float64) {
	r := Sub(p2, p1)
	if lsqr := r.LengthSqr(); lsqr <= 0.0 {
		log.Printf("Assertion Error. Expected: value > 0.0, got: %v", lsqr)
	}
	r = Normalize(r)

	// v is perpendicular to the segment.
	var absV Vect = Vect{math.Abs(-r.Y), math.Abs(r.X)}

	// Separating axis for segment (Gino, p80).
	// |dot(v, p1 - c)| > dot(|v|, h)

	// Build a bounding box for the segment.
	var segmentAABB AABB
	{
		//p1 + maxFraction * (p2 - p1)
		t := Add(p1, Mult(Sub(p2, p1), maxFraction))
		segmentAABB.Lower = Min(p1, t)
		segmentAABB.Upper = Max(p1, t)
	}

	dt._stack[0] = dt._root
	dt._stackCount = 1

	for dt._stackCount > 0 {
		dt._stackCount--
		nodeId := dt._stack[dt._stackCount]

		if nodeId == nullNode {
			continue
		}

		node := dt._nodes[nodeId]

		if !TestOverlap(node.aabb, segmentAABB) {
			continue
		}

		// Separating axis for segment (Gino, p80).
		// |dot(v, p1 - c)| > dot(|v|, h)

		var c, h Vect
		c = node.aabb.Center()
		h = node.aabb.Extents()

		// |dot(v, p1 - c)| > dot(|v|, h)
		var separation float64
		{
			t1 := math.Abs(Dot(Vect{-r.Y, r.X}, Sub(p1, c)))
			t2 := Dot(absV, h)
			separation = t1 - t2
		}

		if separation > 0.0 {
			continue
		}

		if node.isLeaf() {
			value := callback(p1, p2, maxFraction, nodeId)

			if value == 0.0 {
				// the client has terminated the raycast.
				return
			}

			if value > 0.0 {
				// Update segment bounding box.
				maxFraction = value
				t := Add(p1, Mult(Sub(p2, p1), maxFraction))
				segmentAABB.Lower = Min(p1, t)
				segmentAABB.Upper = Max(p1, t)
			}
		} else {
			dt._stack[dt._stackCount] = node.child1
			dt._stackCount++
			dt._stack[dt._stackCount] = node.child2
			dt._stackCount++
		}
	}
}

func (dt *DynamicTree) countLeaves(nodeId int) int {
	if nodeId == nullNode {
		return 0
	}

	if nodeId < 0 || nodeId > dt._nodeCapacity {
		log.Printf("Assertion Error: Expected: 0 <= value < %v, got: %v", dt._nodeCapacity, nodeId)
	}
	node := dt._nodes[nodeId]

	if node.isLeaf() {
		if node.leafCount != 1 {
			log.Printf("Assertion Error: Expected: 1, got: %v", node.leafCount)
		}
		return 1
	}

	count1 := dt.countLeaves(node.child1)
	count2 := dt.countLeaves(node.child2)

	count := count1 + count2
	if count != node.leafCount {
		log.Printf("Assertion Error: Expected: %v, got: %v", node.leafCount, count)
	}
	return count
}

func (dt *DynamicTree) validate() {
	dt.countLeaves(dt._root)
}

func (dt *DynamicTree) allocateNode() int {
	if dt._freeList == nullNode {
		//create a new slice with double the capacity
		dt._nodes = append(dt._nodes, make([]DynamicTreeNode, dt._nodeCapacity)...)
		dt._nodeCapacity *= 2

		for i := dt._nodeCount; i < dt._nodeCapacity-1; i++ {
			dt._nodes[i].parentOrNext = i + 1
		}
		dt._nodes[dt._nodeCapacity-1].parentOrNext = nullNode
		dt._freeList = dt._nodeCount
	}
	nodeId := dt._freeList
	dt._freeList = dt._nodes[nodeId].parentOrNext
	dt._nodes[nodeId].parentOrNext = nullNode
	dt._nodes[nodeId].child1 = nullNode
	dt._nodes[nodeId].child2 = nullNode
	dt._nodes[nodeId].leafCount = 0
	dt._nodeCount++

	return nodeId
}

func (dt *DynamicTree) freeNode(nodeId int) {
	if nodeId < 0 || nodeId > dt._nodeCapacity {
		log.Printf("Assertion Error: Expected: 0 <= value < %v, got: %v", dt._nodeCapacity, nodeId)
	}
	if dt._nodeCount < 0 {
		log.Printf("Assertion Error: Expected: value >= 0, got: %v", dt._nodeCount)
	}
	dt._nodes[nodeId].parentOrNext = dt._freeList
	dt._freeList = nodeId
	dt._nodeCount--
}

func (dt *DynamicTree) insertLeaf(leaf int) {
	dt._insertionCount++
	if dt._root == nullNode {
		dt._root = leaf
		dt._nodes[dt._root].parentOrNext = nullNode
		return
	}

	// Find the best sibling for this node
	var leafAABB AABB = dt._nodes[leaf].aabb
	silbling := dt._root

	for !dt._nodes[silbling].isLeaf() {
		child1 := dt._nodes[silbling].child1
		child2 := dt._nodes[silbling].child2

		// Expand the node's AABB.
		dt._nodes[silbling].aabb = Combine(dt._nodes[silbling].aabb, leafAABB)
		dt._nodes[silbling].leafCount++

		silblingArea := dt._nodes[silbling].aabb.Perimeter()
		parentAABB := Combine(dt._nodes[silbling].aabb, leafAABB)
		parentArea := parentAABB.Perimeter()
		cost1 := 2.0 * parentArea

		inheritanceCost := 2.0 * (parentArea - silblingArea)

		var cost2 float64
		if dt._nodes[child1].isLeaf() {
			aabb := Combine(leafAABB, dt._nodes[child1].aabb)
			cost2 = aabb.Perimeter() + inheritanceCost
		} else {
			aabb := Combine(leafAABB, dt._nodes[child1].aabb)
			oldArea := dt._nodes[child1].aabb.Perimeter()
			newArea := aabb.Perimeter()
			cost2 = (newArea - oldArea) + inheritanceCost
		}

		var cost3 float64
		if dt._nodes[child2].isLeaf() {
			aabb := Combine(leafAABB, dt._nodes[child2].aabb)
			cost3 = aabb.Perimeter() + inheritanceCost
		} else {
			aabb := Combine(leafAABB, dt._nodes[child2].aabb)
			oldArea := dt._nodes[child2].aabb.Perimeter()
			newArea := aabb.Perimeter()
			cost3 = (newArea - oldArea) + inheritanceCost
		}

		if cost1 < cost2 && cost1 < cost3 {
			break
		}

		dt._nodes[silbling].aabb = Combine(dt._nodes[silbling].aabb, leafAABB)

		if cost2 < cost3 {
			silbling = child1
		} else {
			silbling = child2
		}
	}

	oldParent := dt._nodes[silbling].parentOrNext
	newParent := dt.allocateNode()
	dt._nodes[newParent].parentOrNext = oldParent
	//dt._nodes[newParent].userData automatically assigned 0
	dt._nodes[newParent].aabb = Combine(leafAABB, dt._nodes[silbling].aabb)
	dt._nodes[newParent].leafCount = dt._nodes[silbling].leafCount + 1

	if oldParent != nullNode {
		if dt._nodes[oldParent].child1 == silbling {
			dt._nodes[oldParent].child1 = newParent
		} else {
			dt._nodes[oldParent].child2 = newParent
		}

		dt._nodes[newParent].child1 = silbling
		dt._nodes[newParent].child2 = leaf
		dt._nodes[silbling].parentOrNext = newParent
		dt._nodes[leaf].parentOrNext = newParent
	} else {
		dt._nodes[newParent].child1 = silbling
		dt._nodes[newParent].child2 = leaf
		dt._nodes[silbling].parentOrNext = newParent
		dt._nodes[leaf].parentOrNext = newParent
		dt._root = newParent
	}
}

func (dt *DynamicTree) removeLeaf(leaf int) {
	if leaf == dt._root {
		dt._root = nullNode
		return
	}

	parent := dt._nodes[leaf].parentOrNext
	grandParent := dt._nodes[parent].parentOrNext
	var silbling int

	if dt._nodes[parent].child1 == leaf {
		silbling = dt._nodes[parent].child2
	} else {
		silbling = dt._nodes[parent].child1
	}

	if grandParent != nullNode {
		// Destroy parent and connect sibling to grandParent.
		if dt._nodes[grandParent].child1 == parent {
			dt._nodes[grandParent].child1 = silbling
		} else {
			dt._nodes[grandParent].child2 = silbling
		}
		dt._nodes[silbling].parentOrNext = grandParent
		dt.freeNode(parent)

		// Adjust ancestor bounds.
		parent = grandParent
		for parent != nullNode {
			dt._nodes[parent].aabb = Combine(dt._nodes[dt._nodes[parent].child1].aabb,
				dt._nodes[dt._nodes[parent].child2].aabb)
			if dt._nodes[parent].leafCount <= 0 {
				log.Printf("Assertion Error: Expected: value > 0, got: %v", dt._nodes[parent].leafCount)
			}
			dt._nodes[parent].leafCount--
			parent = dt._nodes[parent].parentOrNext
		}
	} else {
		dt._root = silbling
		dt._nodes[silbling].parentOrNext = nullNode
		dt.freeNode(parent)
	}
}

func (dt *DynamicTree) computeHeight(nodeId int) int {
	if nodeId == nullNode {
		return 0
	}

	if nodeId < 0 || nodeId > dt._nodeCapacity {
		log.Printf("Assertion Error: Expected: 0 <= value < %v, got: %v", dt._nodeCapacity, nodeId)
	}

	node := dt._nodes[nodeId]
	height1 := dt.computeHeight(node.child1)
	height2 := dt.computeHeight(node.child2)
	if height1 > height2 {
		return 1 + height1
	} else {
		return 1 + height2
	}
	panic("Never reached")
}

func (dt *DynamicTree) GetNodes() []DynamicTreeNode {
	return dt._nodes
}
