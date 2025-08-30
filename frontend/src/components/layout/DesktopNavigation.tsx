import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import { 
  Home, 
  FileText, 
  Plus, 
  Users, 
  Settings, 
  LogOut,
  Search,
  Bell,
  User
} from 'lucide-react'
import { useAuth } from '@hooks/useAuth'
import { Button } from '@components/ui/Button'
import { Avatar } from '@components/ui/Avatar'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@components/ui/DropdownMenu'

export const DesktopNavigation: React.FC = () => {
  const { user, logout } = useAuth()
  const location = useLocation()

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', icon: Home },
    { name: 'My Content', href: '/content', icon: FileText },
    { name: 'Create', href: '/create', icon: Plus },
    { name: 'Collaborations', href: '/collaborations', icon: Users },
  ]

  const isActive = (path: string) => {
    return location.pathname === path
  }

  const handleLogout = () => {
    logout()
  }

  return (
    <div className="hidden md:flex md:w-64 md:flex-col md:fixed md:inset-y-0">
      <div className="flex-1 flex flex-col min-h-0 bg-white border-r border-gray-200">
        {/* Logo */}
        <div className="flex items-center h-16 flex-shrink-0 px-4 bg-gradient-to-r from-blue-600 to-purple-600">
          <Link to="/dashboard" className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-8 h-8 bg-white rounded-lg flex items-center justify-center">
                <span className="text-blue-600 font-bold text-lg">S</span>
              </div>
            </div>
            <div className="ml-3">
              <h1 className="text-white text-xl font-bold">Open-Same</h1>
            </div>
          </Link>
        </div>

        {/* Navigation */}
        <div className="flex-1 flex flex-col overflow-y-auto">
          <nav className="flex-1 px-2 py-4 space-y-1">
            {navigation.map((item) => {
              const Icon = item.icon
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`group flex items-center px-2 py-2 text-sm font-medium rounded-md transition-colors ${
                    isActive(item.href)
                      ? 'bg-blue-100 text-blue-700 border-r-2 border-blue-600'
                      : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                  }`}
                >
                  <Icon
                    className={`mr-3 flex-shrink-0 h-5 w-5 ${
                      isActive(item.href) ? 'text-blue-500' : 'text-gray-400 group-hover:text-gray-500'
                    }`}
                  />
                  {item.name}
                </Link>
              )
            })}
          </nav>

          {/* User Section */}
          <div className="flex-shrink-0 border-t border-gray-200 p-4">
            <div className="flex items-center">
              <Avatar
                src={user?.avatar}
                alt={user?.username || 'User'}
                fallback={user?.username?.charAt(0)?.toUpperCase() || 'U'}
                className="w-8 h-8"
              />
              <div className="ml-3 min-w-0 flex-1">
                <p className="text-sm font-medium text-gray-900 truncate">
                  {user?.fullName || user?.username}
                </p>
                <p className="text-sm text-gray-500 truncate">
                  {user?.email}
                </p>
              </div>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="sm" className="ml-2">
                    <Settings className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48">
                  <DropdownMenuItem asChild>
                    <Link to="/profile" className="flex items-center">
                      <User className="mr-2 h-4 w-4" />
                      Profile
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem asChild>
                    <Link to="/settings" className="flex items-center">
                      <Settings className="mr-2 h-4 w-4" />
                      Settings
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={handleLogout} className="text-red-600">
                    <LogOut className="mr-2 h-4 w-4" />
                    Logout
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}