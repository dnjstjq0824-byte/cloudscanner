package affiliate

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrLinkNotFound       = errors.New("affiliate link not found")
	ErrLinkInactive       = errors.New("affiliate link is inactive")
	ErrInvalidSignature   = errors.New("invalid webhook signature or subid")
	ErrUnmatchedClick     = errors.New("click not found for attribution")
	ErrOpenRedirectFailed = errors.New("open redirect protection blocked unapproved destination")
)

// RedirectResult contains the final destination URL and click metadata.
type RedirectResult struct {
	DestinationURL     string
	ClickULID          string
	SubIDSignature     string
	IsBot              bool
	RedirectStatusCode int
	FTCBadgeRequired   bool
}

// PostbackPayload represents an incoming webhook payload from an affiliate network.
type PostbackPayload struct {
	NetworkName          string          `json:"network_name"`
	ExternalConversionID string          `json:"external_conversion_id"`
	SubID1               string          `json:"sub_id1"` // ClickULID
	SubID2               string          `json:"sub_id2"` // HMAC signature
	OrderAmount          decimal.Decimal `json:"order_amount"`
	CommissionAmount     decimal.Decimal `json:"commission_amount"`
	Currency             string          `json:"currency"`
	RewardType           RewardType      `json:"reward_type"`
	RawPayload           string          `json:"raw_payload"`
}

// Service manages affiliate tracking, attribution, and double-entry commission ledgers.
type Service struct {
	mu           sync.RWMutex
	links        map[string]*AffiliateLink
	programs     map[string]*AffiliateProgram
	clicks       map[string]*ClickEvent // keyed by ClickULID
	conversions  map[string]*ConversionEvent
	commissions  map[string]*Commission
	adjustments  []CommissionAdjustment
	signingKey   []byte
	ipSalt       []byte
	trustedHosts map[string]bool
}

// NewService instantiates an in-memory Affiliate Service.
func NewService(signingKey, ipSalt []byte) *Service {
	if len(signingKey) == 0 {
		signingKey = []byte("cloudscanr_default_signing_secret_2026")
	}
	if len(ipSalt) == 0 {
		ipSalt = []byte("cloudscanr_ip_privacy_salt_2026")
	}

	svc := &Service{
		links:        make(map[string]*AffiliateLink),
		programs:     make(map[string]*AffiliateProgram),
		clicks:       make(map[string]*ClickEvent),
		conversions:  make(map[string]*ConversionEvent),
		commissions:  make(map[string]*Commission),
		adjustments:  make([]CommissionAdjustment, 0),
		signingKey:   signingKey,
		ipSalt:       ipSalt,
		trustedHosts: make(map[string]bool),
	}

	// Register trusted provider domains for open-redirect prevention
	defaultTrusted := []string{
		"digitalocean.com", "pxf.io",
		"vultr.com",
		"upcloud.com",
		"hetzner.com", "hetzner.cloud",
		"aws.amazon.com", "amazon.com",
		"cloud.google.com", "google.com",
		"azure.microsoft.com", "microsoft.com",
		"linode.com", "akamai.com",
		"ovhcloud.com", "ovh.com",
		"scaleway.com",
		"cloudways.com",
		"contabo.com",
		"impact.com",
		"cj.com",
		"partnerstack.com",
	}
	for _, h := range defaultTrusted {
		svc.trustedHosts[h] = true
	}

	svc.seedDefaultLinks()
	return svc
}

func (s *Service) seedDefaultLinks() {
	doLink := &AffiliateLink{
		ID:                     uuid.NewString(),
		CampaignID:             uuid.NewString(),
		ProviderID:             "01930000-0001-7000-8000-000000000001",
		ProviderSlug:           "digitalocean",
		Slug:                   "digitalocean",
		DestinationTemplateURL: "https://digitalocean.pxf.io/c/cloudscanr/1/sub?subId1={{click_id}}&subId2={{sig}}",
		FallbackCleanURL:       "https://www.digitalocean.com",
		SubIDParamName:         "subId1",
		HTTPHealthStatus:       200,
		IsActive:               true,
		TotalClicks:            0,
		CreatedAt:              time.Now().UTC(),
	}

	vultrLink := &AffiliateLink{
		ID:                     uuid.NewString(),
		CampaignID:             uuid.NewString(),
		ProviderID:             "01930000-0001-7000-8000-000000000003",
		ProviderSlug:           "vultr",
		Slug:                   "vultr",
		DestinationTemplateURL: "https://www.vultr.com/?ref=cloudscanr&subId1={{click_id}}&subId2={{sig}}",
		FallbackCleanURL:       "https://www.vultr.com",
		SubIDParamName:         "subId1",
		HTTPHealthStatus:       200,
		IsActive:               true,
		TotalClicks:            0,
		CreatedAt:              time.Now().UTC(),
	}

	hetznerLink := &AffiliateLink{
		ID:                     uuid.NewString(),
		CampaignID:             uuid.NewString(),
		ProviderID:             "01930000-0001-7000-8000-000000000004",
		ProviderSlug:           "hetzner",
		Slug:                   "hetzner",
		DestinationTemplateURL: "https://hetzner.cloud/?ref=cloudscanr_ref&subId1={{click_id}}",
		FallbackCleanURL:       "https://www.hetzner.com/cloud",
		SubIDParamName:         "subId1",
		HTTPHealthStatus:       200,
		IsActive:               true,
		TotalClicks:            0,
		CreatedAt:              time.Now().UTC(),
	}

	s.links[doLink.Slug] = doLink
	s.links[vultrLink.Slug] = vultrLink
	s.links[hetznerLink.Slug] = hetznerLink
}

// GenerateULID produces a 26-character Crockford Base32 timestamp-prefixed identifier.
func GenerateULID() string {
	const encoding = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	now := uint64(time.Now().UTC().UnixMilli())
	var b [16]byte
	b[0] = byte(now >> 40)
	b[1] = byte(now >> 32)
	b[2] = byte(now >> 24)
	b[3] = byte(now >> 16)
	b[4] = byte(now >> 8)
	b[5] = byte(now)
	_, _ = rand.Read(b[6:])

	var ulid [26]byte
	// 10 chars for 48-bit timestamp
	ulid[0] = encoding[(b[0]&224)>>5]
	ulid[1] = encoding[b[0]&31]
	ulid[2] = encoding[(b[1]&248)>>3]
	ulid[3] = encoding[((b[1]&7)<<2)|((b[2]&192)>>6)]
	ulid[4] = encoding[(b[2]&62)>>1]
	ulid[5] = encoding[((b[2]&1)<<4)|((b[3]&240)>>4)]
	ulid[6] = encoding[((b[3]&15)<<1)|((b[4]&128)>>7)]
	ulid[7] = encoding[(b[4]&124)>>2]
	ulid[8] = encoding[((b[4]&3)<<3)|((b[5]&224)>>5)]
	ulid[9] = encoding[b[5]&31]

	// 16 chars for 80-bit randomness
	ulid[10] = encoding[(b[6]&248)>>3]
	ulid[11] = encoding[((b[6]&7)<<2)|((b[7]&192)>>6)]
	ulid[12] = encoding[(b[7]&62)>>1]
	ulid[13] = encoding[((b[7]&1)<<4)|((b[8]&240)>>4)]
	ulid[14] = encoding[((b[8]&15)<<1)|((b[9]&128)>>7)]
	ulid[15] = encoding[(b[9]&124)>>2]
	ulid[16] = encoding[((b[9]&3)<<3)|((b[10]&224)>>5)]
	ulid[17] = encoding[b[10]&31]
	ulid[18] = encoding[(b[11]&248)>>3]
	ulid[19] = encoding[((b[11]&7)<<2)|((b[12]&192)>>6)]
	ulid[20] = encoding[(b[12]&62)>>1]
	ulid[21] = encoding[((b[12]&1)<<4)|((b[13]&240)>>4)]
	ulid[22] = encoding[((b[13]&15)<<1)|((b[14]&128)>>7)]
	ulid[23] = encoding[(b[14]&124)>>2]
	ulid[24] = encoding[((b[14]&3)<<3)|((b[15]&224)>>5)]
	ulid[25] = encoding[b[15]&31]

	return string(ulid[:])
}

// SignSubID computes an HMAC-SHA256 signature over clickULID:slug to detect tampering.
func (s *Service) SignSubID(clickULID, slug string) string {
	mac := hmac.New(sha256.New, s.signingKey)
	mac.Write([]byte(clickULID + ":" + slug))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifySubIDSignature validates whether the provided signature matches the clickULID and slug.
func (s *Service) VerifySubIDSignature(clickULID, slug, providedSig string) bool {
	expected := s.SignSubID(clickULID, slug)
	return hmac.Equal([]byte(expected), []byte(providedSig))
}

// HashIP creates a privacy-preserving keyed hash of client IP.
func (s *Service) HashIP(ip string) string {
	mac := hmac.New(sha256.New, s.ipSalt)
	mac.Write([]byte(ip))
	return hex.EncodeToString(mac.Sum(nil))
}

// IsBotAgent checks if user agent string indicates automated crawlers.
func IsBotAgent(ua string) bool {
	l := strings.ToLower(ua)
	botKeywords := []string{
		"bot", "crawler", "spider", "curl", "python-requests",
		"wget", "headlesschrome", "postman", "go-http-client",
		"petalbot", "ahrefsbot", "semrushbot",
	}
	for _, kw := range botKeywords {
		if strings.Contains(l, kw) {
			return true
		}
	}
	return false
}

// isTrustedDestination checks against open redirect vulnerabilities.
func (s *Service) isTrustedDestination(targetURL string) bool {
	u, err := url.Parse(targetURL)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return false
	}
	host := strings.ToLower(u.Hostname())
	for trusted := range s.trustedHosts {
		if host == trusted || strings.HasSuffix(host, "."+trusted) {
			return true
		}
	}
	return false
}

// ResolveRedirect processes an outbound deal link (/go/{slug}), validating and signing.
func (s *Service) ResolveRedirect(ctx context.Context, slug, src, userAgent, clientIP string) (*RedirectResult, error) {
	s.mu.RLock()
	link, exists := s.links[slug]
	s.mu.RUnlock()

	if !exists {
		return nil, ErrLinkNotFound
	}
	if !link.IsActive {
		return nil, ErrLinkInactive
	}

	clickULID := "clk_" + GenerateULID()
	sig := s.SignSubID(clickULID, slug)
	isBot := IsBotAgent(userAgent)

	// Build destination URL by template injection
	target := link.DestinationTemplateURL
	target = strings.ReplaceAll(target, "{{click_id}}", url.QueryEscape(clickULID))
	target = strings.ReplaceAll(target, "{{sig}}", url.QueryEscape(sig))

	// If source is provided, append subId2 or src
	if src != "" && !strings.Contains(target, "subId2=") {
		delim := "?"
		if strings.Contains(target, "?") {
			delim = "&"
		}
		target = fmt.Sprintf("%s%ssubId2=%s", target, delim, url.QueryEscape(src))
	}

	// Validate open redirect
	if !s.isTrustedDestination(target) {
		// Fallback to verified clean URL
		if s.isTrustedDestination(link.FallbackCleanURL) {
			target = link.FallbackCleanURL
		} else {
			return nil, ErrOpenRedirectFailed
		}
	}

	// Asynchronous non-blocking click recording
	deviceType := "DESKTOP"
	if strings.Contains(strings.ToLower(userAgent), "mobile") {
		deviceType = "MOBILE"
	} else if strings.Contains(strings.ToLower(userAgent), "tablet") {
		deviceType = "TABLET"
	}
	if isBot {
		deviceType = "BOT_BLOCKED"
	}

	clickEvent := &ClickEvent{
		ID:                 uuid.NewString(),
		ClickULID:          clickULID,
		LinkID:             link.ID,
		ProviderID:         link.ProviderID,
		ProviderSlug:       link.ProviderSlug,
		SubIDSignature:     sig,
		SourcePagePath:     src,
		PlacementSlot:      "deal_redirect",
		CountryCode:        "US", // Default or MaxMind GeoIP
		DeviceType:         deviceType,
		IPSubnetHMAC:       s.HashIP(clientIP),
		IsBotFiltered:      isBot,
		RedirectStatusCode: 302,
		ClickedAt:          time.Now().UTC(),
	}

	s.mu.Lock()
	s.clicks[clickULID] = clickEvent
	link.TotalClicks++
	s.mu.Unlock()

	return &RedirectResult{
		DestinationURL:     target,
		ClickULID:          clickULID,
		SubIDSignature:     sig,
		IsBot:              isBot,
		RedirectStatusCode: 302,
		FTCBadgeRequired:   true,
	}, nil
}

// IngestPostback handles an external affiliate network postback idempotently.
func (s *Service) IngestPostback(ctx context.Context, payload PostbackPayload, networkSecret string) (*Commission, error) {
	// 1. Look up click attribution
	s.mu.RLock()
	click, exists := s.clicks[payload.SubID1]
	s.mu.RUnlock()

	if !exists {
		return nil, ErrUnmatchedClick
	}

	// 2. Validate signature if provided
	if payload.SubID2 != "" && !s.VerifySubIDSignature(payload.SubID1, click.ProviderSlug, payload.SubID2) {
		return nil, ErrInvalidSignature
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 3. Idempotency check: if conversion ID already recorded, return existing commission
	convKey := payload.NetworkName + ":" + payload.ExternalConversionID
	if existingConv, hasConv := s.conversions[convKey]; hasConv {
		for _, comm := range s.commissions {
			if comm.ConversionID == existingConv.ID {
				return comm, nil
			}
		}
	}

	// 4. Record Conversion
	convID := uuid.NewString()
	conv := &ConversionEvent{
		ID:                   convID,
		ProgramID:            uuid.NewString(),
		ClickID:              &click.ID,
		ClickULID:            payload.SubID1,
		ExternalConversionID: payload.ExternalConversionID,
		ExternalSubID:        payload.SubID1,
		ConversionSource:     "POSTBACK_WEBHOOK",
		ConversionStatus:     ConversionStatusApproved,
		OrderAmount:          payload.OrderAmount,
		OrderCurrency:        payload.Currency,
		RawEventJSONB:        payload.RawPayload,
		ConvertedAt:          time.Now().UTC(),
		ReconciledAt:         time.Now().UTC(),
	}
	s.conversions[convKey] = conv

	// 5. Post Double-Entry Commission
	lockExpiry := time.Now().UTC().AddDate(0, 1, 0) // 30 day lock period
	comm := &Commission{
		ID:                  uuid.NewString(),
		ConversionID:        convID,
		ProgramID:           conv.ProgramID,
		RewardType:          payload.RewardType,
		LedgerStatus:        LedgerStatusApproved,
		CommissionAmount:    payload.CommissionAmount,
		Currency:            payload.Currency,
		NormalizedAmountUSD: payload.CommissionAmount, // 1:1 if USD
		LockExpiresAt:       &lockExpiry,
		CreatedAt:           time.Now().UTC(),
	}
	s.commissions[comm.ID] = comm

	return comm, nil
}

// AdjustCommission appends an immutable adjustment (clawback, bonus, fx variance) to a commission.
func (s *Service) AdjustCommission(ctx context.Context, adj CommissionAdjustment) (*CommissionAdjustment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	comm, exists := s.commissions[adj.CommissionID]
	if !exists {
		return nil, errors.New("target commission not found")
	}

	adj.ID = uuid.NewString()
	adj.CreatedAt = time.Now().UTC()
	s.adjustments = append(s.adjustments, adj)

	// If clawback/refund exceeds original commission, void ledger
	newAmount := comm.CommissionAmount.Add(adj.DeltaAmount)
	if newAmount.LessThanOrEqual(decimal.Zero) {
		comm.LedgerStatus = LedgerStatusVoided
	}

	return &adj, nil
}

// GetDisclosurePolicy returns transparent regulatory compliance and editorial policy.
func (s *Service) GetDisclosurePolicy() *DisclosurePolicy {
	return &DisclosurePolicy{
		PolicyVersion:          "2026.1",
		PlatformBrand:          "CloudScanr (cloudscanr.com)",
		IndependentComparison:  true,
		PayForPlacementAllowed: false,
		EditorialIndependence:  "CloudScanr maintains strict algorithmic separation between cloud pricing benchmarks and affiliate partnerships. Rankings are calculated strictly using mathematical Total Cost of Ownership (TCO) models, and no provider can pay to alter rank order or benchmark results.",
		CompensationNotice:     "CloudScanr may receive referral commissions or hosting credits when users click through verified partner links and deploy services. This occurs at zero additional cost to the user and supports ongoing platform research and benchmark automation.",
		KFTCNoticeKR:           "공정거래위원회 '추천·보증 등에 관한 표시·광고 심사지침' 준수: 본 플랫폼의 일부 링크는 제휴 링크를 포함하고 있으며, 이용자가 서비스를 신청할 경우 플랫폼에 소정의 수수료 또는 호스팅 크레딧이 지급될 수 있습니다. 단, 모든 가격 비교 및 추천 순위는 객관적 알고리즘(TCO 및 성능 지수)에 의해서만 산출되며 금전적 대가로 순위가 변동되지 않습니다.",
		LastAuditedAt:          time.Date(2026, 10, 28, 0, 0, 0, 0, time.UTC),
	}
}

// GetLinkBySlug returns the affiliate link details.
func (s *Service) GetLinkBySlug(slug string) (*AffiliateLink, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	link, exists := s.links[slug]
	return link, exists
}

// GetClickByULID looks up a recorded click event.
func (s *Service) GetClickByULID(ulid string) (*ClickEvent, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, exists := s.clicks[ulid]
	return c, exists
}
