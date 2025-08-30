import React, { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { 
  Plus, 
  FileText, 
  Users, 
  TrendingUp, 
  Clock, 
  Star,
  Search,
  Filter,
  Grid,
  List
} from 'lucide-react'
import { useAuth } from '@contexts/AuthContext'
import { useQuery } from '@tanstack/react-query'
import { api, endpoints } from '@services/api'
import { Content, ContentType } from '@types/models'
import { Button } from '@components/ui/Button'
import { Input } from '@components/ui/Input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@components/ui/Card'
import { Badge } from '@components/ui/Badge'
import { Avatar } from '@components/ui/Avatar'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@components/ui/Tabs'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@components/ui/Select'

export const DashboardPage: React.FC = () => {
  const { user } = useAuth()
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedType, setSelectedType] = useState<ContentType | 'all'>('all')
  const [selectedStatus, setSelectedStatus] = useState<string>('all')

  // Fetch user's content
  const { data: contentData, isLoading: contentLoading } = useQuery({
    queryKey: ['user-content', searchQuery, selectedType, selectedStatus],
    queryFn: async () => {
      const params = new URLSearchParams()
      if (searchQuery) params.append('search', searchQuery)
      if (selectedType !== 'all') params.append('type', selectedType)
      if (selectedStatus !== 'all') params.append('status', selectedStatus)
      
      const response = await api.get(`${endpoints.content.list}?${params.toString()}`)
      return response.data.data
    },
    enabled: !!user,
  })

  // Fetch recent collaborations
  const { data: collaborations, isLoading: collaborationsLoading } = useQuery({
    queryKey: ['collaborations'],
    queryFn: async () => {
      const response = await api.get(endpoints.collaboration.list)
      return response.data.data
    },
    enabled: !!user,
  })

  // Fetch public content
  const { data: publicContent, isLoading: publicLoading } = useQuery({
    queryKey: ['public-content'],
    queryFn: async () => {
      const response = await api.get(endpoints.content.public)
      return response.data.data
    },
  })

  const getContentIcon = (type: ContentType) => {
    switch (type) {
      case 'code':
        return '💻'
      case 'diagram':
        return '📊'
      case 'image':
        return '🖼️'
      case 'document':
        return '📄'
      case 'template':
        return '📋'
      default:
        return '📝'
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'published':
        return 'bg-green-100 text-green-800'
      case 'draft':
        return 'bg-yellow-100 text-yellow-800'
      case 'archived':
        return 'bg-gray-100 text-gray-800'
      default:
        return 'bg-blue-100 text-blue-800'
    }
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    const now = new Date()
    const diffTime = Math.abs(now.getTime() - date.getTime())
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    
    if (diffDays === 1) return 'Today'
    if (diffDays === 2) return 'Yesterday'
    if (diffDays < 7) return `${diffDays - 1} days ago`
    return date.toLocaleDateString()
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white border-b border-gray-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
            <p className="text-gray-600">Welcome back, {user?.first_name || user?.username}!</p>
          </div>
          <Link to="/create">
            <Button className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700">
              <Plus className="mr-2 h-4 w-4" />
              Create Content
            </Button>
          </Link>
        </div>
      </div>

      <div className="px-6 py-6">
        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Content</CardTitle>
              <FileText className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{contentData?.total || 0}</div>
              <p className="text-xs text-muted-foreground">
                {contentData?.data?.length || 0} items created
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Collaborations</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{collaborations?.length || 0}</div>
              <p className="text-xs text-muted-foreground">
                Active collaborations
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Public Content</CardTitle>
              <TrendingUp className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {contentData?.data?.filter((c: Content) => c.is_public).length || 0}
              </div>
              <p className="text-xs text-muted-foreground">
                Published items
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Recent Activity</CardTitle>
              <Clock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {contentData?.data?.filter((c: Content) => {
                  const date = new Date(c.updated_at)
                  const now = new Date()
                  const diffDays = Math.ceil(Math.abs(now.getTime() - date.getTime()) / (1000 * 60 * 60 * 24))
                  return diffDays <= 7
                }).length || 0}
              </div>
              <p className="text-xs text-muted-foreground">
                Updated this week
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Main Content */}
        <Tabs defaultValue="my-content" className="space-y-6">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="my-content">My Content</TabsTrigger>
            <TabsTrigger value="collaborations">Collaborations</TabsTrigger>
            <TabsTrigger value="discover">Discover</TabsTrigger>
          </TabsList>

          {/* My Content Tab */}
          <TabsContent value="my-content" className="space-y-6">
            {/* Filters and Search */}
            <div className="flex flex-col sm:flex-row gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                <Input
                  placeholder="Search your content..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>
              
              <Select value={selectedType} onValueChange={(value) => setSelectedType(value as ContentType | 'all')}>
                <SelectTrigger className="w-full sm:w-40">
                  <SelectValue placeholder="Content Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Types</SelectItem>
                  <SelectItem value="text">Text</SelectItem>
                  <SelectItem value="code">Code</SelectItem>
                  <SelectItem value="diagram">Diagram</SelectItem>
                  <SelectItem value="image">Image</SelectItem>
                  <SelectItem value="document">Document</SelectItem>
                  <SelectItem value="template">Template</SelectItem>
                </SelectContent>
              </Select>

              <Select value={selectedStatus} onValueChange={setSelectedStatus}>
                <SelectTrigger className="w-full sm:w-40">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="draft">Draft</SelectItem>
                  <SelectItem value="published">Published</SelectItem>
                  <SelectItem value="archived">Archived</SelectItem>
                </SelectContent>
              </Select>

              <div className="flex border rounded-md">
                <Button
                  variant={viewMode === 'grid' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('grid')}
                  className="rounded-r-none"
                >
                  <Grid className="h-4 w-4" />
                </Button>
                <Button
                  variant={viewMode === 'list' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('list')}
                  className="rounded-l-none"
                >
                  <List className="h-4 w-4" />
                </Button>
              </div>
            </div>

            {/* Content Grid/List */}
            {contentLoading ? (
              <div className="flex items-center justify-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : contentData?.data?.length === 0 ? (
              <Card className="text-center py-12">
                <CardContent>
                  <FileText className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                  <h3 className="text-lg font-medium text-gray-900 mb-2">No content yet</h3>
                  <p className="text-gray-600 mb-4">Start creating your first piece of content</p>
                  <Link to="/create">
                    <Button>Create Content</Button>
                  </Link>
                </CardContent>
              </Card>
            ) : (
              <div className={viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'space-y-4'}>
                {contentData?.data?.map((content: Content) => (
                  <Card key={content.id} className="hover:shadow-md transition-shadow">
                    <CardHeader>
                      <div className="flex items-start justify-between">
                        <div className="flex items-center space-x-2">
                          <span className="text-2xl">{getContentIcon(content.type)}</span>
                          <div>
                            <CardTitle className="text-lg">{content.title}</CardTitle>
                            <CardDescription className="text-sm">
                              {content.description || 'No description'}
                            </CardDescription>
                          </div>
                        </div>
                        <Badge className={getStatusColor(content.status)}>
                          {content.status}
                        </Badge>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-3">
                        <div className="flex items-center justify-between text-sm text-gray-500">
                          <span>Type: {content.type}</span>
                          <span>v{content.version}</span>
                        </div>
                        
                        {content.tags && content.tags.length > 0 && (
                          <div className="flex flex-wrap gap-1">
                            {content.tags.slice(0, 3).map((tag, index) => (
                              <Badge key={index} variant="secondary" className="text-xs">
                                {tag}
                              </Badge>
                            ))}
                            {content.tags.length > 3 && (
                              <Badge variant="secondary" className="text-xs">
                                +{content.tags.length - 3}
                              </Badge>
                            )}
                          </div>
                        )}

                        <div className="flex items-center justify-between text-sm text-gray-500">
                          <span>Updated {formatDate(content.updated_at)}</span>
                          {content.is_public && (
                            <Badge variant="outline" className="text-xs">
                              Public
                            </Badge>
                          )}
                        </div>

                        <div className="flex space-x-2">
                          <Link to={`/editor/${content.id}`} className="flex-1">
                            <Button variant="outline" size="sm" className="w-full">
                              Edit
                            </Button>
                          </Link>
                          <Button variant="outline" size="sm">
                            <Users className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>

          {/* Collaborations Tab */}
          <TabsContent value="collaborations" className="space-y-6">
            {collaborationsLoading ? (
              <div className="flex items-center justify-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : collaborations?.length === 0 ? (
              <Card className="text-center py-12">
                <CardContent>
                  <Users className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                  <h3 className="text-lg font-medium text-gray-900 mb-2">No collaborations yet</h3>
                  <p className="text-gray-600">You haven't been invited to collaborate on any content yet.</p>
                </CardContent>
              </Card>
            ) : (
              <div className="space-y-4">
                {collaborations?.map((collab: any) => (
                  <Card key={collab.id}>
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <Avatar
                            src={collab.content?.user?.avatar}
                            fallback={collab.content?.user?.username?.charAt(0)?.toUpperCase() || 'U'}
                            className="w-10 h-10"
                          />
                          <div>
                            <h4 className="font-medium">{collab.content?.title}</h4>
                            <p className="text-sm text-gray-500">
                              by {collab.content?.user?.username} • {collab.role}
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center space-x-2">
                          <Badge variant="outline">{collab.content?.type}</Badge>
                          <Link to={`/editor/${collab.content?.id}`}>
                            <Button size="sm">Open</Button>
                          </Link>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>

          {/* Discover Tab */}
          <TabsContent value="discover" className="space-y-6">
            {publicLoading ? (
              <div className="flex items-center justify-center py-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : publicContent?.data?.length === 0 ? (
              <Card className="text-center py-12">
                <CardContent>
                  <Star className="mx-auto h-12 w-12 text-gray-400 mb-4" />
                  <h3 className="text-lg font-medium text-gray-900 mb-2">No public content</h3>
                  <p className="text-gray-600">There's no public content available to discover yet.</p>
                </CardContent>
              </Card>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {publicContent?.data?.slice(0, 9).map((content: Content) => (
                  <Card key={content.id} className="hover:shadow-md transition-shadow">
                    <CardHeader>
                      <div className="flex items-center space-x-2">
                        <span className="text-2xl">{getContentIcon(content.type)}</span>
                        <div>
                          <CardTitle className="text-lg">{content.title}</CardTitle>
                          <CardDescription className="text-sm">
                            by {content.user?.username}
                          </CardDescription>
                        </div>
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-3">
                        <p className="text-sm text-gray-600 line-clamp-2">
                          {content.description || 'No description available'}
                        </p>
                        
                        {content.tags && content.tags.length > 0 && (
                          <div className="flex flex-wrap gap-1">
                            {content.tags.slice(0, 3).map((tag, index) => (
                              <Badge key={index} variant="secondary" className="text-xs">
                                {tag}
                              </Badge>
                            ))}
                          </div>
                        )}

                        <div className="flex items-center justify-between text-sm text-gray-500">
                          <span>{content.type}</span>
                          <span>{formatDate(content.created_at)}</span>
                        </div>

                        <Button variant="outline" size="sm" className="w-full">
                          View Content
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}