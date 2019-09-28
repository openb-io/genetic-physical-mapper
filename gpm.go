package main

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/plantimals/genetic-physical-mapper/estimate"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app    = kingpin.New("gpm", "move between genetic coordinates and physical coordinates with speed and confidence")
	est    = app.Command("intervals", "estimate centimorgan span of intervals")
	input  = est.Flag("input", "path to input data, (g)zip or ascii").Required().Short('i').ExistingFile()
	output = est.Flag("output", "path to output data").Required().Short('o').String()
	bases  = est.Flag("bases", "bases per centimorgan").Required().Short('b').Int64()
)

func main() {
	app.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0.0").Author("Rob Long")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case est.FullCommand():
		RunEstimate()
	}
}

func RunEstimate() {

	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Prefix = "estimating genetic coordinates"
	s.Start()

	client := estimate.NewClient(*input, *output, *bases)
	err := client.EstimateIntervals()
	if err != nil {
		panic(err)
	}

	s.Stop()
}
