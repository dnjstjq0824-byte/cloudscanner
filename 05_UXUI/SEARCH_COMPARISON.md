# SEARCH_COMPARISON.md — Global Search, Natural Language Bar & Side-by-Side Comparison UX

**Document ID:** `CS-UX-003`

---

## 1. Command-K Natural Language Search Bar
- Typing `4 vcpu 16gb nvme frankfurt under $35` renders interactive removable filter chips below the input:
  `[vCPU >= 4 ×]` `[RAM >= 16 GB ×]` `[Storage: NVMe ×]` `[Location: Frankfurt, DE ×]` `[Max TCO <= $35/mo ×]`
- Users can toggle between **Engineer Data-Grid Mode** (dense 14-column sortable table modeled after `instances.vantage.sh` + extra TCO columns) and **Guided Card Mode** (highlighting Value Score, Pros/Cons, and Trial Credits).

## 2. Side-by-Side Multi-SKU Comparison Matrix (`/compare`)
- Supports pinning up to **6 SKUs** simultaneously across different providers and native currencies.
- Includes **"Highlight Differences Only"** toggle and **"Winner Per Row"** subtle green badge (lowest TCO, highest vCPU, highest included bandwidth, lowest egress overage rate).
