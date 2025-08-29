import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Toaster } from 'react-hot-toast'
import { AuthProvider } from '@contexts/AuthContext'
import { ThemeProvider } from '@contexts/ThemeContext'
import { CollaborationProvider } from '@contexts/CollaborationContext'
import { DashboardPage } from '@pages/DashboardPage'
import { LoginPage } from '@pages/LoginPage'
import { RegisterPage } from '@pages/RegisterPage'
import { EditorPage } from '@pages/EditorPage'
import { ProfilePage } from '@pages/ProfilePage'
import { MobileNavigation } from '@components/navigation/MobileNavigation'
import { DesktopNavigation } from '@components/navigation/DesktopNavigation'
import { useAuth } from '@hooks/useAuth'
import { useMediaQuery } from '@hooks/useMediaQuery'

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
})

// Protected Route Component
const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading } = useAuth()
  
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 sm:h-12 sm:w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-sm sm:text-base text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }
  
  if (!user) {
    return <Navigate to="/login" replace />
  }
  
  return <>{children}</>
}

// Public Route Component
const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading } = useAuth()
  
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 sm:h-12 sm:w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-sm sm:text-base text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }
  
  if (user) {
    return <Navigate to="/dashboard" replace />
  }
  
  return <>{children}</>
}

// Main App Layout
const AppLayout: React.FC = () => {
  const { user } = useAuth()
  const isMobile = useMediaQuery('(max-width: 768px)')
  
  if (!user) return null
  
  return (
    <div className="min-h-screen bg-gray-50">
      {isMobile ? <MobileNavigation /> : <DesktopNavigation />}
      <main className="pt-0">
        <Routes>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/editor" element={<EditorPage />} />
          <Route path="/editor/:id" element={<EditorPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </main>
    </div>
  )
}

// Main App Component
const App: React.FC = () => {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <AuthProvider>
          <CollaborationProvider>
            <Router>
              <div className="App">
                <Routes>
                  {/* Public Routes */}
                  <Route path="/login" element={
                    <PublicRoute>
                      <LoginPage />
                    </PublicRoute>
                  } />
                  <Route path="/register" element={
                    <PublicRoute>
                      <RegisterPage />
                    </PublicRoute>
                  } />
                  
                  {/* Protected Routes */}
                  <Route path="/*" element={
                    <ProtectedRoute>
                      <AppLayout />
                    </ProtectedRoute>
                  } />
                </Routes>
                
                {/* Global Toast Notifications */}
                <Toaster
                  position={window.innerWidth < 768 ? "top-center" : "top-right"}
                  toastOptions={{
                    duration: 4000,
                    style: {
                      background: '#363636',
                      color: '#fff',
                      fontSize: '14px',
                      padding: '12px 16px',
                      borderRadius: '8px',
                      maxWidth: '90vw',
                    },
                    success: {
                      iconTheme: {
                        primary: '#10B981',
                        secondary: '#fff',
                      },
                    },
                    error: {
                      iconTheme: {
                        primary: '#EF4444',
                        secondary: '#fff',
                      },
                    },
                  }}
                />
              </div>
            </Router>
          </CollaborationProvider>
        </AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  )
}

export default App