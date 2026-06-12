package mq

// Topic is a routing key on the spotnearr.events exchange.
type Topic string

const (
	TopicBusinessFollow   Topic = "user.business.follow"
	TopicBusinessUnfollow Topic = "user.business.unfollow"
)

// BusinessFollowPayload is published for both follow and unfollow events.
type BusinessFollowPayload struct {
	BusinessID uint `json:"business_id"`
	UserID     uint `json:"user_id"`
}
