import React, { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '@hooks/useAuth'
import { Button } from '@components/ui/Button'
import { 
  Home, 
  Edit3, 
  Plus, 
  User, 
  Bell,
  Search,
  Settings,
  LogOut,
  ChevronDown,
  Sparkles
} from 'lucide-react'

export const DesktopNavigation: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false)
  const [isCreateMenuOpen, setIsCreateMenuOpen] = useState(false)

  const navigationItems = [
    {
      name: 'Dashboard',
      path: '/dashboard',
      icon: Home,
      active: location.pathname === '/dashboard'
    },
    {
      name: 'Editor',
      path: '/editor',
      icon: Edit3,
      active: location.pathname.startsWith('/editor')
    }
  ]

  const createOptions = [
    {
      name: 'Document',
      description: 'AI-powered content creation',
      icon: '📄',
      color: 'bg-blue-100 text-blue-600',
      path: '/editor?type=document'
    },
    {
      name: 'Code',
      description: 'AI-generated code',
      icon: '💻',
      color: 'bg-green-100 text-green-600',
      path: '/editor?type=code'
    },
    {
      name: 'Diagram',
      description: 'AI-created visuals',
      icon: '📊',
      color: 'bg-purple-100 text-purple-600',
      path: '/editor?type=diagram'
    },
    {
      name: 'Template',
      description: 'AI-generated templates',
      icon: '🎯',
      color: 'bg-orange-100 text-orange-600',
      path: '/editor?type=template'
    }
  ]

  const handleLogout = async () => {
    try {
      await logout()
      navigate('/login')
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }

  const handleCreate = (option: any) => {
    navigate(option.path)
    setIsCreateMenuOpen(false)
  }

  return (
    <nav className="bg-white border-b border-gray-200 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo and Brand */}
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
                <Sparkles className="w-5 h-5 text-white" />
              </div>
              <span className="text-xl font-bold text-gray-900">Open-Same</span>
            </div>
            
            {/* Navigation Links */}
            <div className="hidden md:flex items-center gap-1 ml-8">
              {navigationItems.map((item) => (
                <button
                  key={item.name}
                  onClick={() => navigate(item.path)}
                  className={`flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors ${
                    item.active
                      ? 'bg-blue-50 text-blue-600'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                  }`}
                >
                  <item.icon className="w-4 h-4" />
                  {item.name}
                </button>
              ))}
            </div>
          </div>

          {/* Right Side Actions */}
          <div className="flex items-center gap-3">
            {/* Search */}
            <Button variant="ghost" size="sm" className="p-2">
              <Search className="w-5 h-5 text-gray-600" />
            </Button>

            {/* Notifications */}
            <Button variant="ghost" size="sm" className="p-2 relative">
              <Bell className="w-5 h-5 text-gray-600" />
              <div className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full"></div>
            </Button>

            {/* Create Button with Dropdown */}
            <div className="relative">
              <Button
                onClick={() => setIsCreateMenuOpen(!isCreateMenuOpen)}
                className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white"
              >
                <Plus className="w-4 h-4" />
                Create
                <ChevronDown className="w-4 h-4" />
              </Button>

              {/* Create Dropdown */}
              {isCreateMenuOpen && (
                <div className="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50">
                  <div className="px-4 py-2 border-b border-gray-200">
                    <h3 className="text-sm font-medium text-gray-900">Create New</h3>
                    <p className="text-xs text-gray-600">Choose your content type</p>
                  </div>
                  
                  <div className="p-2">
                    {createOptions.map((option) => (
                      <button
                        key={option.name}
                        onClick={() => handleCreate(option)}
                        className="w-full flex items-center gap-3 p-3 rounded-lg hover:bg-gray-50 transition-colors text-left"
                      >
                        <div className={`w-10 h-10 rounded-lg flex items-center justify-center text-lg ${option.color}`}>
                          {option.icon}
                        </div>
                        <div>
                          <div className="font-medium text-gray-900">{option.name}</div>
                          <div className="text-sm text-gray-600">{option.description}</div>
                        </div>
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* User Menu */}
            <div className="relative">
              <button
                onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div className="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
                  <span className="text-white font-semibold text-sm">
                    {user?.firstName?.charAt(0) || user?.username?.charAt(0) || 'U'}
                  </span>
                </div>
                <div className="hidden md:block text-left">
                  <div className="text-sm font-medium text-gray-900">
                    {user?.firstName || user?.username}
                  </div>
                  <div className="text-xs text-gray-600">{user?.email}</div>
                </div>
                <ChevronDown className="w-4 h-4 text-gray-400" />
              </button>

              {/* User Dropdown */}
              {isUserMenuOpen && (
                <div className="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50">
                  <div className="px-4 py-2 border-b border-gray-200">
                    <div className="text-sm font-medium text-gray-900">
                      {user?.firstName || user?.username}
                    </div>
                    <div className="text-xs text-gray-600">{user?.email}</div>
                  </div>
                  
                  <div className="py-1">
                    <button
                      onClick={() => {
                        navigate('/profile')
                        setIsUserMenuOpen(false)
                      }}
                      className="w-full flex items-center gap-3 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                    >
                      <User className="w-4 h-4" />
                      Profile
                    </button>
                    
                    <button
                      onClick={() => {
                        navigate('/settings')
                        setIsUserMenuOpen(false)
                      }}
                      className="w-full flex items-center gap-3 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
                    >
                      <Settings className="w-4 h-4" />
                      Settings
                    </button>
                  </div>
                  
                  <div className="border-t border-gray-200 py-1">
                    <button
                      onClick={handleLogout}
                      className="w-full flex items-center gap-3 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                    >
                      <LogOut className="w-4 h-4" />
                      Sign Out
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Backdrop for dropdowns */}
      {(isUserMenuOpen || isCreateMenuOpen) && (
        <div 
          className="fixed inset-0 z-30"
          onClick={() => {
            setIsUserMenuOpen(false)
            setIsCreateMenuOpen(false)
          }}
        />
      )}
    </nav>
  )
}