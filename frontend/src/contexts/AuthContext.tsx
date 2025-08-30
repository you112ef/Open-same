import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@services/api'
import { User, AuthResponse } from '@types/models'

interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, username: string, password: string, firstName?: string, lastName?: string) => Promise<void>
  logout: () => void
  refreshToken: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const queryClient = useQueryClient()

  // Check if user is authenticated on mount
  useEffect(() => {
    const token = localStorage.getItem('access_token')
    if (token) {
      // Verify token and get user info
      api.defaults.headers.common['Authorization'] = `Bearer ${token}`
      checkAuthStatus()
    } else {
      setIsLoading(false)
    }
  }, [])

  const checkAuthStatus = async () => {
    try {
      const response = await api.get('/api/v1/user/profile')
      setUser(response.data.data)
      setIsLoading(false)
    } catch (error) {
      // Token is invalid, clear it
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      delete api.defaults.headers.common['Authorization']
      setUser(null)
      setIsLoading(false)
    }
  }

  const login = async (email: string, password: string) => {
    try {
      const response = await api.post<{ data: AuthResponse }>('/api/v1/auth/login', {
        email,
        password,
      })

      const { access_token, refresh_token, user: userData } = response.data.data

      // Store tokens
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)

      // Set auth header
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`

      // Set user
      setUser(userData)

      // Clear any existing queries
      queryClient.clear()
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Login failed')
    }
  }

  const register = async (email: string, username: string, password: string, firstName?: string, lastName?: string) => {
    try {
      const response = await api.post<{ data: AuthResponse }>('/api/v1/auth/register', {
        email,
        username,
        password,
        first_name: firstName,
        last_name: lastName,
      })

      const { access_token, refresh_token, user: userData } = response.data.data

      // Store tokens
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)

      // Set auth header
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`

      // Set user
      setUser(userData)

      // Clear any existing queries
      queryClient.clear()
    } catch (error: any) {
      throw new Error(error.response?.data?.message || 'Registration failed')
    }
  }

  const logout = () => {
    // Clear tokens
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')

    // Clear auth header
    delete api.defaults.headers.common['Authorization']

    // Clear user
    setUser(null)

    // Clear all queries
    queryClient.clear()
  }

  const refreshToken = async () => {
    try {
      const refresh_token = localStorage.getItem('refresh_token')
      if (!refresh_token) {
        throw new Error('No refresh token available')
      }

      const response = await api.post<{ data: AuthResponse }>('/api/v1/auth/refresh', {
        refresh_token,
      })

      const { access_token, refresh_token: newRefreshToken, user: userData } = response.data.data

      // Store new tokens
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', newRefreshToken)

      // Set auth header
      api.defaults.headers.common['Authorization'] = `Bearer ${access_token}`

      // Set user
      setUser(userData)
    } catch (error: any) {
      // Refresh failed, logout user
      logout()
      throw new Error('Token refresh failed')
    }
  }

  // Set up axios interceptor for automatic token refresh
  useEffect(() => {
    const interceptor = api.interceptors.response.use(
      (response) => response,
      async (error) => {
        const originalRequest = error.config

        if (error.response?.status === 401 && !originalRequest._retry) {
          originalRequest._retry = true

          try {
            await refreshToken()
            // Retry the original request
            return api(originalRequest)
          } catch (refreshError) {
            // Refresh failed, redirect to login
            logout()
            return Promise.reject(refreshError)
          }
        }

        return Promise.reject(error)
      }
    )

    return () => {
      api.interceptors.response.eject(interceptor)
    }
  }, [])

  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated: !!user,
    login,
    register,
    logout,
    refreshToken,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}