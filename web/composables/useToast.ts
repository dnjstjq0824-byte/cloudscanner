export interface Toast {
  id: number
  kind: 'success' | 'error' | 'info'
  message: string
}

let nextId = 1

/** Minimal global toast queue rendered by ToastHost. */
export function useToast() {
  const toasts = useState<Toast[]>('cs-toasts', () => [])

  function push(kind: Toast['kind'], message: string, ttlMs = 4200) {
    const id = nextId++
    toasts.value = [...toasts.value, { id, kind, message }]
    if (import.meta.client) {
      window.setTimeout(() => dismiss(id), ttlMs)
    }
  }

  function dismiss(id: number) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  return {
    toasts,
    success: (message: string) => push('success', message),
    error: (message: string) => push('error', message),
    info: (message: string) => push('info', message),
    dismiss,
  }
}
