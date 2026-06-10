package dtos

import (
	"fmt"
	"strconv"
	"strings"
)

type SearchQuery struct {
	Q       string
	Lat     *float64
	Long    *float64
	RangeKm float64
}

func NewSearchQuery(q, lat, long, rangeKm string) SearchQuery {
	sq := SearchQuery{
		Q:       q,
		Lat:     parseFloat(lat),
		Long:    parseFloat(long),
		RangeKm: 10.0,
	}
	if r := parseFloat(rangeKm); r != nil {
		sq.RangeKm = *r
	}
	return sq
}

func (s *SearchQuery) Validate() error {
	if strings.TrimSpace(s.Q) == "" {
		return fmt.Errorf("q is required")
	}
	return nil
}

func parseFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &v
}
