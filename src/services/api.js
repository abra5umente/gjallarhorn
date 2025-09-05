import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || (import.meta.env.DEV ? '/api' : '/api')

const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
})

// Request interceptor
axiosInstance.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    const message = error.response?.data?.error || error.message || 'An error occurred'
    return Promise.reject(new Error(message))
  }
)

// Service API
export const serviceApi = {
  getServices: () => axiosInstance.get('/services'),
  createService: (data) => axiosInstance.post('/services', data),
  updateService: (id, data) => axiosInstance.put(`/services/${id}`, data),
  deleteService: (id) => axiosInstance.delete(`/services/${id}`),
  getServiceStatus: (id) => axiosInstance.get(`/services/${id}/status`),
}

// Notification API
export const notificationApi = {
  getConfig: () => axiosInstance.get('/notifications/config'),
  updateConfig: (data) => axiosInstance.post('/notifications/config', data),
}

// Combined API object
export const api = {
  ...serviceApi,
  getNotificationConfig: notificationApi.getConfig,
  updateNotificationConfig: notificationApi.updateConfig,
}

export default api
