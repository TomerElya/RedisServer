package server

import (
	"fmt"
	"strings"
)

const (
	str  = '+'
	num  = ':'
	blk  = '$'
	arr  = '*'
	err1 = '-'
)

type Param struct {
	messageType   byte
	value         string
	chainedParams []Param
}

func (p *Param) ToString() string {
	res := strings.Builder{}
	switch p.messageType {
	case arr:
		res.WriteString(fmt.Sprintf("*%d\r\n", len(p.chainedParams)))
		for _, val := range p.chainedParams {
			res.WriteString(val.ToString())
		}
	case str:
		res.WriteString(fmt.Sprintf("+%s\r\n", p.value))
	case blk:
		res.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(p.value), p.value))
	case num:
		res.WriteString(fmt.Sprintf(":%s\r\n", p.value))
	case err1:
		res.WriteString(fmt.Sprintf("-%s\r\n", p.value))
	}
	return res.String()
}

func (p *Param) ToBytes() []byte {
	return []byte(p.ToString())
}
