# COST_MANAGEMENT.md — TirpScanner Own Infrastructure FinOps & Test Account Budget Controls
**Document ID:** `CS-DO-004`
> TirpScanner pivot (2026-11-26): live cloud-VM sandbox testing retired; FinOps now covers deal-scanner infrastructure and CI only. See `TIRPSCANNER_PIVOT.md`.
## 1. Test Infrastructure Budget Cap
- CI ephemeral infrastructure enforces a hard **$5.00/day** and **$25.00/month** ceiling, with a planned automated janitor cron (`cmd/janitor`, unimplemented) that would terminate any test resource tagged `tirpscanner-ci-ephemeral=true` older than 10 minutes. Until implemented, purge ephemeral resources manually after CI runs.
