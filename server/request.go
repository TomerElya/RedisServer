package server

type Request struct {
	action string
	params []Param
	client *Client
}
