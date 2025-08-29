import React, { useState } from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import { useAuth } from '@hooks/useAuth'
import { Button } from '@components/ui/Button'
import { 
  Home, 
  Edit3, 
  Plus, 
  User, 
  Menu,
  X,
  Settings,
  LogOut,
  Bell,
  Search
} from 'lucide-react'

export const MobileNavigation: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()
  const [isMenuOpen, setIsMenuOpen] = useState(false)

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
    },
    {
      name: 'Create',
      path: '/editor',
      icon: Plus,
      active: false,
      onClick: () => navigate('/editor')
    },
    {
      name: 'Profile',
      path: '/profile',
      icon: User,
      active: location.pathname === '/profile'
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

  const handleNavigation = (item: any) => {
    if (item.onClick) {
      item.onClick()
    } else {
      navigate(item.path)
    }
    setIsMenuOpen(false)
  }

  return (
    <>
      {/* Bottom Tab Navigation */}
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-50 sm:hidden">
        <div className="flex items-center justify-around py-2">
          {navigationItems.map((item) => (
            <button
              key={item.name}
              onClick={() => handleNavigation(item)}
              className={`flex flex-col items-center gap-1 py-2 px-3 rounded-lg transition-colors ${
                item.active
                  ? 'text-blue-600 bg-blue-50'
                  : 'text-gray-600 hover:text-blue-600 hover:bg-blue-50'
              }`}
            >
              <item.icon className={`w-5 h-5 ${item.active ? 'text-blue-600' : 'text-gray-600'}`} />
              <span className="text-xs font-medium">{item.name}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Top Header with Menu */}
      <div className="fixed top-0 left-0 right-0 bg-white border-b border-gray-200 z-40 sm:hidden">
        <div className="flex items-center justify-between px-4 py-3">
          <div className="flex items-center gap-3">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="p-2"
            >
              {isMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </Button>
            <div>
              <h1 className="text-lg font-bold text-gray-900">Open-Same</h1>
              <p className="text-xs text-gray-600">AI-Powered Workspace</p>
            </div>
          </div>
          
          <div className="flex items-center gap-2">
            <Button variant="ghost" size="sm" className="p-2">
              <Bell className="w-5 h-5 text-gray-600" />
            </Button>
            <Button variant="ghost" size="sm" className="p-2">
              <Search className="w-5 h-5 text-gray-600" />
            </Button>
          </div>
        </div>
      </div>

      {/* Side Menu Overlay */}
      {isMenuOpen && (
        <div className="fixed inset-0 z-50 sm:hidden">
          {/* Backdrop */}
          <div 
            className="absolute inset-0 bg-black bg-opacity-50"
            onClick={() => setIsMenuOpen(false)}
          />
          
          {/* Menu Panel */}
          <div className="absolute left-0 top-0 bottom-0 w-80 bg-white shadow-xl">
            <div className="p-6">
              {/* User Info */}
              <div className="mb-6 pb-6 border-b border-gray-200">
                <div className="flex items-center gap-3 mb-3">
                  <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center">
                    <span className="text-white font-semibold text-lg">
                      {user?.firstName?.charAt(0) || user?.username?.charAt(0) || 'U'}
                    </span>
                  </div>
                  <div>
                    <h3 className="font-semibold text-gray-900">
                      {user?.firstName || user?.username}
                    </h3>
                    <p className="text-sm text-gray-600">{user?.email}</p>
                  </div>
                </div>
                
                <div className="flex items-center gap-2">
                  <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                  <span className="text-sm text-gray-600">Online</span>
                </div>
              </div>

              {/* Menu Items */}
              <nav className="space-y-2">
                <button
                  onClick={() => handleNavigation({ path: '/dashboard' })}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <Home className="w-5 h-5 text-gray-600" />
                  <span className="font-medium">Dashboard</span>
                </button>
                
                <button
                  onClick={() => handleNavigation({ path: '/editor' })}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <Edit3 className="w-5 h-5 text-gray-600" />
                  <span className="font-medium">Editor</span>
                </button>
                
                <button
                  onClick={() => handleNavigation({ path: '/profile' })}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <User className="w-5 h-5 text-gray-600" />
                  <span className="font-medium">Profile</span>
                </button>
                
                <button
                  onClick={() => navigate('/settings')}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <Settings className="w-5 h-5 text-gray-600" />
                  <span className="font-medium">Settings</span>
                </button>
              </nav>

              {/* Divider */}
              <div className="my-6 border-t border-gray-200" />

              {/* Quick Actions */}
              <div className="space-y-2">
                <h4 className="text-sm font-medium text-gray-900 mb-3">Quick Actions</h4>
                
                <button
                  onClick={() => handleNavigation({ path: '/editor', query: 'type=document' })}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div className="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                    <span className="text-blue-600 font-semibold text-sm">D</span>
                  </div>
                  <div>
                    <span className="font-medium">New Document</span>
                    <p className="text-xs text-gray-600">AI-powered content</p>
                  </div>
                </button>
                
                <button
                  onClick={() => handleNavigation({ path: '/editor', query: 'type=code' })}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div className="w-8 h-8 bg-green-100 rounded-lg flex items-center justify-center">
                    <span className="text-green-600 font-semibold text-sm">C</span>
                  </div>
                  <div>
                    <span className="font-medium">New Code</span>
                    <p className="text-xs text-gray-600">AI-generated code</p>
                  </div>
                </button>
                
                <button
                  onClick={() => handleNavigation({ path: '/editor', query: 'type=diagram' })}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-gray-100 transition-colors"
                >
                  <div className="w-8 h-8 bg-purple-100 rounded-lg flex items-center justify-center">
                    <span className="text-purple-600 font-semibold text-sm">V</span>
                  </div>
                  <div>
                    <span className="font-medium">New Diagram</span>
                    <p className="text-xs text-gray-600">AI-created visuals</p>
                  </div>
                </button>
              </div>

              {/* Logout */}
              <div className="mt-8 pt-6 border-t border-gray-200">
                <button
                  onClick={handleLogout}
                  className="w-full flex items-center gap-3 px-3 py-2 text-left rounded-lg hover:bg-red-50 transition-colors text-red-600"
                >
                  <LogOut className="w-5 h-5" />
                  <span className="font-medium">Sign Out</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Bottom Spacing for Content */}
      <div className="pb-20 sm:hidden" />
      {/* Top Spacing for Content */}
      <div className="pt-20 sm:hidden" />
    </>
  )
}