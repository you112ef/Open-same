// User types
export interface User {
  id: string
  email: string
  username: string
  first_name?: string
  last_name?: string
  avatar?: string
  bio?: string
  is_verified: boolean
  is_active: boolean
  is_admin: boolean
  last_login_at?: string
  email_verified_at?: string
  created_at: string
  updated_at: string
  fullName?: string
}

// Content types
export type ContentType = 'text' | 'code' | 'diagram' | 'image' | 'document' | 'template'
export type ContentStatus = 'draft' | 'published' | 'archived' | 'deleted'

export interface Content {
  id: string
  user_id: string
  title: string
  description?: string
  content: string
  type: ContentType
  status: ContentStatus
  is_public: boolean
  is_template: boolean
  tags: string[]
  metadata: Record<string, any>
  ai_generated: boolean
  ai_model?: string
  ai_prompt?: string
  version: number
  parent_id?: string
  created_at: string
  updated_at: string
  user?: User
  parent?: Content
  versions?: ContentVersion[]
  collaborations?: Collaboration[]
  shared_contents?: SharedContent[]
}

export interface ContentVersion {
  id: string
  content_id: string
  version: number
  content: string
  title?: string
  description?: string
  tags?: string[]
  metadata?: Record<string, any>
  created_by: string
  created_at: string
  content_ref?: Content
  user?: User
}

export interface SharedContent {
  id: string
  content_id: string
  owner_id: string
  shared_with: string
  permission: 'read' | 'write' | 'admin'
  expires_at?: string
  created_at: string
  updated_at: string
  content?: Content
  owner?: User
  shared_user?: User
}

export interface Collaboration {
  id: string
  content_id: string
  user_id: string
  role: 'viewer' | 'editor' | 'admin'
  joined_at: string
  last_active?: string
  is_active: boolean
  created_at: string
  updated_at: string
  content?: Content
  user?: User
}

// Authentication types
export interface AuthRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  username: string
  password: string
  first_name?: string
  last_name?: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: User
}

export interface RefreshRequest {
  refresh_token: string
}

// Content creation and update types
export interface CreateContentRequest {
  title: string
  description?: string
  content?: string
  type: ContentType
  is_public?: boolean
  is_template?: boolean
  tags?: string[]
  metadata?: Record<string, any>
  parent_id?: string
}

export interface UpdateContentRequest {
  title?: string
  description?: string
  content?: string
  type?: ContentType
  status?: ContentStatus
  is_public?: boolean
  is_template?: boolean
  tags?: string[]
  metadata?: Record<string, any>
}

// Pagination types
export interface PaginationParams {
  page: number
  per_page: number
  search?: string
  type?: ContentType
  status?: ContentStatus
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
  has_next: boolean
  has_previous: boolean
}

// API response types
export interface ApiResponse<T> {
  message: string
  data: T
  error?: string
  code?: string
}

export interface ErrorResponse {
  error: string
  code: string
  message: string
  details?: any
}

// AI service types
export interface AIGenerateRequest {
  prompt: string
  type: ContentType
  length?: string
  style?: string
  tone?: string
  language?: string
  context?: string
  parameters?: Record<string, any>
}

export interface AIGenerateResponse {
  content: string
  title?: string
  description?: string
  tags?: string[]
  metadata?: Record<string, any>
  model: string
  usage?: {
    prompt_tokens: number
    completion_tokens: number
    total_tokens: number
  }
  error?: string
}

// WebSocket types
export interface WebSocketMessage {
  type: string
  room_id?: string
  user_id?: string
  username?: string
  content?: string
  data?: Record<string, any>
  timestamp: string
}

export interface CollaborationEvent {
  type: 'user_joined' | 'user_left' | 'content_change' | 'cursor_move' | 'selection_change' | 'chat_message'
  room_id: string
  user_id: string
  username: string
  data?: Record<string, any>
  timestamp: string
}

// Form types
export interface LoginFormData {
  email: string
  password: string
}

export interface RegisterFormData {
  email: string
  username: string
  password: string
  confirmPassword: string
  firstName?: string
  lastName?: string
}

// Filter types
export interface ContentFilters {
  type?: ContentType
  status?: ContentStatus
  search?: string
  tags?: string[]
  is_public?: boolean
  is_template?: boolean
  user_id?: string
}

// Sort types
export interface SortOptions {
  field: 'created_at' | 'updated_at' | 'title' | 'type' | 'status'
  direction: 'asc' | 'desc'
}

// Export all types
export type {
  User,
  Content,
  ContentVersion,
  SharedContent,
  Collaboration,
  AuthRequest,
  RegisterRequest,
  AuthResponse,
  RefreshRequest,
  CreateContentRequest,
  UpdateContentRequest,
  PaginationParams,
  PaginatedResponse,
  ApiResponse,
  ErrorResponse,
  AIGenerateRequest,
  AIGenerateResponse,
  WebSocketMessage,
  CollaborationEvent,
  LoginFormData,
  RegisterFormData,
  ContentFilters,
  SortOptions,
}