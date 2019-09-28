package estimate

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
