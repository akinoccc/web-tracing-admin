import axios, { AxiosInstance, AxiosRequestConfig } from 'axios'
import { useRouter } from 'vue-router'

export function useApi() {
  const router = useRouter()
  
  // 创建 axios 实例
  const api: AxiosInstance = axios.create({
    baseURL: '/',
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json'
    }
  })
  
  // 请求拦截器
  api.interceptors.request.use(
    (config) => {
      // 从 localStorage 获取 token
      const token = localStorage.getItem('token')
      
      // 如果有 token，添加到请求头
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`
      }
      
      return config
    },
    (error) => {
      return Promise.reject(error)
    }
  )
  
  // 响应拦截器
  api.interceptors.response.use(
    (response) => {
      return response
    },
    (error) => {
      // 处理 401 错误
      if (error.response && error.response.status === 401) {
        // 清除 token
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        
        // 跳转到登录页
        router.push('/auth/login')
      }
      
      return Promise.reject(error)
    }
  )
  
  // 封装 GET 请求
  const get = (url: string, params?: any, config?: AxiosRequestConfig) => {
    return api.get(url, { params, ...config })
  }
  
  // 封装 POST 请求
  const post = (url: string, data?: any, config?: AxiosRequestConfig) => {
    return api.post(url, data, config)
  }
  
  // 封装 PUT 请求
  const put = (url: string, data?: any, config?: AxiosRequestConfig) => {
    return api.put(url, data, config)
  }
  
  // 封装 DELETE 请求
  const del = (url: string, config?: AxiosRequestConfig) => {
    return api.delete(url, config)
  }
  
  return {
    api,
    get,
    post,
    put,
    delete: del
  }
}
