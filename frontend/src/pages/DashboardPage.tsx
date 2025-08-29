import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '@hooks/useAuth'
import { useAI } from '@hooks/useAI'
import { Button } from '@components/ui/Button'
import { Card } from '@components/ui/Card'
import { Badge } from '@components/ui/Badge'
import { Input } from '@components/ui/Input'
import { 
  Plus, 
  Sparkles, 
  FileText, 
  Code, 
  Image, 
  Users, 
  Clock, 
  Star,
  TrendingUp,
  Lightbulb,
  Zap,
  Brain,
  Search,
  Filter,
  Grid,
  List,
  Menu,
  X,
  Smartphone
} from 'lucide-react'

interface ContentItem {
  id: string
  title: string
  type: 'document' | 'code' | 'diagram' | 'template'
  description: string
  lastModified: string
  collaborators: number
  isPublic: boolean
  aiGenerated: boolean
}

export const DashboardPage: React.FC = () => {
  const navigate = useNavigate()
  const { user } = useAuth()
  const { isAvailable, modelStatus } = useAI()
  
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedType, setSelectedType] = useState<string>('all')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [recentContent, setRecentContent] = useState<ContentItem[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [activeMobileTab, setActiveMobileTab] = useState<'content' | 'ai' | 'search'>('content')

  // Sample content data (in real app, this would come from API)
  useEffect(() => {
    const loadContent = async () => {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      setRecentContent([
        {
          id: '1',
          title: 'AI-Powered Business Proposal',
          type: 'document',
          description: 'Professional business proposal generated with AI assistance',
          lastModified: '2 hours ago',
          collaborators: 3,
          isPublic: false,
          aiGenerated: true
        },
        {
          id: '2',
          title: 'React Component Library',
          type: 'code',
          description: 'Reusable React components with TypeScript',
          lastModified: '1 day ago',
          collaborators: 2,
          isPublic: true,
          aiGenerated: false
        },
        {
          id: '3',
          title: 'System Architecture Diagram',
          type: 'diagram',
          description: 'Microservices architecture visualization',
          lastModified: '3 days ago',
          collaborators: 5,
          isPublic: false,
          aiGenerated: true
        },
        {
          id: '4',
          title: 'Project Management Template',
          type: 'template',
          description: 'AI-generated project planning template',
          lastModified: '1 week ago',
          collaborators: 1,
          isPublic: true,
          aiGenerated: true
        }
      ])
      setIsLoading(false)
    }

    loadContent()
  }, [])

  const handleCreateNew = (type: string) => {
    navigate(`/editor?type=${type}`)
  }

  const handleContentClick = (content: ContentItem) => {
    navigate(`/editor/${content.id}`)
  }

  const filteredContent = recentContent.filter(content => {
    const matchesSearch = content.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
                         content.description.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesType = selectedType === 'all' || content.type === selectedType
    return matchesSearch && matchesType
  })

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'document': return <FileText className="w-4 h-4 sm:w-5 sm:h-5" />
      case 'code': return <Code className="w-4 h-4 sm:w-5 sm:h-5" />
      case 'diagram': return <Image className="w-4 h-4 sm:w-5 sm:h-5" />
      case 'template': return <FileText className="w-4 h-4 sm:w-5 sm:h-5" />
      default: return <FileText className="w-4 h-4 sm:w-5 sm:h-5" />
    }
  }

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'document': return 'bg-blue-100 text-blue-800'
      case 'code': return 'bg-green-100 text-green-800'
      case 'diagram': return 'bg-purple-100 text-purple-800'
      case 'template': return 'bg-orange-100 text-orange-800'
      default: return 'bg-gray-100 text-gray-800'
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 sm:h-12 sm:w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-sm sm:text-base text-gray-600">Loading your AI-powered workspace...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Mobile Header */}
      <div className="sm:hidden bg-white border-b border-gray-200">
        <div className="flex items-center justify-between px-4 py-3">
          <div className="flex items-center gap-2">
            <Smartphone className="w-5 h-5 text-blue-600" />
            <h1 className="text-lg font-bold text-gray-900">Open-Same</h1>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
          </Button>
        </div>

        {/* Mobile Tab Navigation */}
        <div className="flex bg-white border-b border-gray-200">
          <button
            onClick={() => setActiveMobileTab('content')}
            className={`flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeMobileTab === 'content'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Content
          </button>
          <button
            onClick={() => setActiveMobileTab('ai')}
            className={`flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeMobileTab === 'ai'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            AI Tools
          </button>
          <button
            onClick={() => setActiveMobileTab('search')}
            className={`flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeMobileTab === 'search'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Search
          </button>
        </div>
      </div>

      {/* Desktop Header */}
      <div className="hidden sm:block bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">
                Welcome back, {user?.firstName || user?.username}!
              </h1>
              <p className="text-gray-600 mt-1">
                Your AI-powered content creation workspace
              </p>
            </div>
            
            <div className="flex items-center gap-3">
              {isAvailable && (
                <Badge variant="secondary" className="flex items-center gap-2">
                  <Brain className="w-4 h-4 text-green-500" />
                  AI Available
                </Badge>
              )}
              <Button onClick={() => navigate('/profile')}>
                Profile
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-8">
        {/* AI Quick Actions - Mobile */}
        {isAvailable && activeMobileTab === 'ai' && (
          <div className="sm:hidden mb-6">
            <Card className="p-4 bg-gradient-to-r from-purple-50 to-blue-50 border-purple-200">
              <div className="flex items-center gap-3 mb-4">
                <Sparkles className="w-5 h-5 text-purple-600" />
                <h2 className="text-lg font-semibold text-gray-900">AI Quick Actions</h2>
              </div>
              
              <div className="grid grid-cols-2 gap-3">
                <Button
                  variant="outline"
                  className="h-auto p-3 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('document')}
                >
                  <FileText className="w-6 h-6 text-blue-600" />
                  <span className="text-sm font-medium">AI Document</span>
                </Button>
                
                <Button
                  variant="outline"
                  className="h-auto p-3 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('code')}
                >
                  <Code className="w-6 h-6 text-green-600" />
                  <span className="text-sm font-medium">AI Code</span>
                </Button>
                
                <Button
                  variant="outline"
                  className="h-auto p-3 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('diagram')}
                >
                  <Image className="w-6 h-6 text-purple-600" />
                  <span className="text-sm font-medium">AI Diagram</span>
                </Button>
                
                <Button
                  variant="outline"
                  className="h-auto p-3 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('template')}
                >
                  <FileText className="w-6 h-6 text-orange-600" />
                  <span className="text-sm font-medium">AI Template</span>
                </Button>
              </div>
            </Card>
          </div>
        )}

        {/* AI Quick Actions - Desktop */}
        {isAvailable && (
          <div className="hidden sm:block mb-8">
            <Card className="p-6 bg-gradient-to-r from-purple-50 to-blue-50 border-purple-200">
              <div className="flex items-center gap-3 mb-4">
                <Sparkles className="w-6 h-6 text-purple-600" />
                <h2 className="text-xl font-semibold text-gray-900">AI Quick Actions</h2>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <Button
                  variant="outline"
                  className="h-auto p-4 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('document')}
                >
                  <FileText className="w-8 h-8 text-blue-600" />
                  <span className="font-medium">AI Document</span>
                  <span className="text-sm text-gray-500">Generate professional content</span>
                </Button>
                
                <Button
                  variant="outline"
                  className="h-auto p-4 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('code')}
                >
                  <Code className="w-8 h-8 text-green-600" />
                  <span className="font-medium">AI Code</span>
                  <span className="text-sm text-gray-500">Generate clean code</span>
                </Button>
                
                <Button
                  variant="outline"
                  className="h-auto p-4 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('diagram')}
                >
                  <Image className="w-8 h-8 text-purple-600" />
                  <span className="font-medium">AI Diagram</span>
                  <span className="text-sm text-gray-500">Create visual diagrams</span>
                </Button>
                
                <Button
                  variant="outline"
                  className="h-auto p-4 flex flex-col items-center gap-2 bg-white hover:bg-purple-50"
                  onClick={() => handleCreateNew('template')}
                >
                  <FileText className="w-8 h-8 text-orange-600" />
                  <span className="font-medium">AI Template</span>
                  <span className="text-sm text-gray-500">Generate templates</span>
                </Button>
              </div>
            </Card>
          </div>
        )}

        {/* Search and Filters - Mobile */}
        {activeMobileTab === 'search' && (
          <div className="sm:hidden mb-6">
            <div className="space-y-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
                <Input
                  placeholder="Search your content..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
              
              <div className="flex gap-2">
                <select
                  value={selectedType}
                  onChange={(e) => setSelectedType(e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                >
                  <option value="all">All Types</option>
                  <option value="document">Documents</option>
                  <option value="code">Code</option>
                  <option value="diagram">Diagrams</option>
                  <option value="template">Templates</option>
                </select>
                
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
                >
                  {viewMode === 'grid' ? <List className="w-4 h-4" /> : <Grid className="w-4 h-4" />}
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* Search and Filters - Desktop */}
        <div className="hidden sm:flex flex-col sm:flex-row gap-4 mb-6">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-5 h-5" />
            <Input
              placeholder="Search your content..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-10"
            />
          </div>
          
          <div className="flex gap-2">
            <select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Types</option>
              <option value="document">Documents</option>
              <option value="code">Code</option>
              <option value="diagram">Diagrams</option>
              <option value="template">Templates</option>
            </select>
            
            <Button
              variant="outline"
              size="sm"
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
            >
              {viewMode === 'grid' ? <List className="w-4 h-4" /> : <Grid className="w-4 h-4" />}
            </Button>
          </div>
        </div>

        {/* Content Grid/List */}
        <div className={`${activeMobileTab === 'content' ? 'block' : 'hidden'} sm:block space-y-6`}>
          <div className="flex items-center justify-between">
            <h2 className="text-lg sm:text-xl font-semibold text-gray-900">Recent Content</h2>
            <Button 
              onClick={() => handleCreateNew('document')} 
              className="flex items-center gap-2 text-sm sm:text-base"
              size="sm"
            >
              <Plus className="w-4 h-4" />
              <span className="hidden sm:inline">Create New</span>
              <span className="sm:hidden">New</span>
            </Button>
          </div>

          {filteredContent.length === 0 ? (
            <Card className="p-8 sm:p-12 text-center">
              <FileText className="w-12 h-12 sm:w-16 sm:h-16 text-gray-400 mx-auto mb-4" />
              <h3 className="text-base sm:text-lg font-medium text-gray-900 mb-2">No content found</h3>
              <p className="text-sm sm:text-base text-gray-600 mb-6">
                {searchQuery || selectedType !== 'all' 
                  ? 'Try adjusting your search or filters'
                  : 'Get started by creating your first piece of content'
                }
              </p>
              <Button onClick={() => handleCreateNew('document')} size="sm">
                Create Your First Content
              </Button>
            </Card>
          ) : (
            <div className={viewMode === 'grid' ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6' : 'space-y-3 sm:space-y-4'}>
              {filteredContent.map((content) => (
                <Card
                  key={content.id}
                  className={`p-4 sm:p-6 cursor-pointer hover:shadow-lg transition-shadow ${
                    viewMode === 'list' ? 'flex items-center gap-4' : ''
                  }`}
                  onClick={() => handleContentClick(content)}
                >
                  <div className={`flex items-center gap-3 mb-3 ${viewMode === 'list' ? 'mb-0' : ''}`}>
                    <div className={`p-2 rounded-lg ${getTypeColor(content.type)}`}>
                      {getTypeIcon(content.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1 sm:mb-2">
                        <h3 className="font-semibold text-gray-900 text-sm sm:text-base truncate">{content.title}</h3>
                        {content.aiGenerated && (
                          <Badge variant="secondary" className="text-xs flex-shrink-0">
                            <Sparkles className="w-3 h-3 mr-1" />
                            AI
                          </Badge>
                        )}
                      </div>
                      <p className="text-xs sm:text-sm text-gray-600 line-clamp-2">{content.description}</p>
                    </div>
                  </div>
                  
                  <div className={`flex items-center justify-between text-xs sm:text-sm text-gray-500 ${
                    viewMode === 'list' ? 'ml-12' : ''
                  }`}>
                    <div className="flex items-center gap-2 sm:gap-4">
                      <span className="flex items-center gap-1">
                        <Clock className="w-3 h-3 sm:w-4 sm:h-4" />
                        <span className="hidden sm:inline">{content.lastModified}</span>
                        <span className="sm:hidden">{content.lastModified.split(' ')[0]}</span>
                      </span>
                      <span className="flex items-center gap-1">
                        <Users className="w-3 h-3 sm:w-4 sm:h-4" />
                        {content.collaborators}
                      </span>
                    </div>
                    
                    <div className="flex items-center gap-2">
                      {content.isPublic && (
                        <Badge variant="outline" className="text-xs">
                          Public
                        </Badge>
                      )}
                      <Badge variant="outline" className={`text-xs ${getTypeColor(content.type)}`}>
                        {content.type}
                      </Badge>
                    </div>
                  </div>
                </Card>
              ))}
            </div>
          )}
        </div>

        {/* AI Status - Mobile */}
        {isAvailable && activeMobileTab === 'ai' && (
          <div className="sm:hidden mt-6">
            <Card className="p-4 bg-gray-50">
              <div className="flex items-center gap-3 mb-4">
                <Brain className="w-5 h-5 text-blue-600" />
                <h2 className="text-lg font-semibold text-gray-900">AI Status</h2>
              </div>
              
              <div className="space-y-3">
                {Object.entries(modelStatus).map(([provider, status]) => (
                  <div key={provider} className="flex items-center gap-3 p-3 bg-white rounded-lg">
                    <div className={`w-3 h-3 rounded-full ${
                      status?.available ? 'bg-green-500' : 'bg-red-500'
                    }`} />
                    <div>
                      <p className="font-medium capitalize text-sm">{provider}</p>
                      <p className="text-xs text-gray-600">{status?.model || 'Not configured'}</p>
                    </div>
                  </div>
                ))}
              </div>
            </Card>
          </div>
        )}

        {/* AI Status - Desktop */}
        {isAvailable && (
          <div className="hidden sm:block mt-8">
            <Card className="p-6 bg-gray-50">
              <div className="flex items-center gap-3 mb-4">
                <Brain className="w-6 h-6 text-blue-600" />
                <h2 className="text-xl font-semibold text-gray-900">AI Model Status</h2>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {Object.entries(modelStatus).map(([provider, status]) => (
                  <div key={provider} className="flex items-center gap-3 p-3 bg-white rounded-lg">
                    <div className={`w-3 h-3 rounded-full ${
                      status?.available ? 'bg-green-500' : 'bg-red-500'
                    }`} />
                    <div>
                      <p className="font-medium capitalize">{provider}</p>
                      <p className="text-sm text-gray-600">{status?.model || 'Not configured'}</p>
                    </div>
                  </div>
                ))}
              </div>
            </Card>
          </div>
        )}
      </div>
    </div>
  )
}