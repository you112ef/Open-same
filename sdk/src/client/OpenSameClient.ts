import { EventEmitter } from 'eventemitter3'
import { OpenSameConfig, ClientOptions, AuthOptions } from '../types/config'
import { AuthClient } from './AuthClient'
import { ContentClient } from './ContentClient'
import { UserClient } from './UserClient'
import { CollaborationClient } from './CollaborationClient'
import { WebSocketClient } from './WebSocketClient'

export class OpenSameClient extends EventEmitter {
  private config: OpenSameConfig
  private authClient: AuthClient
  private contentClient: ContentClient
  private userClient: UserClient
  private collaborationClient: CollaborationClient
  private wsClient: WebSocketClient | null = null

  constructor(config: OpenSameConfig, options?: ClientOptions) {
    super()
    
    this.config = {
      apiUrl: 'http://localhost:8080',
      wsUrl: 'ws://localhost:8080/ws',
      timeout: 30000,
      retries: 3,
      ...config
    }

    // Initialize clients
    this.authClient = new AuthClient(this.config)
    this.contentClient = new ContentClient(this.config)
    this.userClient = new UserClient(this.config)
    this.collaborationClient = new CollaborationClient(this.config)

    // Set up event listeners
    this.setupEventListeners()
  }

  /**
   * Get the authentication client
   */
  get auth(): AuthClient {
    return this.authClient
  }

  /**
   * Get the content management client
   */
  get content(): ContentClient {
    return this.contentClient
  }

  /**
   * Get the user management client
   */
  get user(): UserClient {
    return this.userClient
  }

  /**
   * Get the collaboration client
   */
  get collaboration(): CollaborationClient {
    return this.collaborationClient
  }

  /**
   * Get the WebSocket client for real-time features
   */
  get websocket(): WebSocketClient | null {
    return this.wsClient
  }

  /**
   * Initialize the WebSocket connection
   */
  async connectWebSocket(authOptions?: AuthOptions): Promise<WebSocketClient> {
    if (this.wsClient) {
      return this.wsClient
    }

    const token = authOptions?.token || this.authClient.getAccessToken()
    if (!token) {
      throw new Error('Authentication token required for WebSocket connection')
    }

    this.wsClient = new WebSocketClient(this.config.wsUrl, {
      token,
      userId: authOptions?.userId,
      username: authOptions?.username,
    })

    // Forward WebSocket events
    this.wsClient.on('connected', () => this.emit('websocket:connected'))
    this.wsClient.on('disconnected', () => this.emit('websocket:disconnected'))
    this.wsClient.on('error', (error) => this.emit('websocket:error', error))
    this.wsClient.on('message', (message) => this.emit('websocket:message', message))

    await this.wsClient.connect()
    return this.wsClient
  }

  /**
   * Disconnect the WebSocket connection
   */
  async disconnectWebSocket(): Promise<void> {
    if (this.wsClient) {
      await this.wsClient.disconnect()
      this.wsClient = null
    }
  }

  /**
   * Check if the client is authenticated
   */
  isAuthenticated(): boolean {
    return this.authClient.isAuthenticated()
  }

  /**
   * Get the current user information
   */
  async getCurrentUser() {
    if (!this.isAuthenticated()) {
      throw new Error('Not authenticated')
    }
    return this.userClient.getProfile()
  }

  /**
   * Refresh the authentication token
   */
  async refreshToken(): Promise<void> {
    await this.authClient.refreshToken()
  }

  /**
   * Logout and clear all authentication data
   */
  async logout(): Promise<void> {
    try {
      await this.authClient.logout()
    } catch (error) {
      // Continue with cleanup even if logout request fails
    } finally {
      // Clear local storage
      if (typeof window !== 'undefined') {
        localStorage.removeItem('open_same_access_token')
        localStorage.removeItem('open_same_refresh_token')
        localStorage.removeItem('open_same_user')
      }
      
      // Disconnect WebSocket
      await this.disconnectWebSocket()
      
      // Emit logout event
      this.emit('auth:logout')
    }
  }

  /**
   * Get the current configuration
   */
  getConfig(): OpenSameConfig {
    return { ...this.config }
  }

  /**
   * Update the configuration
   */
  updateConfig(newConfig: Partial<OpenSameConfig>): void {
    this.config = { ...this.config, ...newConfig }
    
    // Update client configurations
    this.authClient.updateConfig(this.config)
    this.contentClient.updateConfig(this.config)
    this.userClient.updateConfig(this.config)
    this.collaborationClient.updateConfig(this.config)
    
    this.emit('config:updated', this.config)
  }

  /**
   * Health check for the API
   */
  async healthCheck(): Promise<{ status: string; timestamp: string; version: string }> {
    try {
      const response = await fetch(`${this.config.apiUrl}/health`)
      if (!response.ok) {
        throw new Error(`Health check failed: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      throw new Error(`Health check failed: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  /**
   * Get API status and information
   */
  async getApiInfo(): Promise<{
    version: string
    environment: string
    features: string[]
    limits: Record<string, any>
  }> {
    try {
      const response = await fetch(`${this.config.apiUrl}/api/v1/info`)
      if (!response.ok) {
        throw new Error(`Failed to get API info: ${response.status}`)
      }
      return await response.json()
    } catch (error) {
      throw new Error(`Failed to get API info: ${error instanceof Error ? error.message : 'Unknown error'}`)
    }
  }

  /**
   * Set up event listeners for authentication state changes
   */
  private setupEventListeners(): void {
    // Listen for authentication events
    this.authClient.on('auth:login', (user) => {
      this.emit('auth:login', user)
    })

    this.authClient.on('auth:logout', () => {
      this.emit('auth:logout')
    })

    this.authClient.on('auth:token_refreshed', (tokens) => {
      this.emit('auth:token_refreshed', tokens)
    })

    this.authClient.on('auth:error', (error) => {
      this.emit('auth:error', error)
    })

    // Listen for content events
    this.contentClient.on('content:created', (content) => {
      this.emit('content:created', content)
    })

    this.contentClient.on('content:updated', (content) => {
      this.emit('content:updated', content)
    })

    this.contentClient.on('content:deleted', (contentId) => {
      this.emit('content:deleted', contentId)
    })

    // Listen for collaboration events
    this.collaborationClient.on('collaboration:joined', (collaboration) => {
      this.emit('collaboration:joined', collaboration)
    })

    this.collaborationClient.on('collaboration:left', (collaboration) => {
      this.emit('collaboration:left', collaboration)
    })

    this.collaborationClient.on('collaboration:updated', (collaboration) => {
      this.emit('collaboration:updated', collaboration)
    })
  }

  /**
   * Clean up resources
   */
  async destroy(): Promise<void> {
    // Disconnect WebSocket
    await this.disconnectWebSocket()
    
    // Remove all event listeners
    this.removeAllListeners()
    
    // Clean up clients
    this.authClient.removeAllListeners()
    this.contentClient.removeAllListeners()
    this.userClient.removeAllListeners()
    this.collaborationClient.removeAllListeners()
  }
}

// Export a default instance factory
export const createOpenSameClient = (config: OpenSameConfig, options?: ClientOptions): OpenSameClient => {
  return new OpenSameClient(config, options)
}

// Export the class and factory function
export default OpenSameClient