package services

import (
	"context"
	"log"

	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"github.com/mmcloughlin/geohash"
)

type DealService struct{ base_service }

func NewDealService(b base_service) *DealService {
	return &DealService{base_service: b}
}

func (s *DealService) NearbyDeals(lat, long float64) response.Res {
	ctx := context.Background()
	gh5 := geohash.EncodeWithPrecision(lat, long, 5)
	gc := s.cache.GetDealGeoCache()

	dealIDs, err := gc.GetDealIDs(ctx, gh5)
	if err != nil {
		log.Printf("deals: redis get %s: %v", gh5, err)
	}

	if dealIDs == nil {
		entries, dbErr := s.RepoDeal().GetByGeoHash5(ctx, gh5)
		if dbErr != nil {
			return s.ResponseInternalServer("failed to fetch deals")
		}
		ids := make([]uint, 0, len(entries))
		for _, e := range entries {
			ids = append(ids, e.DealID)
		}
		if err := gc.PopulateDealIDs(ctx, gh5, ids); err != nil {
			log.Printf("deals: repopulate %s: %v", gh5, err)
		}
		return s.ResponseOK("nearby deals", entries)
	}

	if len(dealIDs) == 0 {
		return s.ResponseOK("nearby deals", []struct{}{})
	}
	entries, err := s.RepoDeal().GetByIDs(ctx, dealIDs)
	if err != nil {
		return s.ResponseInternalServer("failed to fetch deal details")
	}
	return s.ResponseOK("nearby deals", entries)
}
