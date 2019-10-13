package estimate

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/plantimals/genetic-physical-mapper/itree"
)

type Client struct {
	inFile   string
	output   string
	basesPer int64
}

func NewClient(inFile string, output string, basesPer int64) *Client {
	return &Client{
		inFile:   inFile,
		output:   output,
		basesPer: basesPer,
	}
}

func (c *Client) EstimateIntervals() error {
	inp, err := os.Open(c.inFile)
	if err != nil {
		return err
	}
	defer inp.Close()
	outp, err := os.OpenFile(c.output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer outp.Close()

	scanner := bufio.NewScanner(inp)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "\t")
		end, err := strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			return err
		}
		start, err := strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			return err
		}
		fmt.Fprintf(outp, "%s\t%v\n", line, float64(end-start)/float64(c.basesPer))
	}
	return nil
}

func (c *Client) InterpolateIntervals(gMap string) error {
	iTree, err := itree.New(gMap)
	if err != nil {
		return err
	}
	fmt.Printf("iTree forest: %v\tsource: %v\n", iTree.ForestSize(), iTree.SourceSize())
	inp, err := os.Open(c.inFile)
	if err != nil {
		return err
	}
	defer inp.Close()
	outp, err := os.OpenFile(c.output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer outp.Close()

	scanner := bufio.NewScanner(inp)
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line, "\t")
		end, err := strconv.ParseInt(data[6], 10, 64)
		if err != nil {
			return err
		}
		cmEnd, err := iTree.Interpolate(data[4], end)
		if err != nil {
			fmt.Println(err)
			continue
		}
		start, err := strconv.ParseInt(data[5], 10, 64)
		if err != nil {
			return err
		}
		cmStart, err := iTree.Interpolate(data[4], start)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Fprintf(outp, "%s\t%v\n", line, cmEnd-cmStart)
	}
	return nil
}
