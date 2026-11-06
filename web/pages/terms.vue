<template>
  <div class="max-w-3xl mx-auto space-y-8 py-8">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Terms of Service</h1>
      <p class="text-xs text-slate-500">Last updated: November 1, 2026 · CloudScanr (cloudscanr.com)</p>
    </div>

    <div class="space-y-6 text-sm text-slate-300 leading-relaxed">
      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">1. The service</h2>
        <p>CloudScanr provides independent cloud price comparison, total-cost-of-ownership estimation, affiliate-linked provider referrals, and Bring-Your-Own-Cloud (BYOC) server lifecycle management. Comparison rankings are computed by deterministic algorithms only; providers cannot pay for placement.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">2. Your cloud accounts (BYOC)</h2>
        <p>You connect cloud accounts that you own or are authorized to administer, using OAuth or scoped API tokens. Tokens are sealed with envelope encryption and used solely to perform actions you explicitly approve. You remain the cloud customer of record: all infrastructure charges are billed by your provider, never by CloudScanr, and CloudScanr applies zero markup.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">3. Billable and destructive actions</h2>
        <p>Server creation, resize, backup enrollment, rebuild, and deletion require an explicit signed cost approval showing estimated hourly and monthly charges. Destructive actions additionally require hostname confirmation and step-up verification. Approvals expire after 15 minutes.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">4. Affiliate disclosure</h2>
        <p>Some outbound links are affiliate links. If you purchase through them, CloudScanr may receive a commission or hosting credit at no extra cost to you. See the <NuxtLink to="/affiliate-disclosure" class="text-sky-400 hover:text-sky-300">Affiliate Disclosure</NuxtLink> for the full policy.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">5. Acceptable use</h2>
        <p>You must not abuse the service, circumvent rate limits or safety gates, probe other tenants' data, or use the platform to operate infrastructure you are not authorized to manage. We may suspend accounts that violate these terms or threaten platform integrity.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">6. Estimates are not invoices</h2>
        <p>TCO calculations are transparent estimates based on provider rate cards, scenario inputs, and freshness indicators. They are not binding quotes; your provider's metering is authoritative.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">7. Limitation of liability</h2>
        <p>To the maximum extent permitted by law, CloudScanr is not liable for indirect damages, provider outages, or infrastructure charges incurred on your connected accounts. The service is provided "as is" without warranties of uninterrupted availability.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">8. Contact</h2>
        <p>Questions about these terms: legal@cloudscanr.com.</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'public',
})

useHead({
  title: 'Terms of Service',
  meta: [{ name: 'description', content: 'CloudScanr Terms of Service: BYOC responsibilities, approval gates, affiliate disclosure, and acceptable use.' }],
})
</script>
