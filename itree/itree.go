package itree

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-collections/go-datastructures/augmentedtree"
)

type iTree struct {
	forest map[string]augmentedtree.Tree
	source map[string]map[int64]float64
}

func New(gMap string) (*iTree, error) {
	g, err := os.Open(gMap)
	if err != nil {
		return nil, err
	}
	defer g.Close()

	forest := make(map[string]augmentedtree.Tree)
	source := make(map[string]map[int64]float64)

	answer := &iTree{
		forest: forest,
		source: source,
	}

	scnnr := bufio.NewScanner(g)
	scnnr.Scan()
	prevLocus, err := ParseLocus(scnnr.Text())
	if err != nil {
		return nil, err
	}
	answer.forest[prevLocus.chr] = augmentedtree.New(1)
	answer.source[prevLocus.chr] = make(map[int64]float64)
	answer.source[prevLocus.chr][prevLocus.pos] = prevLocus.cM

	for scnnr.Scan() {
		locus, err := ParseLocus(scnnr.Text())
		if err != nil {
			return nil, err
		}
		if CheckChrs(prevLocus, locus) {
			interval, err := MakeInterval(prevLocus, locus)
			if err != nil {
				return nil, err
			}
			answer.source[locus.chr][locus.pos] = locus.cM
			s1 := answer.forest[prevLocus.chr].Len()
			answer.forest[prevLocus.chr].Add(interval)
			s2 := answer.forest[prevLocus.chr].Len()
			if s2-s1 == 0 {
				fmt.Printf("s1: %v\ts2: %v\n", s1, s2)
				return nil, fmt.Errorf("failed to add interval %s:%v-%v", prevLocus.chr, prevLocus.pos, locus.pos)
			}
			prevLocus = locus
		} else {
			answer.forest[locus.chr] = augmentedtree.New(1)
			answer.source[locus.chr] = make(map[int64]float64)
			answer.source[locus.chr][locus.pos] = locus.cM
			prevLocus = locus
		}

	}
	return answer, nil
}

func (i *iTree) ForestSize() int {
	total := 0
	for _, tree := range i.forest {
		total += int(tree.Len())
	}
	return total
}

func (i *iTree) SourceSize() int {
	total := 0
	for _, m := range i.source {
		total += len(m)
	}
	return total
}

func CheckChrs(l1 *Locus, l2 *Locus) bool {
	return l1.chr == l2.chr
}

func MakeInterval(l1 *Locus, l2 *Locus) (*GenomicInterval, error) {
	i, err := NewInterval(l1, l2)
	if err != nil {
		return nil, err
	}
	return i, nil
}

type Locus struct {
	chr string
	cM  float64
	pos int64
}

func ParseLocus(line string) (*Locus, error) {
	fields := strings.Split(line, "\t")
	if len(fields) != 4 {
		return nil, fmt.Errorf("found improper number of fields in the genetic map file: %v fields", len(fields))
	}

	cM, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}
	pos, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}
	return &Locus{
		chr: fields[0],
		cM:  cM,
		pos: pos,
	}, nil
}

type GenomicInterval struct {
	id     uint64
	chr    string
	start  int64
	end    int64
	gStart float64
	gEnd   float64
}

func NewInterval(l1 *Locus, l2 *Locus) (*GenomicInterval, error) {

	if l1.chr != l2.chr {
		return nil, fmt.Errorf("chromosome mismatch: %s vs %s", l1.chr, l2.chr)
	}
	if l1.pos > l2.pos {
		return nil, fmt.Errorf("locus 1 has position higher than locus 2: %s %s")
	}
	return &GenomicInterval{
		id:     uint64(l1.pos),
		chr:    l1.chr,
		start:  l1.pos,
		end:    l2.pos,
		gStart: l1.cM,
		gEnd:   l2.cM,
	}, nil
}

func (g GenomicInterval) LowAtDimension(d uint64) int64 {
	return g.start
}

func (g GenomicInterval) HighAtDimension(d uint64) int64 {
	return g.end
}

func (g GenomicInterval) OverlapsAtDimension(i augmentedtree.Interval, d uint64) bool {
	return g.start < i.HighAtDimension(0) && i.LowAtDimension(0) < g.end
}

func (g GenomicInterval) ID() uint64 {
	return g.id
}

func (g GenomicInterval) GeneticStart() float64 {
	return g.gStart
}

func (g GenomicInterval) GeneticEnd() float64 {
	return g.gEnd
}

func (i *iTree) Interpolate(chr string, pos int64) (float64, error) {
	result := i.forest[chr].Query(&GenomicInterval{
		chr:   chr,
		start: pos,
		end:   pos + 1,
	})
	if len(result) < 1 {
		return 0.0, fmt.Errorf("no intersecting interval found for %s:%v", chr, pos)
	}
	return LinearInterpolation(
		float64(result[0].LowAtDimension(0)),
		float64(result[0].HighAtDimension(0)),
		float64(pos),
		i.source[chr][result[0].LowAtDimension(0)],
		i.source[chr][result[0].HighAtDimension(0)],
	), nil
}

func LinearInterpolation(x1 float64, x2 float64, xa float64, y1 float64, y2 float64) float64 {
	fmt.Printf("x1: %v\tx2: %v\txa: %v\ty1: %v\ty2: %v\n", x1, x2, xa, y1, y2)
	return (((xa - x1) / (x2 - x1)) * (y2 - y1)) + y1
}
