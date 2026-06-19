package dtos

import (
	"fmt"
	"strconv"
	"strings"
)

type SearchQuery struct {
	Q           string
	Lat         *float64
	Long        *float64
	RangeKm     float64
	CategoryIDs []uint
	MinPrice    *float64
	MaxPrice    *float64
}

func NewSearchQuery(q, lat, long, rangeKm, categoryIDs, minPrice, maxPrice string) SearchQuery {
	sq := SearchQuery{
		Q:           q,
		Lat:         parseFloat(lat),
		Long:        parseFloat(long),
		RangeKm:     10.0,
		CategoryIDs: parseUintList(categoryIDs),
		MinPrice:    parseFloat(minPrice),
		MaxPrice:    parseFloat(maxPrice),
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
	if s.MinPrice != nil && *s.MinPrice < 0 {
		return fmt.Errorf("min_price must be non-negative")
	}
	if s.MaxPrice != nil && s.MinPrice != nil && *s.MaxPrice < *s.MinPrice {
		return fmt.Errorf("max_price must be >= min_price")
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

func parseUintList(s string) []uint {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseUint(p, 10, 64)
		if err == nil {
			out = append(out, uint(v))
		}
	}
	return out
}
