package handlers

import (
	"net/http"
	"strconv"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type ReviewHandler struct {
	BaseHandler
}

func NewReviewHandler(services *services.Services) *ReviewHandler {
	return &ReviewHandler{BaseHandler: *NewBaserHandler(services)}
}

// reviewBody is the expected JSON body for add/update/delete review requests.
type reviewBody struct {
	TargetID uint   `json:"targetId"`
	Stars    uint8  `json:"stars"`
	Comment  string `json:"comment"`
}

// idFromQuery reads ?id= as a uint (0 if missing or invalid).
func idFromQuery(r *http.Request) uint {
	raw := r.URL.Query().Get("id")
	if raw == "" {
		return 0
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return 0
	}
	return uint(v)
}

// ── Business reviews ─────────────────────────────────────────────────────────

func (rev *ReviewHandler) ReviewGetByBusiness(w http.ResponseWriter, r *http.Request) {
	id := idFromQuery(r)
	rev.Response(w, rev.GetReviewService().GetReviews(models.ReviewTargetBusiness, id))
}

func (rev *ReviewHandler) ReviewBusiness(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetBusiness, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewBusinessUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetBusiness, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewBusinessDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetBusiness, body.TargetID))
}

// ── Product reviews ───────────────────────────────────────────────────────────

func (rev *ReviewHandler) ReviewGetByProduct(w http.ResponseWriter, r *http.Request) {
	id := idFromQuery(r)
	rev.Response(w, rev.GetReviewService().GetReviews(models.ReviewTargetProduct, id))
}

func (rev *ReviewHandler) ReviewProduct(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetProduct, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewProductUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetProduct, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewProductDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetProduct, body.TargetID))
}

// ── Offer reviews ─────────────────────────────────────────────────────────────

func (rev *ReviewHandler) ReviewGetByOffer(w http.ResponseWriter, r *http.Request) {
	offerID, err := rev.GetOfferId(r)
	if err != nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid offer ID")
		return
	}
	rev.Response(w, rev.GetReviewService().GetReviews(models.ReviewTargetOffer, offerID))
}

func (rev *ReviewHandler) ReviewOffer(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetOffer, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewOfferUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetOffer, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewOfferDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetOffer, body.TargetID))
}

// ── Spotlight reviews ─────────────────────────────────────────────────────────

func (rev *ReviewHandler) ReviewGetBySpotlight(w http.ResponseWriter, r *http.Request) {
	spotlightID, err := rev.GetSpotlightId(r)
	if err != nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid spotlight ID")
		return
	}
	rev.Response(w, rev.GetReviewService().GetReviews(models.ReviewTargetSpotlight, spotlightID))
}

func (rev *ReviewHandler) ReviewSpotlight(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetSpotlight, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewSpotlightUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetSpotlight, body.TargetID, body.Stars, body.Comment))
}

func (rev *ReviewHandler) ReviewSpotlightDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	body, err := ParseBody[reviewBody](r)
	if err != nil || body == nil {
		rev.ResponseBadRequestWithMessage(w, "Invalid request body")
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetSpotlight, body.TargetID))
}
