package affiliate

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

// RedirectHandler processes high-speed outbound deal clicks (/go/{slug}).
func RedirectHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Extract slug from path: /go/{slug}
		path := strings.TrimPrefix(r.URL.Path, "/go/")
		slug := strings.Split(path, "/")[0]
		if slug == "" {
			http.Error(w, "Affiliate deal slug required", http.StatusBadRequest)
			return
		}

		src := r.URL.Query().Get("src")
		ua := r.UserAgent()
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
		if comma := strings.Index(ip, ","); comma != -1 {
			ip = strings.TrimSpace(ip[:comma])
		}

		res, err := svc.ResolveRedirect(r.Context(), slug, src, ua, ip)
		if err != nil {
			if err == ErrLinkNotFound {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// Set transparent FTC/KFTC affiliate disclosure header
		w.Header().Set("X-Affiliate-Disclosure", "independent-comparison")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Location", res.DestinationURL)
		w.WriteHeader(http.StatusFound)
	}
}

// DisclosureHandler returns transparent regulatory policies and editorial independence rules.
func DisclosureHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		policy := svc.GetDisclosurePolicy()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": policy,
		})
	}
}

// WebhookPostbackRequest defines expected fields in JSON postback payload.
type WebhookPostbackRequest struct {
	ExternalConversionID string  `json:"external_conversion_id"`
	SubID1               string  `json:"sub_id1"`
	SubID2               string  `json:"sub_id2"`
	OrderAmount          float64 `json:"order_amount"`
	CommissionAmount     float64 `json:"commission_amount"`
	Currency             string  `json:"currency"`
	RewardType           string  `json:"reward_type"`
}

// WebhookHandler processes network postback notifications.
func WebhookHandler(svc *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Path: /api/v1/webhooks/affiliate/{network}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		network := "impact"
		if len(parts) >= 5 {
			network = parts[4]
		}

		var req WebhookPostbackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid postback JSON payload", http.StatusBadRequest)
			return
		}

		rewardType := RewardTypeCash
		if req.RewardType == string(RewardTypeCredit) {
			rewardType = RewardTypeCredit
		}

		payload := PostbackPayload{
			NetworkName:          network,
			ExternalConversionID: req.ExternalConversionID,
			SubID1:               req.SubID1,
			SubID2:               req.SubID2,
			OrderAmount:          decimal.NewFromFloat(req.OrderAmount),
			CommissionAmount:     decimal.NewFromFloat(req.CommissionAmount),
			Currency:             req.Currency,
			RewardType:           rewardType,
			RawPayload:           "{}",
		}

		networkSecret := r.Header.Get("X-Affiliate-Network-Signature")
		comm, err := svc.IngestPostback(r.Context(), payload, networkSecret)
		if err != nil {
			if err == ErrUnmatchedClick {
				http.Error(w, "Attribution click not found", http.StatusNotFound)
				return
			}
			if err == ErrInvalidSignature {
				http.Error(w, "Invalid SubID signature", http.StatusForbidden)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "reconciled",
			"data":   comm,
		})
	}
}

