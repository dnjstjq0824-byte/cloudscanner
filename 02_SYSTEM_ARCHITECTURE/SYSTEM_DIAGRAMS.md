# SYSTEM_DIAGRAMS.md — TirpScanner Architecture & State Machine Diagrams

**Document ID:** `CS-ARCH-002`

> **Pivot addendum (2026-11-26):** Diagrams refreshed for TirpScanner. The BYOC worker and cloud management APIs are removed, while KMS envelope encryption now protects affiliate/integration credentials; the live platform ingests SaaS/AI-tool deals and monetizes via affiliate. The BYOC provisioning state machine is preserved below as retired history.

---

## 1. End-to-End Platform Architecture Diagram

```mermaid
graph TB
    subgraph Public_Internet["Public Internet & Crawlers"]
        User["Browser / Mobile User"]
        Bot["Googlebot / Bingbot"]
        AffNet["Affiliate Networks (Impact / CJ / PartnerStack)"]
    end

    subgraph Edge["Cloudflare Enterprise Edge"]
        WAF["WAF + DDoS + Bot Score"]
        CDN["SSR Cache & Static Assets"]
    end

    subgraph App_Tier["TirpScanner Kubernetes Cluster (Independent VPC)"]
        Web["Nuxt 3 SSR Frontend (cs-web-ssr)"]
        API["Go Core API & Redirect Engine (cs-api-gateway)"]
        W_Deals["Deal & Plan-Price Sync Worker (cs-worker-deals)"]
        W_Aff["Affiliate & Outbox Worker (cs-worker-affiliate)"]
    end

    subgraph Data_Tier["Isolated Data & Security Tier"]
        PG[("PostgreSQL 16 HA (RLS + Ledger)")]
        OS[("OpenSearch 2.17 (Deal Facets)")]
        RD[("Redis 7.4 (Asynq + Cache)")]
    end

    subgraph External_Vendors["SaaS Vendors & Deal Sources"]
        P_API["Vendor Pricing/Deal Feeds & Network Product Catalogs"]
    end

    User --> WAF
    Bot --> WAF
    AffNet -->|Signed Postback| WAF
    WAF --> CDN
    CDN --> Web
    CDN --> API
    Web --> API
    API --> PG
    API --> OS
    API --> RD
    W_Deals --> P_API
    W_Deals --> PG
    W_Deals --> OS
    W_Aff --> PG
    W_Aff --> RD
```

---

## 2. Deal Scan & Price-Drop Detection Flow (live)

```mermaid
stateDiagram-v2
    [*] --> FETCHING: Scheduled ingestion run starts
    FETCHING --> HASH_COMPARING: Raw feed payload archived + SHA-256 hashed
    HASH_COMPARING --> VERIFIED_CURRENT: Payload unchanged — bump verified_at
    HASH_COMPARING --> ANOMALY_CHECKING: Payload changed
    ANOMALY_CHECKING --> QUARANTINED: Price swing > 25% or plan count drop > 20%
    QUARANTINED --> ANOMALY_CHECKING: Admin approves or rejects with reason
    ANOMALY_CHECKING --> NORMALIZED: Snapshot validated & plan prices normalized
    NORMALIZED --> PRICE_DROP_MATCHING: Immutable price_snapshot inserted
    PRICE_DROP_MATCHING --> ALERTS_DISPATCHED: Drops matched to price_alerts; notifications sent
    PRICE_DROP_MATCHING --> INDEXED: No drops — OpenSearch deal index updated
    ALERTS_DISPATCHED --> INDEXED: OpenSearch deal index updated
    VERIFIED_CURRENT --> [*]
    INDEXED --> [*]
```

---

## 3. BYOC Safe Server Provisioning State Machine — RETIRED HISTORY (do not implement)

> Retired 2026-11-26 with the TirpScanner pivot. Preserved for audit traceability only.

```mermaid
stateDiagram-v2
    [*] --> DRAFT: User selects SKU, Region, Image, SSH Key
    DRAFT --> AWAITING_COST_APPROVAL: Calculate TCO & Budget Check
    AWAITING_COST_APPROVAL --> CANCELLED: User declines or 5m token expires
    AWAITING_COST_APPROVAL --> APPROVED_QUEUED: User signs Cost Approval + Idempotency-Key locked
    APPROVED_QUEUED --> PREFLIGHT_VERIFYING: Worker decrypts token in-memory with AAD
    PREFLIGHT_VERIFYING --> FAILED_SAFE: Invalid/Expired/Read-Only Token or Quota/Balance Error
    PREFLIGHT_VERIFYING --> DISPATCHING_TO_PROVIDER: Tag with op-id & POST /servers
    DISPATCHING_TO_PROVIDER --> POLLING_ACTIVE: HTTP 201/202 + external_resource_id saved
    DISPATCHING_TO_PROVIDER --> TIMEOUT_RECONCILING: HTTP Timeout / Network Drop
    TIMEOUT_RECONCILING --> POLLING_ACTIVE: Reconcile finds server with op tag (NO duplicate POST!)
    TIMEOUT_RECONCILING --> FAILED_SAFE: Reconcile confirms zero server created at provider
    POLLING_ACTIVE --> SUCCEEDED: Provider reports RUNNING / ACTIVE + IP assigned
    POLLING_ACTIVE --> FAILED_SAFE: Provider reports ERROR / PROVISION_FAILED
```
