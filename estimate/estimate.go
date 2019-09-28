package estimate

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
