# SYSTEM_DIAGRAMS.md — CloudScanr Architecture & State Machine Diagrams

**Document ID:** `CS-ARCH-002`

---

## 1. End-to-End Platform Architecture Diagram

```mermaid
graph TB
    subgraph Public_Internet["Public Internet & Crawlers"]
        User["Browser / Mobile User"]
        Bot["Googlebot / Bingbot / TreeScanbot"]
        AffNet["Affiliate Networks (Impact / CJ / PartnerStack)"]
    end

    subgraph Edge["Cloudflare Enterprise Edge"]
        WAF["WAF + DDoS + Bot Score"]
        CDN["SSR Cache & Static Assets"]
    end

    subgraph App_Tier["CloudScanr Kubernetes Cluster (Independent VPC)"]
        Web["Nuxt 3 SSR Frontend (cs-web-ssr)"]
        API["Go Core API & Redirect Engine (cs-api-gateway)"]
        W_Price["Pricing Sync Worker (cs-worker-pricing)"]
        W_Cloud["BYOC Cloud Ops Worker (cs-worker-cloud)"]
        W_Aff["Affiliate & Outbox Worker (cs-worker-affiliate)"]
    end

    subgraph Data_Tier["Isolated Data & Security Tier"]
        PG[("PostgreSQL 16 HA (RLS + Ledger)")]
        OS[("OpenSearch 2.17 (SKU Facets)")]
        RD[("Redis 7.4 (Asynq + Cache)")]
        KMS["AWS KMS / Vault (KEK Root)"]
    end

    subgraph External_Clouds["18+ External Cloud Providers"]
        P_API["Official Pricing APIs (DO, Vultr, Hetzner, Linode, GCP, AWS, Azure...)"]
        M_API["User BYOC Management APIs (Droplets, UpCloud, Cloudways v2...)"]
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
    API --> KMS
    W_Price --> P_API
    W_Price --> PG
    W_Price --> OS
    W_Cloud --> KMS
    W_Cloud --> M_API
    W_Cloud --> PG
    W_Aff --> PG
    W_Aff --> RD
```

---

## 2. BYOC Safe Server Provisioning & Timeout Reconciliation State Machine (Article 13 & 24)

```mermaid
stateDiagram-v2
    [*] --> DRAFT: User selects SKU, Region, Image, SSH Key
    DRAFT --> AWAITING_COST_APPROVAL: Calculate 36-field TCO & Budget Check
    AWAITING_COST_APPROVAL --> CANCELLED: User declines or 5m token expires
    AWAITING_COST_APPROVAL --> APPROVED_QUEUED: User signs Cost Approval + Idempotency-Key locked
    APPROVED_QUEUED --> PREFLIGHT_VERIFYING: Worker decrypts token in-memory with AAD
    PREFLIGHT_VERIFYING --> FAILED_SAFE: Invalid/Expired/Read-Only Token or Quota/Balance Error
    PREFLIGHT_VERIFYING --> DISPATCHING_TO_PROVIDER: Tag with cloudscanr-op-id={uuid} & POST /servers
    DISPATCHING_TO_PROVIDER --> POLLING_ACTIVE: HTTP 201/202 + external_resource_id saved
    DISPATCHING_TO_PROVIDER --> TIMEOUT_RECONCILING: HTTP Timeout / Network Drop
    TIMEOUT_RECONCILING --> POLLING_ACTIVE: Reconcile finds server with tag cloudscanr-op-id (NO duplicate POST!)
    TIMEOUT_RECONCILING --> FAILED_SAFE: Reconcile confirms zero server created at provider
    POLLING_ACTIVE --> SUCCEEDED: Provider reports RUNNING / ACTIVE + IP assigned
    POLLING_ACTIVE --> FAILED_SAFE: Provider reports ERROR / PROVISION_FAILED
```
