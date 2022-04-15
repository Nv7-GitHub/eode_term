package main

import (
	"math/big"
	"strconv"
	"time"
)

// github.com/Nv7-Github/Nv7haven/eod/api/data/data.go

type Method int

const (
	MethodGuild Method = iota
	MethodElem
	MethodCombo
	MethodElemInfo
	MethodInv
)

type Message struct {
	Method Method         `json:"method"`
	Params map[string]any `json:"params"`
}

type Resp struct {
	Error *string        `json:"error,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

// github.com/Nv7-Github/Nv7haven/eod/types/types.go
type Element struct {
	ID         int
	Name       string
	Image      string
	Color      int
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  *TimeStamp
	Parents    []int
	Complexity int
	Difficulty int
	UsedIn     int
	TreeSize   int
	Air        *big.Int
	Earth      *big.Int
	Fire       *big.Int
	Water      *big.Int

	Commenter string
	Colorer   string
	Imager    string
}

type TimeStamp struct {
	time.Time
}

func (t *TimeStamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return t.Time.UnmarshalJSON(data)
	}
	t.Time = time.Unix(i, 0)
	return nil
}
