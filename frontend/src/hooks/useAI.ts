import { useState, useCallback, useEffect } from 'react'
import { useAuth } from './useAuth'
import { api } from '@services/api'

interface AIGenerationRequest {
  type: string
  prompt: string
  context?: string
  length?: number
  style?: string
  language?: string
  metadata?: Record<string, any>
}

interface AIGenerationResponse {
  content: string
  title?: string
  description?: string
  metadata?: Record<string, any>
  model: string
  tokens: number
  cost: number
  latency: string
}

interface AISuggestionRequest {
  content: string
  type: string
  context?: string
}

interface AISuggestionResponse {
  suggestions: Array<{
    type: string
    content: string
    confidence: number
    explanation?: string
  }>
}

interface AITemplateRequest {
  type: string
  category: string
  description?: string
}

interface AITemplateResponse {
  name: string
  description: string
  content: string
  type: string
  category: string
  tags: string[]
  metadata?: Record<string, any>
}

interface AIModelStatus {
  openai?: {
    available: boolean
    model: string
  }
  anthropic?: {
    available: boolean
    model: string
  }
  local?: {
    available: boolean
    model: string
  }
}

export const useAI = () => {
  const { user } = useAuth()
  const [isAvailable, setIsAvailable] = useState(false)
  const [modelStatus, setModelStatus] = useState<AIModelStatus>({})
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Check AI availability on mount
  useEffect(() => {
    checkAIAvailability()
  }, [])

  // Check if AI services are available
  const checkAIAvailability = useCallback(async () => {
    try {
      const response = await api.get('/ai/status')
      if (response.data) {
        setModelStatus(response.data)
        // Check if any model is available
        const hasAvailableModel = Object.values(response.data).some(
          (model: any) => model?.available
        )
        setIsAvailable(hasAvailableModel)
      }
    } catch (err) {
      console.warn('AI services not available:', err)
      setIsAvailable(false)
    }
  }, [])

  // Generate AI content
  const generateContent = useCallback(async (
    request: AIGenerationRequest
  ): Promise<AIGenerationResponse | null> => {
    if (!isAvailable || !user) {
      setError('AI services not available or user not authenticated')
      return null
    }

    setIsLoading(true)
    setError(null)

    try {
      const response = await api.post('/ai/generate', {
        ...request,
        user_id: user.id,
        collaboration_id: request.metadata?.collaboration_id
      })

      if (response.data) {
        return response.data
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to generate content'
      setError(errorMessage)
      console.error('AI generation error:', err)
    } finally {
      setIsLoading(false)
    }

    return null
  }, [isAvailable, user])

  // Generate AI suggestions
  const generateSuggestions = useCallback(async (
    request: AISuggestionRequest
  ): Promise<AISuggestionResponse | null> => {
    if (!isAvailable || !user) {
      setError('AI services not available or user not authenticated')
      return null
    }

    setIsLoading(true)
    setError(null)

    try {
      const response = await api.post('/ai/suggestions', {
        ...request,
        user_id: user.id
      })

      if (response.data) {
        return response.data
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to generate suggestions'
      setError(errorMessage)
      console.error('AI suggestions error:', err)
    } finally {
      setIsLoading(false)
    }

    return null
  }, [isAvailable, user])

  // Generate AI template
  const generateTemplate = useCallback(async (
    request: AITemplateRequest
  ): Promise<AITemplateResponse | null> => {
    if (!isAvailable || !user) {
      setError('AI services not available or user not authenticated')
      return null
    }

    setIsLoading(true)
    setError(null)

    try {
      const response = await api.post('/ai/templates', {
        ...request,
        user_id: user.id
      })

      if (response.data) {
        return response.data
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.message || 'Failed to generate template'
      setError(errorMessage)
      console.error('AI template error:', err)
    } finally {
      setIsLoading(false)
    }

    return null
  }, [isAvailable, user])

  // Improve existing content
  const improveContent = useCallback(async (
    content: string,
    improvementType: 'grammar' | 'style' | 'clarity' | 'professional'
  ): Promise<string | null> => {
    if (!isAvailable || !user) {
      setError('AI services not available or user not authenticated')
      return null
    }

    setIsLoading(true)
    setError(null)

    try {
      const prompt = `Improve the following content for better ${improvementType}:\n\n${content}`
      
      const response = await generateContent({
        type: 'document',
        prompt: prompt,
        context: content,
        length: content.length + 100,
        style: improvementType,
        metadata: { improvement_type: improvementType }
      })

      return response?.content || null
    } catch (err) {
      const errorMessage = 'Failed to improve content'
      setError(errorMessage)
      console.error('Content improvement error:', err)
      return null
    } finally {
      setIsLoading(false)
    }
  }, [isAvailable, user, generateContent])

  // Translate content
  const translateContent = useCallback(async (
    content: string,
    targetLanguage: string,
    sourceLanguage?: string
  ): Promise<string | null> => {
    if (!isAvailable || !user) {
      setError('AI services not available or user not authenticated')
      return null
    }

    setIsLoading(true)
    setError(null)

    try {
      const prompt = `Translate the following content to ${targetLanguage}${sourceLanguage ? ` from ${sourceLanguage}` : ''}:\n\n${content}`
      
      const response = await generateContent({
        type: 'document',
        prompt: prompt,
        context: content,
        length: content.length + 50,
        language: targetLanguage,
        metadata: { 
          translation: true,
          target_language: targetLanguage,
          source_language: sourceLanguage
        }
      })

      return response?.content || null
    } catch (err) {
      const errorMessage = 'Failed to translate content'
      setError(errorMessage)
      console.error('Content translation error:', err)
      return null
    } finally {
      setIsLoading(false)
    }
  }, [isAvailable, user, generateContent])

  // Summarize content
  const summarizeContent = useCallback(async (
    content: string,
    maxLength: number = 200
  ): Promise<string | null> => {
    if (!isAvailable || !user) {
      setError('AI services not available or user not authenticated')
      return null
    }

    setIsLoading(true)
    setError(null)

    try {
      const prompt = `Summarize the following content in ${maxLength} words or less:\n\n${content}`
      
      const response = await generateContent({
        type: 'document',
        prompt: prompt,
        context: content,
        length: maxLength,
        style: 'concise',
        metadata: { summary: true, max_length: maxLength }
      })

      return response?.content || null
    } catch (err) {
      const errorMessage = 'Failed to summarize content'
      setError(errorMessage)
      console.error('Content summarization error:', err)
      return null
    } finally {
      setIsLoading(false)
    }
  }, [isAvailable, user, generateContent])

  // Generate code with AI
  const generateCode = useCallback(async (
    description: string,
    language: string,
    framework?: string
  ): Promise<string | null> => {
    if (!isAvailable || !user) {
      setError('AI services not available or user not authenticated')
      return null
    }

    setIsLoading(true)
    setError(null)

    try {
      const prompt = `Generate ${language} code${framework ? ` using ${framework}` : ''} for: ${description}`
      
      const response = await generateContent({
        type: 'code',
        prompt: prompt,
        language: language,
        length: 1000,
        style: 'clean',
        metadata: { 
          code_generation: true,
          language: language,
          framework: framework
        }
      })

      return response?.content || null
    } catch (err) {
      const errorMessage = 'Failed to generate code'
      setError(errorMessage)
      console.error('Code generation error:', err)
      return null
    } finally {
      setIsLoading(false)
    }
  }, [isAvailable, user, generateContent])

  // Clear error
  const clearError = useCallback(() => {
    setError(null)
  }, [])

  // Retry AI connection
  const retryConnection = useCallback(() => {
    checkAIAvailability()
  }, [checkAIAvailability])

  return {
    // State
    isAvailable,
    modelStatus,
    isLoading,
    error,
    
    // Core AI functions
    generateContent,
    generateSuggestions,
    generateTemplate,
    
    // Specialized AI functions
    improveContent,
    translateContent,
    summarizeContent,
    generateCode,
    
    // Utility functions
    clearError,
    retryConnection,
    checkAIAvailability
  }
}