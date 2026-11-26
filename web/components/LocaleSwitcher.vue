<template>
  <div class="flex items-center gap-2 text-xs">
    <label class="flex items-center gap-1 text-ink-secondary">
      <span class="hidden xl:inline">{{ $t('selector.language') }}</span>
      <select :value="locale" class="input !w-auto !py-1 !text-xs" :aria-label="$t('selector.language')" @change="switchLocale(($event.target as HTMLSelectElement).value)">
        <option v-for="loc in locales" :key="loc.code" :value="loc.code">{{ loc.name }}</option>
      </select>
    </label>
    <label class="flex items-center gap-1 text-ink-secondary">
      <span class="hidden xl:inline">{{ $t('selector.currency') }}</span>
      <select v-model="market.currency" class="input !w-auto !py-1 !text-xs" :aria-label="$t('selector.currency')" @change="market.setCurrency(market.currency)">
        <option v-for="code in currencies" :key="code" :value="code">{{ code }}</option>
      </select>
    </label>
    <label class="flex items-center gap-1 text-ink-secondary" :title="$t('selector.taxDisclaimer')">
      <span class="hidden xl:inline">{{ $t('selector.tax') }}</span>
      <select
        :value="market.taxIncluded ? 'incl' : 'excl'"
        class="input !w-auto !py-1 !text-xs"
        :aria-label="$t('selector.tax')"
        @change="market.setTaxIncluded(($event.target as HTMLSelectElement).value === 'incl')"
      >
        <option value="excl">{{ $t('selector.taxExcluded') }}</option>
        <option value="incl">{{ $t('selector.taxIncluded') }}</option>
      </select>
    </label>
  </div>
</template>

<script setup lang="ts">
const { locale, locales, setLocale } = useI18n()
const localePath = useLocalePath()
const switchLocalePath = useSwitchLocalePath()
const market = useMarketStore()

const currencies = ['USD', 'EUR', 'KRW', 'JPY', 'GBP']

type LocaleCode = 'en' | 'ko' | 'ja' | 'de' | 'es'

async function switchLocale(code: string) {
  // The select only offers valid codes; narrow for the typed i18n APIs.
  const target = code as LocaleCode
  await setLocale(target)
  await navigateTo(switchLocalePath(target) || localePath('/'))
}

onMounted(() => {
  market.fetchFx()
})
</script>
