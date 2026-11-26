<template>
  <Teleport to="body">
    <div
      v-if="open"
      class="fixed inset-0 z-[100] flex items-end sm:items-center justify-center p-4 bg-black/60"
      role="dialog"
      aria-modal="true"
      :aria-label="title"
      @click.self="$emit('cancel')"
      @keydown.esc="$emit('cancel')"
    >
      <div class="card w-full max-w-md p-6 space-y-4 shadow-pop">
        <h3 class="text-base font-bold text-ink">{{ title }}</h3>
        <p class="text-sm text-ink-secondary">{{ message }}</p>
        <slot />
        <div v-if="error" class="text-xs text-rose-500" role="alert">{{ error }}</div>
        <div class="flex justify-end gap-2 pt-1">
          <button class="btn-ghost !py-2 !text-xs" :disabled="busy" @click="$emit('cancel')">{{ $t('common.cancel') }}</button>
          <button :class="danger ? 'btn-danger' : 'btn-primary'" class="!py-2 !text-xs" :disabled="busy" @click="$emit('confirm')">
            {{ busy ? $t('common.working') : (confirmLabel || $t('common.confirm')) }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{ open: boolean; title: string; message: string; confirmLabel?: string; danger?: boolean; busy?: boolean; error?: string }>(),
  { confirmLabel: '', danger: false, busy: false, error: '' }
)
defineEmits<{ (e: 'confirm'): void; (e: 'cancel'): void }>()
</script>
