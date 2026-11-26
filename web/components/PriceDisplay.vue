<template>
  <span class="tabular-nums" :title="quote.tip">
    {{ quote.converted }}
    <span class="opacity-70 text-[0.9em]">({{ quote.native }})</span>
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{
  amount: number
  currency: string
  cycle?: 'month' | 'year'
}>()

const { dual } = useMoney()
const quote = computed(() => dual(props.amount, props.currency, props.cycle))
</script>
