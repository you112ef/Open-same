import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'

// Create axios instance
export const api: AxiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    // Add auth token if available
    const token = localStorage.getItem('access_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // Add request ID for tracking
    config.headers['X-Request-ID'] = generateRequestId()

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config

    // Handle 401 errors (unauthorized)
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        // Try to refresh token
        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          const response = await axios.post(
            `${process.env.REACT_APP_API_URL || 'http://localhost:8080'}/api/v1/auth/refresh`,
            { refresh_token: refreshToken }
          )

          const { access_token } = response.data.data
          localStorage.setItem('access_token', access_token)
          api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`

          // Retry original request
          originalRequest.headers.Authorization = `Bearer ${access_token}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        // Refresh failed, clear tokens
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        delete api.defaults.headers.common['Authorization']
        
        // Redirect to login
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    // Handle other errors
    if (error.response?.status === 403) {
      // Forbidden - user doesn't have permission
      console.error('Access forbidden:', error.response.data)
    } else if (error.response?.status === 404) {
      // Not found
      console.error('Resource not found:', error.response.data)
    } else if (error.response?.status >= 500) {
      // Server error
      console.error('Server error:', error.response.data)
    }

    return Promise.reject(error)
  }
)

// Generate unique request ID
const generateRequestId = (): string => {
  return Math.random().toString(36).substring(2) + Date.now().toString(36)
}

// API endpoints
export const endpoints = {
  // Authentication
  auth: {
    login: '/api/v1/auth/login',
    register: '/api/v1/auth/register',
    refresh: '/api/v1/auth/refresh',
    logout: '/api/v1/auth/logout',
  },

  // User management
  user: {
    profile: '/api/v1/user/profile',
    updateProfile: '/api/v1/user/profile',
    deleteAccount: '/api/v1/user/account',
  },

  // Content management
  content: {
    create: '/api/v1/content',
    list: '/api/v1/content',
    get: (id: string) => `/api/v1/content/${id}`,
    update: (id: string) => `/api/v1/content/${id}`,
    delete: (id: string) => `/api/v1/content/${id}`,
    public: '/api/v1/content/public',
  },

  // Collaboration
  collaboration: {
    list: '/api/v1/collaborations',
    update: (id: string) => `/api/v1/collaborations/${id}`,
    delete: (id: string) => `/api/v1/collaborations/${id}`,
    addCollaborator: (id: string) => `/api/v1/content/${id}/collaborate`,
    shareContent: (id: string) => `/api/v1/content/${id}/share`,
  },

  // AI services
  ai: {
    generate: '/api/v1/ai/generate',
    improve: '/api/v1/ai/improve',
    summarize: '/api/v1/ai/summarize',
  },

  // WebSocket
  websocket: process.env.REACT_APP_WS_URL || 'ws://localhost:8080/ws',
}

// API helper functions
export const apiHelpers = {
  // Handle API errors
  handleError: (error: any): string => {
    if (error.response?.data?.message) {
      return error.response.data.message
    }
    if (error.message) {
      return error.message
    }
    return 'An unexpected error occurred'
  },

  // Handle API responses
  handleResponse: <T>(response: AxiosResponse<T>): T => {
    return response.data
  },

  // Create query parameters
  createQueryParams: (params: Record<string, any>): string => {
    const searchParams = new URLSearchParams()
    
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        if (Array.isArray(value)) {
          value.forEach(item => searchParams.append(key, item.toString()))
        } else {
          searchParams.append(key, value.toString())
        }
      }
    })

    return searchParams.toString()
  },

  // Pagination helper
  createPaginationParams: (page: number, perPage: number, filters?: Record<string, any>) => {
    const params: Record<string, any> = {
      page,
      per_page: perPage,
    }

    if (filters) {
      Object.assign(params, filters)
    }

    return params
  },
}

// Export default api instance
export default api