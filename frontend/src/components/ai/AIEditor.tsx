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
  AlertCircle,
  Menu,
  X,
  Smartphone
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
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [activeTab, setActiveTab] = useState<'editor' | 'ai' | 'suggestions'>('editor')
  
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
        setActiveTab('suggestions')
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
      setActiveTab('suggestions')
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
        setActiveTab('suggestions')
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
      setActiveTab('editor')
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
    setActiveTab('ai')
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
      <Card className="p-4 sm:p-6 text-center">
        <AlertCircle className="w-8 h-8 sm:w-12 sm:h-12 text-yellow-500 mx-auto mb-4" />
        <h3 className="text-base sm:text-lg font-semibold mb-2">AI Features Unavailable</h3>
        <p className="text-sm sm:text-base text-gray-600 mb-4">
          AI services are not configured. Please check your API keys and configuration.
        </p>
        <Button variant="outline" onClick={() => window.location.reload()}>
          Retry Connection
        </Button>
      </Card>
    )
  }

  return (
    <div className="space-y-4 sm:space-y-6">
      {/* Mobile Header with Tabs */}
      <div className="sm:hidden">
        <div className="flex items-center justify-between bg-white border-b border-gray-200 px-4 py-3">
          <div className="flex items-center gap-2">
            <Smartphone className="w-5 h-5 text-blue-600" />
            <span className="font-medium text-gray-900">AI Editor</span>
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
            onClick={() => setActiveTab('editor')}
            className={`flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'editor'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Editor
          </button>
          <button
            onClick={() => setActiveTab('ai')}
            className={`flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'ai'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            AI Tools
          </button>
          <button
            onClick={() => setActiveTab('suggestions')}
            className={`flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-colors ${
              activeTab === 'suggestions'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            }`}
          >
            Suggestions
          </button>
        </div>
      </div>

      {/* Desktop Layout */}
      <div className="hidden sm:block">
        {/* AI Generation Panel */}
        <Card className="p-4 sm:p-6">
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
                    className="text-xs sm:text-sm"
                  >
                    {quickPrompt}
                  </Button>
                ))}
              </div>
            </div>

            {/* Generate Button */}
            <div className="flex flex-col sm:flex-row gap-3">
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
      </div>

      {/* Mobile AI Panel */}
      <div className={`sm:hidden ${activeTab === 'ai' ? 'block' : 'hidden'}`}>
        <Card className="p-4">
          <div className="flex items-center gap-2 mb-4">
            <Sparkles className="w-5 h-5 text-purple-500" />
            <h3 className="text-lg font-semibold">AI Generation</h3>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">
                What do you want to create?
              </label>
              <Textarea
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder={`${currentType.defaultPrompt}...`}
                className="min-h-[80px] text-sm"
                disabled={isGenerating}
              />
            </div>

            <div className="grid grid-cols-2 gap-2">
              {quickPrompts.slice(0, 4).map((quickPrompt, index) => (
                <Button
                  key={index}
                  variant="outline"
                  size="sm"
                  onClick={() => handleQuickPrompt(quickPrompt)}
                  disabled={isGenerating}
                  className="text-xs py-2"
                >
                  {quickPrompt}
                </Button>
              ))}
            </div>

            <Button
              onClick={handleGenerate}
              disabled={!prompt.trim() || isGenerating}
              className="w-full"
            >
              {isGenerating ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Generating...
                </>
              ) : (
                <>
                  <Wand2 className="w-4 h-4 mr-2" />
                  Generate
                </>
              )}
            </Button>
          </div>
        </Card>
      </div>

      {/* Content Editor */}
      <div className={`${activeTab === 'editor' ? 'block' : 'hidden'} sm:block`}>
        <Card className="p-4 sm:p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Content Editor</h3>
            <div className="flex items-center gap-2">
              <Brain className="w-4 h-4 text-blue-500" />
              <span className="text-sm text-gray-600 hidden sm:inline">AI-Powered</span>
            </div>
          </div>

          <Textarea
            value={content}
            onChange={(e) => onContentChange(e.target.value)}
            placeholder={`Start writing your ${contentType}...`}
            className="min-h-[300px] sm:min-h-[400px] font-mono text-sm sm:text-base"
            disabled={isGenerating}
          />

          <div className="flex flex-col sm:flex-row justify-between items-center mt-4 gap-3">
            <div className="text-sm text-gray-500">
              {content.length} characters
            </div>
            <Button onClick={onSave} variant="outline" className="w-full sm:w-auto">
              Save
            </Button>
          </div>
        </Card>
      </div>

      {/* AI Suggestions */}
      <div className={`${activeTab === 'suggestions' ? 'block' : 'hidden'} sm:block`}>
        {showSuggestions && suggestions.length > 0 && (
          <Card className="p-4 sm:p-6">
            <div className="flex items-center gap-2 mb-4">
              <Lightbulb className="w-5 h-5 text-yellow-500" />
              <h3 className="text-lg font-semibold">AI Suggestions</h3>
            </div>

            <div className="space-y-3">
              {suggestions.map((suggestion, index) => (
                <div
                  key={index}
                  className={`p-3 sm:p-4 border rounded-lg cursor-pointer transition-colors ${
                    selectedSuggestion === index
                      ? 'border-blue-500 bg-blue-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                  onClick={() => setSelectedSuggestion(index)}
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-2">
                        <Badge variant="outline" className="text-xs">
                          {suggestion.type}
                        </Badge>
                        {suggestion.confidence > 0.8 && (
                          <CheckCircle className="w-4 h-4 text-green-500 flex-shrink-0" />
                        )}
                      </div>
                      <p className="text-sm text-gray-700 break-words">{suggestion.content}</p>
                      {suggestion.explanation && (
                        <p className="text-xs text-gray-500 mt-2 break-words">
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
                      className="flex-shrink-0"
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
      </div>

      {/* Content Type Examples */}
      <div className="hidden sm:block">
        <Card className="p-4 sm:p-6">
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

      {/* Mobile Examples */}
      <div className="sm:hidden">
        <Card className="p-4">
          <h3 className="text-lg font-semibold mb-3">Quick Examples</h3>
          <div className="grid grid-cols-2 gap-2">
            {currentType.examples.slice(0, 4).map((example, index) => (
              <Button
                key={index}
                variant="outline"
                size="sm"
                onClick={() => {
                  setPrompt(`${currentType.defaultPrompt} ${example}`)
                  setActiveTab('ai')
                }}
                className="text-xs py-2 text-center"
              >
                {example}
              </Button>
            ))}
          </div>
        </Card>
      </div>
    </div>
  )
}