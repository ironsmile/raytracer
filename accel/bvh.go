package accel

import (
	"fmt"
	"math"
	"sort"

	"github.com/ironsmile/raytracer/bbox"
	"github.com/ironsmile/raytracer/geometry"
	"github.com/ironsmile/raytracer/primitive"
)

// BVH is an approach for ray intersection acceleration based on primitive subdivision,
// where the primitives are partitioned into a hierarchy of disjoint sets.
// BVH stands for Bounding Volume Hierarchies.
type BVH struct {
	Base

	maxPrimsInNode int
	nodes          []linearBVHNode
	primitives     []primitive.Primitive
}

// bvhPrimitiveInfo is a struct used for building the traverse tree in BVH
type bvhPrimitiveInfo struct {
	centroid        geometry.Vector
	primitiveNumber int
	bounds          *bbox.BBox
}

// bvhBuildNode represents a node in the BVH. All nodes store a BBox , which stores
// the bounds of all of the children beneath the node. Each interior node stores pointers to
// its two children in children . Interior nodes also record the coordinate axis along which
// primitives were sorted for distribution to their two children; this information is used to
// improve the performance of the traversal algorithm. Leaf nodes need to record which
// primitive or primitives are stored in them.
type bvhBuildNode struct {
	bounds          *bbox.BBox
	childern        [2]*bvhBuildNode
	splitAxis       int
	firstPrimOffset int
	nPrimitives     int
}

func (bn *bvhBuildNode) InitLeaf(first, n int, bb *bbox.BBox) {
	bn.firstPrimOffset = first
	bn.nPrimitives = n
	bn.bounds = bb
}

func (bn *bvhBuildNode) InitInterior(axis int, c0, c1 *bvhBuildNode) {
	bn.childern[0] = c0
	bn.childern[1] = c1
	bn.bounds = bbox.Union(c0.bounds, c1.bounds)
	bn.splitAxis = axis
	bn.nPrimitives = 0
}

// linearBVHNode represents a node in the compact tree of bvh
type linearBVHNode struct {
	bounds      bbox.BBox
	offset      uint32
	nPrimitives uint8
	axis        uint8
}

// NewBVH returns a new BVH structure which would accelerate the intersection of the
// primitives `p`. The mp arguments is the number of primitives that can be in any leaf node.
func NewBVH(p []primitive.Primitive, mp uint8) *BVH {

	bvh := &BVH{
		maxPrimsInNode: int(math.Min(float64(mp), 255.0)),
	}

	bvh.primitives = FullyRefinePrimitives(p)

	// Nothing else to do, this would be an empty BVH
	if len(bvh.primitives) == 0 {
		return bvh
	}

	// Building the BVH tree from its primitives
	var buildData []bvhPrimitiveInfo

	for i, prim := range bvh.primitives {
		bb := prim.GetWorldBBox()
		buildData = append(buildData, newBVHPrimitiveInfo(i, bb))
	}

	//WIP
	var totalNodes int
	root, orderedPrims := bvh.bvhRecursiveBuild(buildData, &totalNodes, nil)
	bvh.primitives = orderedPrims

	bvh.nodes = make([]linearBVHNode, totalNodes)

	var offset uint32
	bvh.flattenBVHTree(root, &offset)

	fmt.Printf("Final BVH has %d nodes\n", totalNodes)

	return bvh
}

func (bvh *BVH) flattenBVHTree(node *bvhBuildNode, offset *uint32) uint32 {
	linearNode := &bvh.nodes[*offset]
	linearNode.bounds = *node.bounds
	myOffset := *offset
	(*offset)++

	if node.nPrimitives > 0 {
		linearNode.offset = uint32(node.firstPrimOffset)
		linearNode.nPrimitives = uint8(node.nPrimitives)
	} else {
		linearNode.axis = uint8(node.splitAxis)
		linearNode.nPrimitives = 0
		bvh.flattenBVHTree(node.childern[0], offset)
		linearNode.offset = bvh.flattenBVHTree(node.childern[1], offset)
	}

	return myOffset
}

func (bvh *BVH) bvhRecursiveBuild(
	buildData []bvhPrimitiveInfo,
	totalNodes *int,
	orderedPrims []primitive.Primitive,
) (*bvhBuildNode, []primitive.Primitive) {
	(*totalNodes)++

	var centroidBound *bbox.BBox
	var start, dim, mid int
	end := len(buildData)

	node := &bvhBuildNode{}
	// bounds of all primitives in bvh node
	var bb *bbox.BBox
	for i := start; i < end; i++ {
		bb = bbox.Union(bb, buildData[i].bounds)
	}

	nPrimitives := end - start
	if nPrimitives == 1 {
		// Create leaf bvhBuildNode
		firstPrimOffset := len(orderedPrims)
		for i := start; i < end; i++ {
			primNum := buildData[i].primitiveNumber
			orderedPrims = append(orderedPrims, bvh.primitives[primNum])
		}
		node.InitLeaf(firstPrimOffset, nPrimitives, bb)
		return node, orderedPrims
	}

	// build of primitive centroid and dim
	for i := start; i < end; i++ {
		centroidBound = bbox.UnionPoint(centroidBound, buildData[i].centroid)
	}
	dim = centroidBound.MaximumExtend()
	// partition primitives in two sets and build children
	mid = (start + end) / 2
	if centroidBound.Max.ByAxis(dim) == centroidBound.Min.ByAxis(dim) {

		if nPrimitives <= bvh.maxPrimsInNode {
			firstPrimOffset := len(orderedPrims)
			for i := start; i < end; i++ {
				primNum := buildData[i].primitiveNumber
				orderedPrims = append(orderedPrims, bvh.primitives[primNum])
			}
			node.InitLeaf(firstPrimOffset, nPrimitives, bb)
			return node, orderedPrims
		}

		c0, c0ordered := bvh.bvhRecursiveBuild(buildData[start:mid], totalNodes, orderedPrims)
		c1, c1ordered := bvh.bvhRecursiveBuild(buildData[mid:], totalNodes, c0ordered)
		orderedPrims = c1ordered
		node.InitInterior(dim, c0, c1)
		return node, orderedPrims
	}

	// Partition primitives base on the Serfice Area Heuristic split method
	if nPrimitives <= 4 {
		// Instead of sort one can use C++'s std::nth_element-like function to partinion
		// the slice in two parts in O(n).
		sort.Slice(buildData, func(i, j int) bool {
			return buildData[i].centroid.ByAxis(dim) < buildData[j].centroid.ByAxis(dim)
		})
	} else {
		const nBuckets = 12
		var buckets [nBuckets]bucketInfo
		for i := start; i < end; i++ {
			b := int(nBuckets *
				((buildData[i].centroid.ByAxis(dim) - centroidBound.Min.ByAxis(dim)) /
					(centroidBound.Max.ByAxis(dim) - centroidBound.Min.ByAxis(dim))))

			if b == nBuckets {
				b = nBuckets - 1
			}

			buckets[b].count++
			buckets[b].bounds = bbox.Union(buckets[b].bounds, buildData[i].bounds)
		}

		const traversalCost = 0.85 // relative to ray-bbox intersection

		var cost [nBuckets - 1]float64

		for i := 0; i < nBuckets-1; i++ {
			var b0, b1 = bbox.Null(), bbox.Null()
			var count0, count1 int
			for j := 0; j < i; j++ {
				b0 = bbox.Union(b0, buckets[j].bounds)
				count0 += buckets[j].count
			}
			for j := i + 1; j < nBuckets; j++ {
				b1 = bbox.Union(b1, buckets[j].bounds)
				count1 += buckets[j].count
			}

			cost[i] = traversalCost + (float64(count0)*b0.SurfaceArea()+
				float64(count1)*b1.SurfaceArea())/
				bb.SurfaceArea()
		}

		var minCost = cost[0]
		var minCostSplit int

		for i := 1; i < nBuckets-1; i++ {
			if minCost > cost[i] {
				continue
			}
			minCost = cost[i]
			minCostSplit = i
		}

		// if nPrimitives > bvh.maxPrimsInNode check can be moved in the leaf creation in
		// the first if nPrimitives == 1
		if nPrimitives > bvh.maxPrimsInNode || minCost < float64(nPrimitives) {
			mid = partitionPrims(buildData,
				compareToBucket(minCostSplit, nBuckets, dim, centroidBound))
		} else {
			firstPrimOffset := len(orderedPrims)
			for i := start; i < end; i++ {
				primNum := buildData[i].primitiveNumber
				orderedPrims = append(orderedPrims, bvh.primitives[primNum])
			}
			node.InitLeaf(firstPrimOffset, nPrimitives, bb)
			return node, orderedPrims
		}
	}

	c0, c0ordered := bvh.bvhRecursiveBuild(buildData[start:mid], totalNodes, orderedPrims)
	c1, c1ordered := bvh.bvhRecursiveBuild(buildData[mid:], totalNodes, c0ordered)
	orderedPrims = c1ordered
	node.InitInterior(dim, c0, c1)

	return node, orderedPrims
}

// Intersect implements the Primitive interface
func (bvh *BVH) Intersect(ray geometry.Ray, in *primitive.Intersection) bool {
	if bvh.nodes == nil {
		return false
	}
	var hit bool

	// Used by the super-duper optimized ray-box intersection test
	// origin := ray.At(ray.Mint)
	invDir := geometry.NewVector(1/ray.Direction.X, 1/ray.Direction.Y, 1/ray.Direction.Z)
	dirIsNeg := [3]bool{
		invDir.X < 0,
		invDir.Y < 0,
		invDir.Z < 0,
	}

	var todoOffset, nodeNum uint32
	var todo [256]uint32
	for {
		node := &bvh.nodes[nodeNum]
		if node.bounds.IntersectPOptimized(&ray, &invDir, dirIsNeg) {
			if node.nPrimitives > 0 {
				// intersect with all primitives
				for i := uint32(0); i < uint32(node.nPrimitives); i++ {
					if bvh.primitives[node.offset+i].Intersect(ray, in) {
						hit = true
						ray.Maxt = in.DfGeometry.Distance
					}
				}
				if todoOffset == 0 {
					break
				}
				todoOffset--
				nodeNum = todo[todoOffset]
			} else {
				// put far bvh node on todo stack
				if dirIsNeg[node.axis] {
					todo[todoOffset] = nodeNum + 1
					todoOffset++
					nodeNum = node.offset
				} else {
					todo[todoOffset] = node.offset
					todoOffset++
					nodeNum = nodeNum + 1
				}
			}
		} else {
			if todoOffset == 0 {
				break
			}
			todoOffset--
			nodeNum = todo[todoOffset]
		}
	}

	return hit
}

// IntersectP implements the Primitive interface
func (bvh *BVH) IntersectP(ray geometry.Ray) bool {
	if bvh.nodes == nil {
		return false
	}

	// Used by the super-duper optimized ray-box intersection test
	// origin := ray.At(ray.Mint)
	invDir := geometry.NewVector(1/ray.Direction.X, 1/ray.Direction.Y, 1/ray.Direction.Z)
	dirIsNeg := [3]bool{
		invDir.X < 0,
		invDir.Y < 0,
		invDir.Z < 0,
	}

	var todoOffset, nodeNum uint32
	var todo [256]uint32
	for {
		node := &bvh.nodes[nodeNum]
		if node.bounds.IntersectPOptimized(&ray, &invDir, dirIsNeg) {
			if node.nPrimitives > 0 {
				// intersect with all primitives
				for i := uint32(0); i < uint32(node.nPrimitives); i++ {
					if bvh.primitives[node.offset+i].IntersectP(ray) {
						return true
					}
				}
				if todoOffset == 0 {
					break
				}
				todoOffset--
				nodeNum = todo[todoOffset]
			} else {
				// put far bvh node on todo stack
				if dirIsNeg[node.axis] {
					todo[todoOffset] = nodeNum + 1
					todoOffset++
					nodeNum = node.offset
				} else {
					todo[todoOffset] = node.offset
					todoOffset++
					nodeNum = nodeNum + 1
				}
			}
		} else {
			if todoOffset == 0 {
				break
			}
			todoOffset--
			nodeNum = todo[todoOffset]
		}
	}

	return false
}

func newBVHPrimitiveInfo(pn int, b *bbox.BBox) bvhPrimitiveInfo {
	return bvhPrimitiveInfo{
		primitiveNumber: pn,
		bounds:          b,
		centroid:        b.Min.MultiplyScalar(0.5).Plus(b.Max.MultiplyScalar(0.5)),
	}
}

type bucketInfo struct {
	count  int
	bounds *bbox.BBox
}

type bvhSplitFunction func(bvhPrimitiveInfo) bool

func compareToBucket(splitBucket, nBuckets, dim int, centroidBounds *bbox.BBox) bvhSplitFunction {
	return func(p bvhPrimitiveInfo) bool {
		b := int(float64(nBuckets) * (p.centroid.ByAxis(dim) - centroidBounds.Min.ByAxis(dim)) /
			(centroidBounds.Max.ByAxis(dim) - centroidBounds.Min.ByAxis(dim)))
		if b == nBuckets {
			b = nBuckets - 1
		}
		return b <= splitBucket
	}
}

// Basically implement the C++'s std::partition function
func partitionPrims(p []bvhPrimitiveInfo, comp bvhSplitFunction) int {
	placeIn := len(p) - 1
	var splitIndex int

	for i := 0; i < len(p) && i < placeIn; i++ {
		if comp(p[i]) {
			splitIndex = i + 1
			continue
		}
		for j := placeIn; j > i; j-- {
			placeIn = j
			if comp(p[j]) {
				p[i], p[placeIn] = p[placeIn], p[i]
				splitIndex = i + 1
				break
			}
		}
	}

	return splitIndex
}
