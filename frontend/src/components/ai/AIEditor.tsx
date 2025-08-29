import React, { useState, useRef, useEffect } from 'react'
import { useAI } from '@hooks/useAI'
import { Button } from '@components/ui/Button'
import { Textarea } from '@components/ui/Textarea'
import { Select } from '@components/ui/Select'
import { Badge } from '@components/ui/Badge'
import { Card } from '@components/ui/Card'
import { 
  Sparkles, 
  Wand2, 
  Lightbulb, 
  Code, 
  FileText, 
  Image, 
  Brain,
  Loader2,
  CheckCircle,
  AlertCircle
} from 'lucide-react'

interface AIEditorProps {
  content: string
  contentType: 'document' | 'code' | 'diagram' | 'template'
  onContentChange: (content: string) => void
  onSave: () => void
  collaborationId?: string
}

export const AIEditor: React.FC<AIEditorProps> = ({
  content,
  contentType,
  onContentChange,
  onSave,
  collaborationId
}) => {
  const [prompt, setPrompt] = useState('')
  const [isGenerating, setIsGenerating] = useState(false)
  const [suggestions, setSuggestions] = useState<any[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [selectedSuggestion, setSelectedSuggestion] = useState<number | null>(null)
  
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const { generateContent, generateSuggestions, isAvailable } = useAI()

  // AI content types with appropriate prompts
  const contentTypes = {
    document: {
      icon: FileText,
      label: 'Document',
      defaultPrompt: 'Write a professional document about',
      examples: ['business proposal', 'research paper', 'technical manual', 'creative story']
    },
    code: {
      icon: Code,
      label: 'Code',
      defaultPrompt: 'Generate code for',
      examples: ['React component', 'API endpoint', 'database query', 'algorithm']
    },
    diagram: {
      icon: Image,
      label: 'Diagram',
      defaultPrompt: 'Create a diagram showing',
      examples: ['system architecture', 'user flow', 'data model', 'process workflow']
    },
    template: {
      icon: FileText,
      label: 'Template',
      defaultPrompt: 'Create a template for',
      examples: ['email campaign', 'project plan', 'meeting agenda', 'report format']
    }
  }

  const currentType = contentTypes[contentType]

  // Generate AI content
  const handleGenerate = async () => {
    if (!prompt.trim() || !isAvailable) return

    setIsGenerating(true)
    try {
      const result = await generateContent({
        type: contentType,
        prompt: prompt,
        context: content,
        length: 500,
        style: 'professional'
      })

      if (result?.content) {
        onContentChange(result.content)
        setPrompt('')
        // Show success feedback
        setSuggestions([{
          type: 'generation',
          content: 'Content generated successfully!',
          confidence: 1.0
        }])
        setShowSuggestions(true)
        setTimeout(() => setShowSuggestions(false), 3000)
      }
    } catch (error) {
      console.error('AI generation failed:', error)
      setSuggestions([{
        type: 'error',
        content: 'Failed to generate content. Please try again.',
        confidence: 0
      }])
      setShowSuggestions(true)
    } finally {
      setIsGenerating(false)
    }
  }

  // Generate AI suggestions
  const handleGenerateSuggestions = async () => {
    if (!content.trim() || !isAvailable) return

    setIsGenerating(true)
    try {
      const result = await generateSuggestions({
        content: content,
        type: contentType
      })

      if (result?.suggestions) {
        setSuggestions(result.suggestions)
        setShowSuggestions(true)
      }
    } catch (error) {
      console.error('AI suggestions failed:', error)
    } finally {
      setIsGenerating(false)
    }
  }

  // Apply suggestion
  const handleApplySuggestion = (suggestion: any) => {
    if (suggestion.content) {
      onContentChange(suggestion.content)
      setShowSuggestions(false)
      setSelectedSuggestion(null)
    }
  }

  // Quick prompt templates
  const quickPrompts = [
    'Improve this content',
    'Make it more professional',
    'Add more details',
    'Simplify the language',
    'Make it more engaging'
  ]

  const handleQuickPrompt = (quickPrompt: string) => {
    setPrompt(`${quickPrompt}: ${content.substring(0, 100)}...`)
  }

  // Auto-save on content change
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (content.trim()) {
        onSave()
      }
    }, 2000)

    return () => clearTimeout(timeoutId)
  }, [content, onSave])

  if (!isAvailable) {
    return (
      <Card className="p-6 text-center">
        <AlertCircle className="w-12 h-12 text-yellow-500 mx-auto mb-4" />
        <h3 className="text-lg font-semibold mb-2">AI Features Unavailable</h3>
        <p className="text-gray-600 mb-4">
          AI services are not configured. Please check your API keys and configuration.
        </p>
        <Button variant="outline" onClick={() => window.location.reload()}>
          Retry Connection
        </Button>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* AI Generation Panel */}
      <Card className="p-6">
        <div className="flex items-center gap-2 mb-4">
          <Sparkles className="w-5 h-5 text-purple-500" />
          <h3 className="text-lg font-semibold">AI Content Generation</h3>
          <Badge variant="secondary" className="ml-auto">
            {currentType.label}
          </Badge>
        </div>

        <div className="space-y-4">
          {/* Prompt Input */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Describe what you want to create
            </label>
            <Textarea
              ref={textareaRef}
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              placeholder={`${currentType.defaultPrompt}...`}
              className="min-h-[100px]"
              disabled={isGenerating}
            />
          </div>

          {/* Quick Prompts */}
          <div>
            <label className="block text-sm font-medium mb-2">
              Quick Actions
            </label>
            <div className="flex flex-wrap gap-2">
              {quickPrompts.map((quickPrompt, index) => (
                <Button
                  key={index}
                  variant="outline"
                  size="sm"
                  onClick={() => handleQuickPrompt(quickPrompt)}
                  disabled={isGenerating}
                >
                  {quickPrompt}
                </Button>
              ))}
            </div>
          </div>

          {/* Generate Button */}
          <div className="flex gap-3">
            <Button
              onClick={handleGenerate}
              disabled={!prompt.trim() || isGenerating}
              className="flex-1"
            >
              {isGenerating ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Generating...
                </>
              ) : (
                <>
                  <Wand2 className="w-4 h-4 mr-2" />
                  Generate Content
                </>
              )}
            </Button>

            <Button
              variant="outline"
              onClick={handleGenerateSuggestions}
              disabled={!content.trim() || isGenerating}
            >
              <Lightbulb className="w-4 h-4 mr-2" />
              Get Suggestions
            </Button>
          </div>
        </div>
      </Card>

      {/* Content Editor */}
      <Card className="p-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold">Content Editor</h3>
          <div className="flex items-center gap-2">
            <Brain className="w-4 h-4 text-blue-500" />
            <span className="text-sm text-gray-600">AI-Powered</span>
          </div>
        </div>

        <Textarea
          value={content}
          onChange={(e) => onContentChange(e.target.value)}
          placeholder={`Start writing your ${contentType}...`}
          className="min-h-[400px] font-mono"
          disabled={isGenerating}
        />

        <div className="flex justify-between items-center mt-4">
          <div className="text-sm text-gray-500">
            {content.length} characters
          </div>
          <Button onClick={onSave} variant="outline">
            Save
          </Button>
        </div>
      </Card>

      {/* AI Suggestions */}
      {showSuggestions && suggestions.length > 0 && (
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <Lightbulb className="w-5 h-5 text-yellow-500" />
            <h3 className="text-lg font-semibold">AI Suggestions</h3>
          </div>

          <div className="space-y-3">
            {suggestions.map((suggestion, index) => (
              <div
                key={index}
                className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                  selectedSuggestion === index
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 hover:border-gray-300'
                }`}
                onClick={() => setSelectedSuggestion(index)}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <Badge variant="outline" className="text-xs">
                        {suggestion.type}
                      </Badge>
                      {suggestion.confidence > 0.8 && (
                        <CheckCircle className="w-4 h-4 text-green-500" />
                      )}
                    </div>
                    <p className="text-sm text-gray-700">{suggestion.content}</p>
                    {suggestion.explanation && (
                      <p className="text-xs text-gray-500 mt-2">
                        {suggestion.explanation}
                      </p>
                    )}
                  </div>
                  <Button
                    size="sm"
                    onClick={(e) => {
                      e.stopPropagation()
                      handleApplySuggestion(suggestion)
                    }}
                  >
                    Apply
                  </Button>
                </div>
              </div>
            ))}
          </div>

          <div className="flex justify-end mt-4">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowSuggestions(false)}
            >
              Close
            </Button>
          </div>
        </Card>
      )}

      {/* Content Type Examples */}
      <Card className="p-6">
        <h3 className="text-lg font-semibold mb-4">Examples</h3>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
          {currentType.examples.map((example, index) => (
            <Button
              key={index}
              variant="outline"
              size="sm"
              onClick={() => setPrompt(`${currentType.defaultPrompt} ${example}`)}
              className="text-left h-auto p-3"
            >
              <div className="text-sm font-medium">{example}</div>
            </Button>
          ))}
        </div>
      </Card>
    </div>
  )
}