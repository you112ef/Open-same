import React, { useState } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { 
  Home, 
  FileText, 
  Plus, 
  Users, 
  Settings, 
  LogOut,
  Menu,
  X,
  User
} from 'lucide-react'
import { useAuth } from '@hooks/useAuth'
import { Button } from '@components/ui/Button'
import { Avatar } from '@components/ui/Avatar'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@components/ui/Sheet'

export const MobileNavigation: React.FC = () => {
  const { user, logout } = useAuth()
  const location = useLocation()
  const [isOpen, setIsOpen] = useState(false)

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
    setIsOpen(false)
  }

  const handleNavigation = () => {
    setIsOpen(false)
  }

  return (
    <div className="md:hidden">
      {/* Top Bar */}
      <div className="fixed top-0 left-0 right-0 z-50 bg-white border-b border-gray-200 px-4 py-3">
        <div className="flex items-center justify-between">
          {/* Logo */}
          <Link to="/dashboard" className="flex items-center">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-lg">S</span>
            </div>
            <span className="ml-2 text-lg font-bold text-gray-900">Open-Same</span>
          </Link>

          {/* Menu Button */}
          <Sheet open={isOpen} onOpenChange={setIsOpen}>
            <SheetTrigger asChild>
              <Button variant="ghost" size="sm">
                <Menu className="h-6 w-6" />
              </Button>
            </SheetTrigger>
            <SheetContent side="right" className="w-80">
              <SheetHeader>
                <SheetTitle>Menu</SheetTitle>
              </SheetHeader>
              
              {/* User Info */}
              <div className="py-6 border-b border-gray-200">
                <div className="flex items-center">
                  <Avatar
                    src={user?.avatar}
                    alt={user?.username || 'User'}
                    fallback={user?.username?.charAt(0)?.toUpperCase() || 'U'}
                    className="w-12 h-12"
                  />
                  <div className="ml-4">
                    <p className="text-sm font-medium text-gray-900">
                      {user?.fullName || user?.username}
                    </p>
                    <p className="text-sm text-gray-500">
                      {user?.email}
                    </p>
                  </div>
                </div>
              </div>

              {/* Navigation */}
              <nav className="py-4 space-y-2">
                {navigation.map((item) => {
                  const Icon = item.icon
                  return (
                    <Link
                      key={item.name}
                      to={item.href}
                      onClick={handleNavigation}
                      className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                        isActive(item.href)
                          ? 'bg-blue-100 text-blue-700'
                          : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                      }`}
                    >
                      <Icon
                        className={`mr-3 flex-shrink-0 h-5 w-5 ${
                          isActive(item.href) ? 'text-blue-500' : 'text-gray-400'
                        }`}
                      />
                      {item.name}
                    </Link>
                  )
                })}
              </nav>

              {/* User Actions */}
              <div className="py-4 border-t border-gray-200 space-y-2">
                <Link
                  to="/profile"
                  onClick={handleNavigation}
                  className="flex items-center px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-md transition-colors"
                >
                  <User className="mr-3 h-5 w-5 text-gray-400" />
                  Profile
                </Link>
                <Link
                  to="/settings"
                  onClick={handleNavigation}
                  className="flex items-center px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-md transition-colors"
                >
                  <Settings className="mr-3 h-5 w-5 text-gray-400" />
                  Settings
                </Link>
                <button
                  onClick={handleLogout}
                  className="flex items-center w-full px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 hover:text-red-700 rounded-md transition-colors"
                >
                  <LogOut className="mr-3 h-5 w-5" />
                  Logout
                </button>
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>

      {/* Bottom Navigation */}
      <div className="fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 px-2 py-2">
        <div className="flex justify-around">
          {navigation.map((item) => {
            const Icon = item.icon
            return (
              <Link
                key={item.name}
                to={item.href}
                className={`flex flex-col items-center py-2 px-3 text-xs font-medium rounded-md transition-colors ${
                  isActive(item.href)
                    ? 'text-blue-600 bg-blue-50'
                    : 'text-gray-600 hover:text-blue-600'
                }`}
              >
                <Icon
                  className={`h-5 w-5 mb-1 ${
                    isActive(item.href) ? 'text-blue-600' : 'text-gray-400'
                  }`}
                />
                <span>{item.name}</span>
              </Link>
            )
          })}
        </div>
      </div>

      {/* Spacer for top bar */}
      <div className="h-16" />
      {/* Spacer for bottom navigation */}
      <div className="h-20" />
    </div>
  )
}