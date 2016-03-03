package ts3

import "strings"

type message struct {
	Contents string
	Error    error
}

type Response struct {
	Values map[string]string
}

func newResponse(input string) (*Response, error) {
	resp := &Response{}
	resp.Values = make(map[string]string)

	for _, chunk := range strings.Split(input, " ") {
		if err := parseChunk(resp, chunk); err != nil {
		}
	}

	return resp, nil
}

func parseChunk(resp *Response, chunk string) error {
	idx := strings.Index(chunk, "=")
	if idx == -1 {
		return
	}

	left, right := chunk[:idx], chunk[idx+1:]
	resp.Values[left] = right

	return nil
}
