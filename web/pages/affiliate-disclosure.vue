<template>
  <div class="max-w-4xl mx-auto space-y-8 py-4">
    <div class="space-y-3">
      <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-slate-800 border border-slate-700 text-slate-300 text-xs font-semibold">
        <span>Policy Version 2026.1 • Regulatory Compliance</span>
      </div>
      <h1 class="text-3xl font-extrabold text-white tracking-tight">Affiliate & Editorial Independence Disclosure</h1>
      <p class="text-sm text-slate-400">
        In compliance with the U.S. Federal Trade Commission (FTC) 16 CFR Part 255 and the Korea Fair Trade Commission (KFTC) 심사지침.
      </p>
    </div>

    <!-- Core Disclosures -->
    <div class="space-y-6 text-sm text-slate-300 leading-relaxed bg-slate-900/60 p-8 rounded-2xl border border-slate-800">
      <section class="space-y-2">
        <h2 class="text-base font-bold text-white">1. Strict Algorithmic Separation (Article 10 Compliance)</h2>
        <p>
          CloudScanr operates with complete algorithmic separation between our cloud pricing benchmarks, Total Cost of Ownership (TCO) calculations, and commercial affiliate partnerships. <strong>No cloud provider can pay to alter their ranking, score, or placement order in our search and comparison engine.</strong> All rankings are computed deterministically from objective hardware specifications and public rate cards.
        </p>
      </section>

      <section class="space-y-2 border-t border-slate-800 pt-6">
        <h2 class="text-base font-bold text-white">2. Affiliate Relationships & Outbound Deal Routing</h2>
        <p>
          When you click certain outbound buttons (such as "Deploy Deal" or links under <code class="text-sky-400 bg-slate-800 px-1.5 py-0.5 rounded">/go/{slug}</code>) and register or launch a server, CloudScanr may receive an affiliate referral commission or hosting credit from the respective cloud provider.
        </p>
        <p class="text-slate-400 text-xs">
          This referral compensation comes at <strong>zero additional cost to you</strong>. In many cases, our partner links provide special promotional signup credits (e.g. $200 free trial credits for new accounts) that reduce your starting infrastructure expense.
        </p>
      </section>

      <section class="space-y-2 border-t border-slate-800 pt-6">
        <h2 class="text-base font-bold text-white">3. 공정거래위원회 '추천·보증 등에 관한 표시·광고 심사지침' 고지</h2>
        <p class="text-slate-300">
          본 사이트의 일부 상품 링크 및 제휴 링크(<code class="text-sky-400 bg-slate-800 px-1.5 py-0.5 rounded">/go/*</code>)를 통해 클라우드 서비스 가입 또는 서버 배포가 이루어지는 경우, CloudScanr는 해당 클라우드 공급자(예: DigitalOcean, Vultr, Hetzner 등)로부터 일정액의 수수료 또는 서비스 이용 크레딧을 제공받을 수 있습니다.
        </p>
        <p class="text-slate-400 text-xs">
          단, 플랫폼 내 모든 가격 정보, 인스턴스 성능 지수, 36개 차원 TCO 총소유비용 산출 및 추천 순위는 객관적 하드웨어 규격과 수학적 알고리즘에 의해서만 결정되며, 경제적 이해관계에 따라 순위가 변경되거나 왜곡되지 않습니다.
        </p>
      </section>

      <section class="space-y-2 border-t border-slate-800 pt-6">
        <h2 class="text-base font-bold text-white">4. Privacy-Preserving Attribution</h2>
        <p class="text-slate-400 text-xs">
          To prevent click fraud without tracking personal browsing history, CloudScanr hashes client IP subnets using SHA-256 with rotating salts. We do not store raw IP addresses or sell user data to advertising brokers.
        </p>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'public'
})
</script>
