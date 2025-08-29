export { OpenSameClient } from './client/OpenSameClient'
export { CollaborationClient } from './client/CollaborationClient'
export { ContentClient } from './client/ContentClient'
export { UserClient } from './client/UserClient'
export { AuthClient } from './client/AuthClient'

export type {
  OpenSameConfig,
  ClientOptions,
  AuthOptions
} from './types/config'

export type {
  User,
  Content,
  Collaboration,
  ContentVersion,
  SharedContent,
  CreateContentRequest,
  UpdateContentRequest,
  CollaborationRequest,
  ShareContentRequest
} from './types/models'

export type {
  ApiResponse,
  PaginatedResponse,
  ErrorResponse,
  SuccessResponse
} from './types/api'

export type {
  CollaborationEvent,
  CollaborationMessage,
  CursorPosition,
  Selection,
  Change,
  Operation
} from './types/collaboration'

export type {
  WebSocketMessage,
  WebSocketEvent,
  ConnectionStatus
} from './types/websocket'

// Re-export commonly used utilities
export { EventEmitter } from 'eventemitter3'
export { default as axios } from 'axios'

// Version information
export const VERSION = '1.0.0'
export const API_VERSION = 'v1'