package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	userdb "github.com/Developer-Aadesh/spotnearr-database/user"
	userpostgres "github.com/Developer-Aadesh/spotnearr-database/user/postgres"
	"github.com/atharvyadav96k/spotnearr/pkg/mq"
	"gorm.io/gorm"
)

type DealConsumer struct {
	sub      *mq.Subscriber
	dealRepo *userpostgres.DealRepository
	cache    *cache.Cache
}

func NewDealConsumer(sub *mq.Subscriber, db *gorm.DB, c *cache.Cache) *DealConsumer {
	return &DealConsumer{
		sub:      sub,
		dealRepo: userpostgres.NewDealRepository(db),
		cache:    c,
	}
}

func (d *DealConsumer) Run(ctx context.Context) error {
	if err := d.sub.Subscribe(ctx, "user.deal.sync", mq.TopicDealSync, d.handleSync); err != nil {
		return err
	}
	return d.sub.Subscribe(ctx, "user.deal.delete", mq.TopicDealDelete, d.handleDelete)
}

func (d *DealConsumer) handleSync(ctx context.Context, body []byte) error {
	var p mq.DealSyncPayload
	if err := json.Unmarshal(body, &p); err != nil {
		return err
	}
	gc := d.cache.GetDealGeoCache()
	for _, loc := range p.Locations {
		entry := userdb.DealEntry{
			DealID:       p.DealID,
			GeoHash5:     loc.GeoHash5,
			BusinessID:   p.BusinessID,
			Name:         p.Title,
			Price:        p.MinOrderValue,
			DealPrice:    p.DiscountValue,
			DiscountType: p.DiscountType,
			Active:       p.Active,
			ExpiresAt:    p.ExpiresAt,
		}
		if err := d.dealRepo.Upsert(ctx, entry); err != nil {
			log.Printf("deal_consumer: upsert deal %d/%s: %v", p.DealID, loc.GeoHash5, err)
			return err
		}
		if err := gc.AddDeal(ctx, loc.GeoHash5, p.DealID); err != nil {
			log.Printf("deal_consumer: cache add %d/%s: %v", p.DealID, loc.GeoHash5, err)
		}
	}
	return nil
}

func (d *DealConsumer) handleDelete(ctx context.Context, body []byte) error {
	var p mq.DealDeletePayload
	if err := json.Unmarshal(body, &p); err != nil {
		return err
	}
	hashes, err := d.dealRepo.GetGeoHash5ForDeal(ctx, p.DealID)
	if err != nil {
		log.Printf("deal_consumer: get hashes for deal %d: %v", p.DealID, err)
	}
	if err := d.dealRepo.DeleteByDealID(ctx, p.DealID); err != nil {
		return err
	}
	gc := d.cache.GetDealGeoCache()
	for _, gh5 := range hashes {
		if err := gc.RemoveDeal(ctx, gh5, p.DealID); err != nil {
			log.Printf("deal_consumer: cache remove %d/%s: %v", p.DealID, gh5, err)
		}
	}
	return nil
}
