/**
 * Toast notification composable
 */

import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: number
  message: string
  type: ToastType
  duration: number
  createdAt: number
}

const toasts = ref<Toast[]>([])
let toastIdCounter = 0

export function useToast() {
  function addToast(message: string, type: ToastType = 'info', duration: number = 3000) {
    const id = ++toastIdCounter
    const toast: Toast = {
      id,
      message,
      type,
      duration,
      createdAt: Date.now(),
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

  function success(message: string, duration?: number) {
    return addToast(message, 'success', duration)
  }

  function error(message: string, duration?: number) {
    return addToast(message, 'error', duration)
  }

  function warning(message: string, duration?: number) {
    return addToast(message, 'warning', duration)
  }

  function info(message: string, duration?: number) {
    return addToast(message, 'info', duration)
  }

  function clear() {
    toasts.value = []
  }

  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info,
    clear,
  }
}
