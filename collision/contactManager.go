package collision

//doubly linked list of arbiters
type ArbiterList struct {
	Arbiter *Arbiter
}

//new arbiters are inserted at the front of the list
func (arbList *ArbiterList) Add(arb *Arbiter) {
	if arbList.Arbiter != nil {
		arbList.Arbiter.Prev = arb
		arb.Next = arbList.Arbiter
	}
	arbList.Arbiter = arb
}

func (arbList *ArbiterList) Remove(arb *Arbiter) {
	if arb.Prev != nil {
		arb.Prev.Next = arb.Next
	}
	if arb.Next != nil {
		arb.Next.Prev = arb.Prev
	}

	if arbList.Arbiter == arb {
		if arb.Prev != nil {
			arbList.Arbiter = arb.Prev
		} else {
			arbList.Arbiter = arb.Next
		}
	}
}

type ContactManager struct {
	BroadPhase  *BroadPhase
	ArbiterList ArbiterList
	Space       *Space
}

func newContactManager(space *Space) *ContactManager {
	cm := new(ContactManager)
	cm.Space = space
	cm.BroadPhase = space.BroadPhase

	return cm
}

func (cm *ContactManager) addPair(proxyA, proxyB *shapeProxy) {
	shapeA := proxyA.Shape
	shapeB := proxyB.Shape
	bodyA := shapeA.Body
	bodyB := shapeB.Body

	// Are the shapes on the same body?
	if bodyA == bodyB {
		return
	}

	// Does a arbiter already exist?
	edge := bodyB.arbiterList
	for edge != nil {
		if edge.Other == bodyA {
			sA := edge.Arbiter.ShapeA
			sB := edge.Arbiter.ShapeB

			if sA == shapeA && sB == shapeB {
				//arbiter exists
				return
			}
			if sA == shapeB && sB == shapeA {
				//arbiter exists
				return
			}
		}
		edge = edge.Next
	}

	// Does a joint override collision? Is at least one body dynamic?
	if !bodyB.shouldCollide(bodyA) {
		return
	}

	if !shouldCollide(shapeA, shapeB) {
		return
	}

	// Check user filtering.
	arbiterFilter := cm.Space.Callbacks.ShouldCollide
	if arbiterFilter != nil && arbiterFilter(shapeA, shapeB) == false {
		return
	}

	// Call the factory.
	arb := CreateArbiter(shapeA, shapeB)

	// Contact creation may swap shapes.
	shapeA = arb.ShapeA
	shapeB = arb.ShapeB
	bodyA = shapeA.Body
	bodyB = shapeB.Body

	// Insert into the world.
	cm.ArbiterList.Add(arb)

	// Connect to island graph.

	// Connect to body A
	arb.nodeA.Arbiter = arb
	arb.nodeA.Other = bodyB
	arb.nodeA.Prev = nil
	arb.nodeA.Next = bodyA.arbiterList

	if bodyA.arbiterList != nil {
		bodyA.arbiterList.Prev = arb.nodeA
	}
	bodyA.arbiterList = arb.nodeA

	// Connect to body B
	arb.nodeB.Arbiter = arb
	arb.nodeB.Other = bodyA
	arb.nodeB.Prev = nil
	arb.nodeB.Next = bodyB.arbiterList

	if bodyB.arbiterList != nil {
		bodyB.arbiterList.Prev = arb.nodeB
	}
	bodyB.arbiterList = arb.nodeB
}

func (cm *ContactManager) findNewContacts() {
	cm.BroadPhase.updatePairs(
		func(proxyA, proxyB *shapeProxy) {
			cm.addPair(proxyA, proxyB)
		})
}

func (cm *ContactManager) destroy(arbiter *Arbiter) {
	sA := arbiter.ShapeA
	sB := arbiter.ShapeB
	bodyA := sA.Body
	bodyB := sB.Body

	//remove the arbiter from our list
	cm.ArbiterList.Remove(arbiter)

	// Remove from body 1
	if arbiter.nodeA.Prev != nil {
		arbiter.nodeA.Prev.Next = arbiter.nodeA.Next
	}

	if arbiter.nodeA.Next != nil {
		arbiter.nodeA.Next.Prev = arbiter.nodeA.Prev
	}

	if arbiter.nodeA == bodyA.arbiterList {
		bodyA.arbiterList = arbiter.nodeA.Next
	}

	// Remove from body 2
	if arbiter.nodeB.Prev != nil {
		arbiter.nodeB.Prev.Next = arbiter.nodeB.Next
	}

	if arbiter.nodeB.Next != nil {
		arbiter.nodeB.Next.Prev = arbiter.nodeB.Prev
	}

	if arbiter.nodeB == bodyB.arbiterList {
		bodyB.arbiterList = arbiter.nodeB.Next
	}

	arbiter.destroy()
}

func (cm *ContactManager) collide() {
	for arb := cm.ArbiterList.Arbiter; arb != nil; arb = arb.Next {
		shapeA := arb.ShapeA
		shapeB := arb.ShapeB
		bodyA := shapeA.Body
		bodyB := shapeB.Body

		//if !bodyA.Awake() && !bodyB.Awake() {continue}

		if !bodyB.shouldCollide(bodyA) {
			cm.destroy(arb)
			continue
		}

		if !shouldCollide(shapeA, shapeB) {
			cm.destroy(arb)
			continue
		}

		arbiterFilter := cm.Space.Callbacks.ShouldCollide
		if arbiterFilter != nil && arbiterFilter(shapeA, shapeB) == false {
			cm.destroy(arb)
			return
		}

		proxyIdA := shapeA.proxy.ProxyId
		proxyIdB := shapeB.proxy.ProxyId
		overlap := cm.BroadPhase.testOverlap(proxyIdA, proxyIdB)
		// Here we destroy arbiters that cease to overlap in the broad-phase.
		if !overlap {
			cm.destroy(arb)
			continue
		}

		arb.update()

		if arb.NumContacts <= 0 {
			cm.destroy(arb)
		} else {
			collisionCallback := cm.Space.Callbacks.OnCollision
			if collisionCallback != nil {
				collisionCallback(arb)
			}
		}

	}
}

func shouldCollide(sA, sB *Shape) bool {
	return true
}
