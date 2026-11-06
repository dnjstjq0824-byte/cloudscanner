<template>
  <div class="max-w-3xl mx-auto space-y-8 py-8">
    <div class="space-y-2">
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Privacy Policy</h1>
      <p class="text-xs text-slate-500">Last updated: November 1, 2026 · CloudScanr (cloudscanr.com)</p>
    </div>

    <div class="space-y-6 text-sm text-slate-300 leading-relaxed">
      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">1. Data we collect</h2>
        <ul class="list-disc list-inside space-y-1">
          <li>Account data: email, display name, Argon2id password hash (never plain-text passwords).</li>
          <li>Organization and role membership for tenant isolation.</li>
          <li>BYOC credential metadata: KMS key reference, encrypted blobs, token fingerprints, granted scopes. Raw provider tokens are sealed at submission and never logged.</li>
          <li>Operational telemetry: request logs with IP-derived HMAC hashes (not raw IPs) for affiliate click attribution and abuse prevention.</li>
        </ul>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">2. How we use data</h2>
        <p>We use data to operate the comparison service, attribute affiliate conversions, execute only your explicitly approved cloud operations, enforce budgets and safety gates, and prevent fraud. We do not sell personal data and we do not use your infrastructure data for advertising.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">3. Data sharing</h2>
        <p>Data is shared only with infrastructure subprocessors required to run the service (hosting, database, email delivery) and with affiliate networks in the form of pseudonymous click identifiers (SubID/ClickID) when you click an outbound deal link. No account credentials are ever shared.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">4. Security</h2>
        <p>Passwords use Argon2id; sessions are short-lived signed tokens; BYOC secrets use KMS envelope encryption with per-connection authentication context; tenant data is isolated per organization and enforced at the database layer. Security incidents follow the published incident runbook.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">5. Retention and deletion</h2>
        <p>Account data is retained while your account is active. Disconnecting a cloud account immediately revokes and wipes its stored credential blobs. You may request full account and data deletion at privacy@cloudscanr.com; audit logs required for financial compliance are retained in anonymized form where the law requires.</p>
      </section>

      <section class="space-y-2">
        <h2 class="text-lg font-bold text-white">6. Your rights</h2>
        <p>Depending on your jurisdiction (including GDPR and Korean PIPA rights), you may request access, correction, export, or deletion of your personal data. Contact privacy@cloudscanr.com and we will respond within 30 days.</p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'public',
})

useHead({
  title: 'Privacy Policy',
  meta: [{ name: 'description', content: 'CloudScanr Privacy Policy: data collection, BYOC credential handling, retention, and your rights.' }],
})
</script>
