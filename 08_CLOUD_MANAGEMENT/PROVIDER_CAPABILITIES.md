# PROVIDER_CAPABILITIES.md — Verified Provider Management Capability Matrix (Article 13)
**Document ID:** `CS-CM-006`

| Provider Slug | Auth Mode | List Regions/Plans/Images | Create Server | Power On/Off/Reboot/Cycle | Resize | Backups & Snapshots | Firewalls & SSH Keys | Tagging / Idempotency | Delete Server |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| `digitalocean` | OAuth2 / Scoped PAT | YES | YES (`POST /v2/droplets`) | YES (`/v2/droplets/{id}/actions`) | YES | YES | YES | YES (`tags[]`) | YES |
| `upcloud` | Sub-Account API Token | YES | YES (`POST /1.3/server`) | YES (`/1.3/server/{uuid}/start\|stop\|restart`) | YES | YES | YES | YES (`labels[]`) | YES |
| `cloudways` | API v2 OAuth Bearer | YES | YES (`POST /api/v2/server`) | YES (`/api/v2/server/manage/state`) | YES (`scaleServer`) | YES | YES | YES (Op Polling) | YES |
| `vultr` | Bearer API Key + IP ACL | YES | YES (`POST /v2/instances`) | YES (`/v2/instances/{id}/halt\|start\|reboot`) | YES | YES | YES | YES (`tags[]`) | YES |
| `hetzner` | Project API Token | YES | YES (`POST /v1/servers`) | YES (`/v1/servers/{id}/actions/*`) | YES (`change_type`) | YES | YES | YES (`labels{}`) | YES |
| `akamai_linode` | OAuth2 / PAT | YES | YES (`POST /v4/linode/instances`) | YES (`/v4/linode/instances/{id}/boot\|shutdown\|reboot`) | YES | YES | YES | YES (`tags[]`) | YES |
| `gcp` | OAuth2 / Workload Identity | YES | YES (`instances.insert`) | YES (`start`, `stop`, `reset`) | YES | YES | YES | YES (`requestId` UUID) | YES |
| `aws` | IAM STS `AssumeRole` | YES | YES (`RunInstances`) | YES (`Start/Stop/RebootInstances`) | YES | YES | YES | YES (`ClientToken` + Tags) | YES |
| `azure` | Entra ID OAuth2 RBAC | YES | YES (`virtualMachines/createOrUpdate`) | YES (`start`, `powerOff`, `restart`) | YES | YES | YES | YES (`tags{}`) | YES |
| `ovhcloud` | OVH Consumer Key / OAuth2 | YES | YES (`POST /cloud/project/{id}/instance`) | YES | YES | YES | YES | YES | YES |
| `scaleway` | IAM Secret Key | YES | YES (`POST /instance/v1/zones/{z}/servers`) | YES | YES | YES | YES | YES (`tags[]`) | YES |
| `contabo` | OAuth2 Token | YES | YES (`POST /v1/compute/instances`) | YES | YES | YES | YES | YES (`x-request-id`) | YES |
| `hostinger` | Bearer PAT | YES | YES (`/api/vps/v1/virtual-machines`) | YES | YES | YES | YES | YES | YES |
| `kamatera` | ClientID + Secret | YES | YES (`POST /service/server`) | YES | YES | YES | YES | YES | YES |
| `ionos` | Bearer Token | YES | YES (`POST /cloudapi/v6/datacenters/{id}/servers`) | YES | YES | YES | YES | YES | YES |
| `oracle_cloud` | OCI RSA Signing Key | YES | YES (`LaunchInstance`) | YES (`InstanceAction`) | YES | YES | YES | YES (`opc-retry-token`) | YES |
| `runpod` | GraphQL Bearer Key | YES | YES (`podFindAndDeployOnDemand`) | YES (`podResume`, `podStop`) | YES | YES | YES | YES | YES (`podTerminate`) |
| `cloudflare` | Scoped API Token | N/A (Edge/R2) | N/A (R2 Bucket / Worker) | N/A | N/A | N/A | YES (WAF) | YES | YES (Bucket) |
