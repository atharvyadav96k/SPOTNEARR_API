package handlers

import (
	"net/http"
	"strconv"

	"github.com/atharvyadav96k/SPOTNEARR_API/dtos"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type ReviewHandler struct {
	BaseHandler
}

func NewReviewHandler(services *services.Services) *ReviewHandler {
	return &ReviewHandler{BaseHandler: *NewBaserHandler(services)}
}

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

func (rev *ReviewHandler) reviewWrite(w http.ResponseWriter, r *http.Request, svcFn func(uint, uint, uint8, string)) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	svcFn(userID, dto.TargetID, dto.Stars, dto.Comment)
}

func (rev *ReviewHandler) reviewDelete(w http.ResponseWriter, r *http.Request, svcFn func(uint, uint)) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	svcFn(userID, dto.TargetID)
}

func (rev *ReviewHandler) ReviewGetByBusiness(w http.ResponseWriter, r *http.Request) {
	rev.Response(w, rev.GetReviewService().GetReviews(models.ReviewTargetBusiness, idFromQuery(r)))
}

func (rev *ReviewHandler) ReviewBusiness(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetBusiness, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewBusinessUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetBusiness, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewBusinessDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetBusiness, dto.TargetID))
}

func (rev *ReviewHandler) ReviewGetByProduct(w http.ResponseWriter, r *http.Request) {
	rev.Response(w, rev.GetReviewService().GetReviews(models.ReviewTargetProduct, idFromQuery(r)))
}

func (rev *ReviewHandler) ReviewProduct(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetProduct, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewProductUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetProduct, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewProductDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetProduct, dto.TargetID))
}

func (rev *ReviewHandler) ReviewGetByOffer(w http.ResponseWriter, r *http.Request) {
	offerID, err := rev.GetOfferId(r)
	if err != nil {
		rev.ResponseBadRequestWithMessage(w, "invalid offer ID")
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
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetOffer, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewOfferUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetOffer, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewOfferDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetOffer, dto.TargetID))
}

func (rev *ReviewHandler) ReviewGetBySpotlight(w http.ResponseWriter, r *http.Request) {
	spotlightID, err := rev.GetSpotlightId(r)
	if err != nil {
		rev.ResponseBadRequestWithMessage(w, "invalid spotlight ID")
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
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().AddReview(userID, models.ReviewTargetSpotlight, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewSpotlightUpdate(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().UpdateReview(userID, models.ReviewTargetSpotlight, dto.TargetID, dto.Stars, dto.Comment))
}

func (rev *ReviewHandler) ReviewSpotlightDelete(w http.ResponseWriter, r *http.Request) {
	userID := rev.ClaimGetUserId(r)
	if userID == 0 {
		rev.ResponseBadRequest(w)
		return
	}
	var dto dtos.ReviewRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		rev.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	rev.Response(w, rev.GetReviewService().DeleteReview(userID, models.ReviewTargetSpotlight, dto.TargetID))
}
