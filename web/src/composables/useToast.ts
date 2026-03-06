import { ref } from 'vue'

type ToastType = 'success' | 'error' | 'warning' | 'info'
type ToastNotifier = (message: string, duration?: number) => number

interface Toast {
  id: number
  message: string
  type: ToastType
  duration: number
}

export interface ToastController {
  success: ToastNotifier
  error: ToastNotifier
  warn: ToastNotifier
  info: ToastNotifier
}

const MAX_TOASTS = 5
const toasts = ref<Toast[]>([])
let toastIdCounter = 0

export function useToast() {
  function addToast(message: string, type: ToastType = 'info', duration: number = 3000) {
    const id = ++toastIdCounter
    const toast: Toast = { id, message, type, duration }

    // 限制最大同时显示 5 个 toast，超过则移除最旧的
    if (toasts.value.length >= MAX_TOASTS) {
      toasts.value.shift()
    }

    toasts.value.push(toast)

    if (duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, duration)
    }

    return id
  }

  function removeToast(id: number) {
    const index = toasts.value.findIndex((toast) => toast.id === id)
    if (index > -1) {
      toasts.value.splice(index, 1)
    }
  }

  const toast: ToastController = {
    success(message, duration) {
      return addToast(message, 'success', duration)
    },
    error(message, duration) {
      return addToast(message, 'error', duration)
    },
    warn(message, duration) {
      return addToast(message, 'warning', duration)
    },
    info(message, duration) {
      return addToast(message, 'info', duration)
    },
  }

  function clear() {
    toasts.value = []
  }

  return { toasts, addToast, removeToast, toast, clear }
}
