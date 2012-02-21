package collision

import (
	. "github.com/teomat/mater/aabb"
	. "github.com/teomat/mater/dyntree"
	"github.com/teomat/mater/vect"
	"math"
	"sort"
)

type pair struct {
	ProxyIdA, ProxyIdB int
}

type pairSlice []pair

func (p pairSlice) Len() int {
	return len(p)
}

func (p pairSlice) Less(i, j int) bool {
	p1, p2 := &p[i], &p[j]
	if p1.ProxyIdA < p2.ProxyIdA {
		return false
	}

	if p1.ProxyIdA == p2.ProxyIdA {
		if p1.ProxyIdB < p2.ProxyIdB {
			return false
		}
		if p1.ProxyIdB == p2.ProxyIdB {
			return true
		}
	}

	return true
}

func (p pairSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// The broad-phase is used for computing pairs and performing volume queries and ray casts.
// This broad-phase does not persist pairs. Instead, this reports potentially new pairs.
// It is up to the client to consume the new pairs and to track subsequent overlap.
type broadPhase struct {
	_moveBuffer   []int
	_moveCapacity int
	_moveCount    int

	_pairBuffer        pairSlice
	_pairCapacity      int
	_pairCount         int
	_proxyCount        int
	_queryCallbackFunc func(int) bool
	_queryProxyId      int
	_tree              *DynamicTree
}

func newBroadPhase() *broadPhase {
	dtb := new(broadPhase)
	dtb._queryCallbackFunc = func(proxyId int) bool {
		return dtb.queryCallback(proxyId)
	}

	dtb._pairCapacity = 16
	dtb._pairBuffer = make([]pair, dtb._pairCapacity)

	dtb._moveCapacity = 16
	dtb._moveBuffer = make([]int, dtb._moveCapacity)

	dtb._tree = NewDynamicTree()

	return dtb
}

func (dtb *broadPhase) proxyCount() int {
	return dtb._proxyCount
}

// Create a proxy with an initial AABB. Pairs are not reported until
// UpdatePairs is called.
func (dtb *broadPhase) addProxy(proxy shapeProxy) int {
	proxyId := dtb._tree.AddProxy(proxy.AABB, proxy)
	dtb._proxyCount++
	dtb.bufferMove(proxyId)
	return proxyId
}

// Destroy a proxy. It is up to the client to remove any pairs.
func (dtb *broadPhase) removeProxy(proxyId int) {
	dtb.unBufferMove(proxyId)
	dtb._proxyCount--
	dtb._tree.RemoveProxy(proxyId)
}

func (dtb *broadPhase) moveProxy(proxyId int, aabb AABB, displacement vect.Vect) {
	buffer := dtb._tree.MoveProxy(proxyId, aabb, displacement)
	_ = buffer
	//buffering everything for now
	//if buffer {
		dtb.bufferMove(proxyId)
	//}
}

// Get the AABB for a proxy.
func (dtb *broadPhase) getFatAABB(proxyId int) AABB {
	return dtb._tree.GetFatAABB(proxyId)
}

// Get user data from a proxy. Returns null if the id is invalid.
func (dtb *broadPhase) getProxy(proxyId int) shapeProxy {
	return dtb._tree.GetUserData(proxyId).(shapeProxy)
}

// Test overlap of fat AABBs.
func (dtb *broadPhase) testOverlap(proxyIdA, proxyIdB int) bool {
	aabbA := dtb._tree.GetFatAABB(proxyIdA)
	aabbB := dtb._tree.GetFatAABB(proxyIdB)
	return TestOverlap(aabbA, aabbB)
}

// Update the pairs. This results in pair callbacks. This can only add pairs.
func (dtb *broadPhase) updatePairs(callback func(proxyA, proxyB *shapeProxy)) {
	// Reset pair buffer
	dtb._pairCount = 0
	// Perform tree queries for all moving proxies.
	for j := 0; j < dtb._moveCount; j++ {
		dtb._queryProxyId = dtb._moveBuffer[j]
		if dtb._queryProxyId == -1 {
			continue
		}

		// We have to query the tree with the fat AABB so that
		// we don't fail to create a pair that may touch later.
		fatAABB := dtb._tree.GetFatAABB(dtb._queryProxyId)

		// Query tree, create pairs and add them pair buffer.
		dtb._tree.Query(dtb._queryCallbackFunc, fatAABB)
	}
	// Reset move buffer
	dtb._moveCount = 0

	// Sort the pair buffer to expose duplicates.
	sort.Sort(dtb._pairBuffer)

	// Send the pairs back to the client.
	i := 0
	for i < dtb._pairCount {
		primaryPair := dtb._pairBuffer[i]
		proxyA := dtb.getProxy(primaryPair.ProxyIdA)
		proxyB := dtb.getProxy(primaryPair.ProxyIdB)

		callback(&proxyA, &proxyB)
		i++

		// Skip any duplicate pairs.
		for i < dtb._pairCount {
			pair := dtb._pairBuffer[i]
			if pair.ProxyIdA != primaryPair.ProxyIdA || pair.ProxyIdB != primaryPair.ProxyIdB {
				break
			}
			i++
		}
	}

	// Try to keep the tree balanced.
	//dtb._tree.Rebalance(4)
}

// Query an AABB for overlapping proxies. The callback class
// is called for each proxy that overlaps the supplied AABB.
func (dtb *broadPhase) query(callback func(int) bool, aabb AABB) {
	dtb._tree.Query(callback, aabb)
}

// Ray-cast against the proxies in the tree. This relies on the callback
// to perform a exact ray-cast in the case were the proxy contains a shape.
// The callback also performs the any collision filtering. This has performance
// roughly equal to k * log(n), where k is the number of collisions and n is the
// number of proxies in the tree.
func (dtb *broadPhase) rayCast(callback func(a, b vect.Vect, fraction float64, proxyId int) float64, input *RayCastInput) {
	dtb._tree.RayCast(callback, input.Point1, input.Point2, input.MaxFraction)
}

func (dtb *broadPhase) touchProxy(proxyId int) {
	dtb.bufferMove(proxyId)
}

// Compute the height of the embedded tree.
func (dtb *broadPhase) computeHeight() int {
	return dtb._tree.ComputeHeight()
}

func (dtb *broadPhase) bufferMove(proxyId int) {
	if dtb._moveCount == dtb._moveCapacity {
		dtb._moveBuffer = append(dtb._moveBuffer, make([]int, dtb._moveCapacity)...)
		dtb._moveCapacity *= 2
	}

	dtb._moveBuffer[dtb._moveCount] = proxyId
	dtb._moveCount++
}

func (dtb *broadPhase) unBufferMove(proxyId int) {
	for i := 0; i < dtb._moveCount; i++ {
		if dtb._moveBuffer[i] == proxyId {
			dtb._moveBuffer[i] = -1
			return
		}
	}
}

func (dtb *broadPhase) queryCallback(proxyId int) bool {
	// A proxy cannot form a pair with itself.
	if proxyId == dtb._queryProxyId {
		return true
	}

	// Grow the pair buffer as needed.
	if dtb._pairCount == dtb._pairCapacity {
		dtb._pairBuffer = append(make([]pair, dtb._pairCapacity), dtb._pairBuffer...)
		dtb._pairCapacity *= 2
	}

	dtb._pairBuffer[dtb._pairCount].ProxyIdA = int(math.Min(float64(proxyId), float64(dtb._queryProxyId)))
	dtb._pairBuffer[dtb._pairCount].ProxyIdB = int(math.Max(float64(proxyId), float64(dtb._queryProxyId)))
	dtb._pairCount++

	return true
}
