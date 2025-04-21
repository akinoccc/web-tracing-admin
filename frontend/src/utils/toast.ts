import { toast } from 'vue-sonner'

type ToastType = 'success' | 'error' | 'info' | 'warning'

/**
 * 显示提示信息
 * @param message 提示消息
 * @param type 提示类型 'success' | 'error' | 'info' | 'warning'
 */
export function showToast(message: string, type: ToastType = 'info') {
  switch (type) {
    case 'success':
      toast.success(message)
      break
    case 'error':
      toast.error(message)
      break
    case 'warning':
      toast.warning(message)
      break
    case 'info':
    default:
      toast(message)
      break
  }
}

export { toast }