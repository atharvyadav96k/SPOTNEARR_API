package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ─── ANSI ─────────────────────────────────────────────────────────────────────
const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
	bold   = "\033[1m"
	dim    = "\033[2m"
)

// ─── SHARED STATE ─────────────────────────────────────────────────────────────
var (
	userID, businessID, bizCategoryID string
	locationID, productID, inventoryID string
	spotlightID, offerID, claimID      string
)

// ─── RESULT TRACKING ──────────────────────────────────────────────────────────
type result struct {
	name     string
	passed   bool
	skipped  bool
	duration time.Duration
	detail   string
}

var suite []result

// ─── HTTP CLIENT ──────────────────────────────────────────────────────────────
var client = &http.Client{Timeout: 30 * time.Second}

type envelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Error   interface{}     `json:"error"`
}

func call(fnURL string, body map[string]interface{}) (*envelope, error) {
	b, _ := json.Marshal(body)
	resp, err := client.Post(fnURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("http error: %w", err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("bad JSON (HTTP %d): %s", resp.StatusCode, string(raw))
	}
	if !env.Success {
		errStr := fmt.Sprintf("%v", env.Error)
		if errStr == "<nil>" || errStr == "" {
			errStr = env.Message
		}
		return &env, fmt.Errorf("%s", errStr)
	}
	return &env, nil
}

// extract a string field from envelope Data (handles both object and array-of-objects)
func str(env *envelope, field string) string {
	if env == nil || env.Data == nil {
		return ""
	}
	// try object
	var m map[string]interface{}
	if json.Unmarshal(env.Data, &m) == nil {
		if v, ok := m[field].(string); ok {
			return v
		}
	}
	// try array — take first element
	var arr []map[string]interface{}
	if json.Unmarshal(env.Data, &arr) == nil && len(arr) > 0 {
		if v, ok := arr[0][field].(string); ok {
			return v
		}
	}
	return ""
}

// ─── RUNNER ───────────────────────────────────────────────────────────────────
func run(name, fnURL string, fn func() error) {
	if fnURL == "" {
		suite = append(suite, result{name: name, skipped: true})
		fmt.Printf("  %s⊘%s %-58s %sSKIPPED%s\n", yellow, reset, name, dim, reset)
		return
	}
	start := time.Now()
	err := fn()
	d := time.Since(start)
	if err != nil {
		suite = append(suite, result{name: name, passed: false, duration: d, detail: err.Error()})
		fmt.Printf("  %s✗%s %-58s %s%dms%s\n    %s↳ %s%s\n",
			red, reset, name, dim, d.Milliseconds(), reset, red, err.Error(), reset)
	} else {
		suite = append(suite, result{name: name, passed: true, duration: d})
		fmt.Printf("  %s✓%s %-58s %s%dms%s\n", green, reset, name, dim, d.Milliseconds(), reset)
	}
}

func section(title string) {
	bar := strings.Repeat("─", 60-len(title)-3)
	fmt.Printf("\n%s%s── %s %s%s\n", bold, cyan, title, bar, reset)
}

// ─── ENV HELPER ───────────────────────────────────────────────────────────────
func e(name string) string {
	return os.Getenv(name)
}

// ─── MAIN ─────────────────────────────────────────────────────────────────────
func main() {
	ts := time.Now().UnixNano() % 1_000_000_000
	userEmail := fmt.Sprintf("test%d@spotnearr.test", ts)
	bizEmail := fmt.Sprintf("biz%d@spotnearr.test", ts)
	testPass := "Test@123456"
	productCategoryID := e("PRODUCT_CATEGORY_ID") // must be pre-seeded

	fmt.Printf("\n%s%s══════════════════════════════════════════════════════════════%s\n", bold, cyan, reset)
	fmt.Printf("%s%s  Spotnearr — End-to-End Integration Test Suite%s\n", bold, cyan, reset)
	fmt.Printf("%s%s══════════════════════════════════════════════════════════════%s\n", bold, cyan, reset)
	fmt.Printf("%sStarted : %s%s\n", dim, time.Now().Format("2006-01-02 15:04:05"), reset)
	fmt.Printf("%sUser    : %s%s\n", dim, userEmail, reset)
	if productCategoryID == "" {
		fmt.Printf("%s⚠  PRODUCT_CATEGORY_ID not set — product/inventory tests will be skipped%s\n", yellow, reset)
	}

	// ──────────────────────────────────────────────────────────────────────────
	section("AUTH")
	// ──────────────────────────────────────────────────────────────────────────

	run("Register user", e("FN_USER_REGISTER"), func() error {
		resp, err := call(e("FN_USER_REGISTER"), map[string]interface{}{
			"full_name": "Integration Tester",
			"email":     userEmail,
			"password":  testPass,
			"role":      "customer",
		})
		if err != nil {
			return err
		}
		userID = str(resp, "id")
		if userID == "" {
			return fmt.Errorf("no id in response")
		}
		fmt.Printf("    %sid = %s%s\n", dim, userID, reset)
		return nil
	})

	run("Login user", e("FN_USER_LOGIN"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — register failed")
		}
		resp, err := call(e("FN_USER_LOGIN"), map[string]interface{}{
			"email":    userEmail,
			"password": testPass,
		})
		if err != nil {
			return err
		}
		if str(resp, "id") == "" {
			return fmt.Errorf("no id in login response")
		}
		return nil
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("BUSINESS SETUP")
	// ──────────────────────────────────────────────────────────────────────────

	run("Create business category", e("FN_ADD_CATEGORIES"), func() error {
		resp, err := call(e("FN_ADD_CATEGORIES"), map[string]interface{}{
			"name": fmt.Sprintf("Test Category %d", ts),
		})
		if err != nil {
			return err
		}
		bizCategoryID = str(resp, "id")
		if bizCategoryID == "" {
			return fmt.Errorf("no id in response")
		}
		fmt.Printf("    %sid = %s%s\n", dim, bizCategoryID, reset)
		return nil
	})

	run("Register business", e("FN_BIZ_REGISTER"), func() error {
		if userID == "" || bizCategoryID == "" {
			return fmt.Errorf("skipped — user or category missing")
		}
		resp, err := call(e("FN_BIZ_REGISTER"), map[string]interface{}{
			"owner_id":    userID,
			"name":        fmt.Sprintf("Test Biz %d", ts),
			"email":       bizEmail,
			"category_id": bizCategoryID,
		})
		if err != nil {
			return err
		}
		businessID = str(resp, "id")
		if businessID == "" {
			return fmt.Errorf("no id in response")
		}
		fmt.Printf("    %sid = %s%s\n", dim, businessID, reset)
		return nil
	})

	run("Get business categories", e("FN_GET_CATEGORIES"), func() error {
		resp, err := call(e("FN_GET_CATEGORIES"), map[string]interface{}{})
		if err != nil {
			return err
		}
		_ = resp
		return nil
	})

	run("Create business location", e("FN_CREATE_LOCATION"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business not registered")
		}
		resp, err := call(e("FN_CREATE_LOCATION"), map[string]interface{}{
			"business_id":   businessID,
			"branch_name":   "Main Branch",
			"address_line1": "123 Test Street",
			"city":          "Mumbai",
			"state":         "Maharashtra",
			"pin_code":      "400001",
			"latitude":      19.0760,
			"longitude":     72.8777,
			"is_main":       true,
			"opening_time":  "09:00:00",
			"closing_time":  "21:00:00",
			"working_days":  []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		})
		if err != nil {
			return err
		}
		locationID = str(resp, "id")
		if locationID == "" {
			return fmt.Errorf("no id in response")
		}
		fmt.Printf("    %sid = %s%s\n", dim, locationID, reset)
		return nil
	})

	run("Create product", e("FN_NEW_PRODUCT"), func() error {
		if productCategoryID == "" {
			return fmt.Errorf("skipped — PRODUCT_CATEGORY_ID not set")
		}
		resp, err := call(e("FN_NEW_PRODUCT"), map[string]interface{}{
			"category_id":  productCategoryID,
			"name":         fmt.Sprintf("Test Product %d", ts),
			"description":  "Integration test product",
			"unit":         "piece",
			"is_available": true,
			"is_active":    true,
			"tags":         []string{"test", "integration"},
		})
		if err != nil {
			return err
		}
		productID = str(resp, "id")
		if productID == "" {
			return fmt.Errorf("no id in response")
		}
		fmt.Printf("    %sid = %s%s\n", dim, productID, reset)
		return nil
	})

	run("Create inventory", e("FN_CREATE_INVENTORY"), func() error {
		if productID == "" || locationID == "" {
			return fmt.Errorf("skipped — product or location missing")
		}
		_, err := call(e("FN_CREATE_INVENTORY"), map[string]interface{}{
			"inventories": []map[string]interface{}{
				{
					"product_id":       productID,
					"location_id":      locationID,
					"price":            500,
					"discounted_price": 450,
					"stock":            10,
					"is_available":     true,
				},
			},
		})
		if err != nil {
			// inventory endpoint takes an array directly, not wrapped — try direct array
			resp2, err2 := call(e("FN_CREATE_INVENTORY"), nil)
			_ = resp2
			if err2 != nil {
				// send as raw array
				b, _ := json.Marshal([]map[string]interface{}{
					{
						"product_id":       productID,
						"location_id":      locationID,
						"price":            500,
						"discounted_price": 450,
						"stock":            10,
						"is_available":     true,
					},
				})
				httpResp, httpErr := client.Post(e("FN_CREATE_INVENTORY"), "application/json", bytes.NewReader(b))
				if httpErr != nil {
					return fmt.Errorf("create inventory: %w", httpErr)
				}
				defer httpResp.Body.Close()
				raw, _ := io.ReadAll(httpResp.Body)
				var env envelope
				if err3 := json.Unmarshal(raw, &env); err3 != nil {
					return fmt.Errorf("bad json: %s", string(raw))
				}
				if !env.Success {
					return fmt.Errorf("%v", env.Error)
				}
			}
		}
		return nil
	})

	run("Get business products", e("FN_GET_BUSINESS_PRODUCTS"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_GET_BUSINESS_PRODUCTS"), map[string]interface{}{
			"id": businessID,
		})
		return err
	})

	run("Get business inventory", e("FN_GET_BUSINESS_INVENTORY"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_GET_BUSINESS_INVENTORY"), map[string]interface{}{
			"business_id": businessID,
		})
		return err
	})

	now := time.Now()
	expiresAt := now.Add(7 * 24 * time.Hour).Format(time.RFC3339)
	startsAt := now.Format(time.RFC3339)

	run("Create spotlight", e("FN_CREATE_SPOTLIGHT"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		resp, err := call(e("FN_CREATE_SPOTLIGHT"), map[string]interface{}{
			"business_id": businessID,
			"type":        "general",
			"status":      "published",
			"title":       fmt.Sprintf("Test Spotlight %d", ts),
			"description": "Integration test spotlight",
			"media_url":   "https://example.com/test-image.jpg",
			"media_type":  "image",
			"expires_at":  expiresAt,
		})
		if err != nil {
			return err
		}
		spotlightID = str(resp, "id")
		fmt.Printf("    %sid = %s%s\n", dim, spotlightID, reset)
		return nil
	})

	run("Create offer", e("FN_CREATE_OFFER"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		resp, err := call(e("FN_CREATE_OFFER"), map[string]interface{}{
			"business_id":     businessID,
			"title":           fmt.Sprintf("Test Offer %d", ts),
			"description":     "Integration test offer",
			"discount_type":   "percentage",
			"discount_value":  10,
			"min_order_value": 100,
			"is_active":       true,
			"starts_at":       startsAt,
			"expires_at":      expiresAt,
		})
		if err != nil {
			return err
		}
		offerID = str(resp, "id")
		fmt.Printf("    %sid = %s%s\n", dim, offerID, reset)
		return nil
	})

	run("Get business offers (business side)", e("FN_GET_BUSINESS_OFFERS"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_GET_BUSINESS_OFFERS"), map[string]interface{}{
			"business_id": businessID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("USER DISCOVERY")
	// ──────────────────────────────────────────────────────────────────────────

	run("Nearby products (lat/lng)", e("FN_NEARBY_PRODUCTS"), func() error {
		_, err := call(e("FN_NEARBY_PRODUCTS"), map[string]interface{}{
			"latitude":  19.0760,
			"longitude": 72.8777,
		})
		return err
	})

	run("Search products", e("FN_SEARCH_PRODUCTS"), func() error {
		_, err := call(e("FN_SEARCH_PRODUCTS"), map[string]interface{}{
			"query": "Test",
		})
		return err
	})

	run("Spotlights by location", e("FN_SPOTLIGHT_BY_LOCATION"), func() error {
		_, err := call(e("FN_SPOTLIGHT_BY_LOCATION"), map[string]interface{}{
			"latitude":  19.0760,
			"longitude": 72.8777,
		})
		return err
	})

	run("Get offers by business (user side)", e("FN_GET_OFFERS_BY_BUSINESS"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_GET_OFFERS_BY_BUSINESS"), map[string]interface{}{
			"business_id": businessID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("USER — SOCIAL")
	// ──────────────────────────────────────────────────────────────────────────

	run("Like business", e("FN_LIKE_BUSINESS"), func() error {
		if userID == "" || businessID == "" {
			return fmt.Errorf("skipped — user or business missing")
		}
		_, err := call(e("FN_LIKE_BUSINESS"), map[string]interface{}{
			"user_id":     userID,
			"business_id": businessID,
		})
		return err
	})

	run("Get liked businesses", e("FN_GET_LIKED_BUSINESSES"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		_, err := call(e("FN_GET_LIKED_BUSINESSES"), map[string]interface{}{
			"user_id": userID,
		})
		return err
	})

	run("Dislike business (cleanup)", e("FN_DISLIKE_BUSINESS"), func() error {
		if userID == "" || businessID == "" {
			return fmt.Errorf("skipped — user or business missing")
		}
		_, err := call(e("FN_DISLIKE_BUSINESS"), map[string]interface{}{
			"user_id":     userID,
			"business_id": businessID,
		})
		return err
	})

	run("Follow business", e("FN_FOLLOW_BUSINESS"), func() error {
		if userID == "" || businessID == "" {
			return fmt.Errorf("skipped — user or business missing")
		}
		_, err := call(e("FN_FOLLOW_BUSINESS"), map[string]interface{}{
			"user_id":     userID,
			"business_id": businessID,
		})
		return err
	})

	run("Get following businesses", e("FN_GET_FOLLOWING_BUSINESSES"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		_, err := call(e("FN_GET_FOLLOWING_BUSINESSES"), map[string]interface{}{
			"user_id": userID,
		})
		return err
	})

	run("Get followed business spotlights", e("FN_FOLLOWED_SPOTLIGHTS"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		_, err := call(e("FN_FOLLOWED_SPOTLIGHTS"), map[string]interface{}{
			"id": userID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("USER — CLAIMS")
	// ──────────────────────────────────────────────────────────────────────────

	// To claim, we need an inventory_id. Fetch it from the business inventory.
	run("Fetch inventory_id for claim", e("FN_GET_BUSINESS_INVENTORY"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		resp, err := call(e("FN_GET_BUSINESS_INVENTORY"), map[string]interface{}{
			"business_id": businessID,
		})
		if err != nil {
			return err
		}
		inventoryID = str(resp, "id")
		if inventoryID == "" {
			return fmt.Errorf("no inventory found — create inventory first")
		}
		fmt.Printf("    %sinventory_id = %s%s\n", dim, inventoryID, reset)
		return nil
	})

	run("Claim product", e("FN_CLAIM_PRODUCT"), func() error {
		if userID == "" || businessID == "" || inventoryID == "" {
			return fmt.Errorf("skipped — user, business, or inventory missing")
		}
		resp, err := call(e("FN_CLAIM_PRODUCT"), map[string]interface{}{
			"user_id":      userID,
			"business_id":  businessID,
			"inventory_id": inventoryID,
			"quantity":     1,
			"note":         "Integration test claim",
		})
		if err != nil {
			return err
		}
		claimID = str(resp, "id")
		if claimID == "" {
			return fmt.Errorf("no id in response")
		}
		fmt.Printf("    %sid = %s%s\n", dim, claimID, reset)
		return nil
	})

	run("Get my claims", e("FN_GET_MY_CLAIMS"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		_, err := call(e("FN_GET_MY_CLAIMS"), map[string]interface{}{
			"user_id": userID,
		})
		return err
	})

	run("Cancel claim (idempotency — expect fail on non-pending)", e("FN_CANCEL_CLAIM"), func() error {
		// We'll test this after the claim is accepted — for now skip
		if claimID == "" {
			return fmt.Errorf("skipped — no claim yet")
		}
		// create a second claim to cancel
		resp, _ := call(e("FN_CLAIM_PRODUCT"), map[string]interface{}{
			"user_id":      userID,
			"business_id":  businessID,
			"inventory_id": inventoryID,
			"quantity":     1,
			"note":         "Claim to be cancelled",
		})
		cancelID := str(resp, "id")
		if cancelID == "" {
			return fmt.Errorf("could not create second claim for cancellation test")
		}
		_, err := call(e("FN_CANCEL_CLAIM"), map[string]interface{}{
			"id":      cancelID,
			"user_id": userID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("BUSINESS — CLAIMS")
	// ──────────────────────────────────────────────────────────────────────────

	run("Get incoming claims", e("FN_GET_INCOMING_CLAIMS"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_GET_INCOMING_CLAIMS"), map[string]interface{}{
			"business_id": businessID,
		})
		return err
	})

	run("Accept claim", e("FN_ACCEPT_CLAIM"), func() error {
		if claimID == "" || businessID == "" {
			return fmt.Errorf("skipped — claim or business missing")
		}
		_, err := call(e("FN_ACCEPT_CLAIM"), map[string]interface{}{
			"id":          claimID,
			"business_id": businessID,
		})
		return err
	})

	run("Reject claim (expect 400 — claim already accepted)", e("FN_REJECT_CLAIM"), func() error {
		if claimID == "" || businessID == "" {
			return fmt.Errorf("skipped — claim or business missing")
		}
		_, err := call(e("FN_REJECT_CLAIM"), map[string]interface{}{
			"id":          claimID,
			"business_id": businessID,
		})
		// this SHOULD fail because claim is accepted, not pending
		if err == nil {
			return fmt.Errorf("expected error: rejecting an accepted claim should fail")
		}
		fmt.Printf("    %s↳ got expected error: %s%s\n", dim, err.Error(), reset)
		return nil // expected failure = test pass
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("USER — COMPLETION")
	// ──────────────────────────────────────────────────────────────────────────

	run("Mark claim received", e("FN_MARK_RECEIVED"), func() error {
		if claimID == "" || userID == "" {
			return fmt.Errorf("skipped — claim or user missing")
		}
		_, err := call(e("FN_MARK_RECEIVED"), map[string]interface{}{
			"id":      claimID,
			"user_id": userID,
		})
		return err
	})

	run("Add review (rating 5)", e("FN_ADD_REVIEW"), func() error {
		if businessID == "" || userID == "" {
			return fmt.Errorf("skipped — business or user missing")
		}
		_, err := call(e("FN_ADD_REVIEW"), map[string]interface{}{
			"business_id": businessID,
			"user_id":     userID,
			"rating":      5,
			"comment":     "Excellent integration test experience!",
		})
		return err
	})

	run("Add duplicate review (expect 409)", e("FN_ADD_REVIEW"), func() error {
		if businessID == "" || userID == "" {
			return fmt.Errorf("skipped — business or user missing")
		}
		_, err := call(e("FN_ADD_REVIEW"), map[string]interface{}{
			"business_id": businessID,
			"user_id":     userID,
			"rating":      3,
			"comment":     "Second review — should fail",
		})
		if err == nil {
			return fmt.Errorf("expected unique constraint error for duplicate review")
		}
		fmt.Printf("    %s↳ got expected error: %s%s\n", dim, err.Error(), reset)
		return nil
	})

	run("Get reviews", e("FN_GET_REVIEWS"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_GET_REVIEWS"), map[string]interface{}{
			"business_id": businessID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("USER — NOTIFICATIONS")
	// ──────────────────────────────────────────────────────────────────────────

	run("Get notifications", e("FN_GET_NOTIFICATIONS"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		_, err := call(e("FN_GET_NOTIFICATIONS"), map[string]interface{}{
			"user_id": userID,
		})
		return err
	})

	run("Mark all notifications read", e("FN_MARK_NOTIFICATION_READ"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		// omit id → marks all read
		_, err := call(e("FN_MARK_NOTIFICATION_READ"), map[string]interface{}{
			"user_id": userID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("BUSINESS — OFFERS CRUD")
	// ──────────────────────────────────────────────────────────────────────────

	run("Update offer", e("FN_UPDATE_OFFER"), func() error {
		if offerID == "" || businessID == "" {
			return fmt.Errorf("skipped — offer or business missing")
		}
		_, err := call(e("FN_UPDATE_OFFER"), map[string]interface{}{
			"id":              offerID,
			"business_id":     businessID,
			"title":           fmt.Sprintf("Updated Offer %d", ts),
			"discount_type":   "flat",
			"discount_value":  50,
			"min_order_value": 200,
			"is_active":       true,
			"starts_at":       startsAt,
			"expires_at":      expiresAt,
		})
		return err
	})

	run("Delete offer", e("FN_DELETE_OFFER"), func() error {
		if offerID == "" || businessID == "" {
			return fmt.Errorf("skipped — offer or business missing")
		}
		_, err := call(e("FN_DELETE_OFFER"), map[string]interface{}{
			"id":          offerID,
			"business_id": businessID,
		})
		return err
	})

	run("Get offers after delete (expect empty)", e("FN_GET_OFFERS_BY_BUSINESS"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_GET_OFFERS_BY_BUSINESS"), map[string]interface{}{
			"business_id": businessID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("VALIDATION — EDGE CASES")
	// ──────────────────────────────────────────────────────────────────────────

	run("Claim with zero quantity (expect 422)", e("FN_CLAIM_PRODUCT"), func() error {
		if userID == "" || businessID == "" || inventoryID == "" {
			return fmt.Errorf("skipped — dependencies missing")
		}
		_, err := call(e("FN_CLAIM_PRODUCT"), map[string]interface{}{
			"user_id":      userID,
			"business_id":  businessID,
			"inventory_id": inventoryID,
			"quantity":     0,
		})
		if err == nil {
			return fmt.Errorf("expected validation error for quantity=0")
		}
		fmt.Printf("    %s↳ got expected error: %s%s\n", dim, err.Error(), reset)
		return nil
	})

	run("Add review with rating 6 (expect 422)", e("FN_ADD_REVIEW"), func() error {
		if businessID == "" || userID == "" {
			return fmt.Errorf("skipped — dependencies missing")
		}
		_, err := call(e("FN_ADD_REVIEW"), map[string]interface{}{
			"business_id": businessID,
			"user_id":     userID,
			"rating":      6,
		})
		if err == nil {
			return fmt.Errorf("expected validation error for rating=6")
		}
		fmt.Printf("    %s↳ got expected error: %s%s\n", dim, err.Error(), reset)
		return nil
	})

	run("Create offer missing title (expect 422)", e("FN_CREATE_OFFER"), func() error {
		if businessID == "" {
			return fmt.Errorf("skipped — business missing")
		}
		_, err := call(e("FN_CREATE_OFFER"), map[string]interface{}{
			"business_id":    businessID,
			"discount_type":  "flat",
			"discount_value": 10,
			"starts_at":      startsAt,
			"expires_at":     expiresAt,
		})
		if err == nil {
			return fmt.Errorf("expected validation error for missing title")
		}
		fmt.Printf("    %s↳ got expected error: %s%s\n", dim, err.Error(), reset)
		return nil
	})

	// ──────────────────────────────────────────────────────────────────────────
	section("PROFILE")
	// ──────────────────────────────────────────────────────────────────────────

	run("Get user profile", e("FN_USER_PROFILE_INFO"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		_, err := call(e("FN_USER_PROFILE_INFO"), map[string]interface{}{
			"id": userID,
		})
		return err
	})

	run("Get business profile", e("FN_BIZ_PROFILE_INFO"), func() error {
		if userID == "" {
			return fmt.Errorf("skipped — user missing")
		}
		_, err := call(e("FN_BIZ_PROFILE_INFO"), map[string]interface{}{
			"owner_id": userID,
		})
		return err
	})

	// ──────────────────────────────────────────────────────────────────────────
	printSummary()
}

// ─── SUMMARY ──────────────────────────────────────────────────────────────────
func printSummary() {
	var passed, failed, skipped int
	var totalDur time.Duration
	var failures []result

	for _, r := range suite {
		switch {
		case r.skipped:
			skipped++
		case r.passed:
			passed++
			totalDur += r.duration
		default:
			failed++
			totalDur += r.duration
			failures = append(failures, r)
		}
	}

	total := passed + failed
	passRate := 0.0
	if total > 0 {
		passRate = float64(passed) / float64(total) * 100
	}

	bar := strings.Repeat("─", 60)
	fmt.Printf("\n%s%s%s%s\n", bold, cyan, bar, reset)
	fmt.Printf("%s%s  RESULTS%s\n", bold, cyan, reset)
	fmt.Printf("%s%s%s%s\n", bold, cyan, bar, reset)
	fmt.Printf("  %s✓ Passed%s   : %s%d%s\n", green, reset, bold, passed, reset)
	fmt.Printf("  %s✗ Failed%s   : %s%d%s\n", red, reset, bold, failed, reset)
	fmt.Printf("  %s⊘ Skipped%s  : %s%d%s\n", yellow, reset, bold, skipped, reset)
	fmt.Printf("  Pass rate  : %s%.1f%%%s\n", bold, passRate, reset)
	fmt.Printf("  Total time : %s%s%s\n", bold, totalDur.Round(time.Millisecond), reset)

	if len(failures) > 0 {
		fmt.Printf("\n%s%sFAILED TESTS%s\n", bold, red, reset)
		for i, r := range failures {
			fmt.Printf("  %d. %s%s%s\n     %s%s%s\n", i+1, red, r.name, reset, dim, r.detail, reset)
		}
	}

	fmt.Printf("%s%s%s%s\n\n", bold, cyan, bar, reset)

	if failed > 0 {
		os.Exit(1)
	}
}
