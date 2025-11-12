'use client'
import { useState } from "react";

interface AnalysisResponse {
  htmlVersion: string
  pageTitle: string
  headings: {
    h1: number
    h2: number
    h3: number
    h4: number
    h5: number
    h6: number
  }
  linkData: {
    internalLinks: number
    externalLinks: number
    unAccessibleLinks: number
  }
  loginFormAvailability: boolean
  error?: string
}

export default function Home() {
  const [url, setUrl] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [result, setResult] = useState<AnalysisResponse | null>(null)


  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!url.trim()) {
      setError('Please enter a URL')
      return
    }

    // Basic URL validation
    try {
      new URL(url)
    } catch {
      setError('Please enter a valid URL')
      return
    }
    setError(null)
    setResult(null)

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_APP_API_URL}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: url.trim() }),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data: AnalysisResponse = await response.json()

      if (data.error) {
        throw new Error(data.error)
      }
      console.log(`Data Loaded from API ${JSON.stringify(data)}`)
      setResult(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to analyze URL')
    }

  }
  return (
    <div className="min-h-screen bg-white py-8">
      <div className="max-w-4xl mx-auto px-4">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Web Analyser - Golang Web UI</h1>
          <p className="text-gray-600">
            Analyze any website's with golang API
          </p>
        </div>
        <div>
          <form onSubmit={handleSubmit}>
            <div>
              <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-2">
                Website URL
              </label>
              <input
                type="url"
                id="url"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="https://example.com"
                className="w-full px-3 py-2 border border-gray-300"
              />
              <div className="flex gap-20">
                <button
                  type="submit"
                  className="flex-1 bg-blue-600 text-white py-2 px-4"
                >
                  Analyze Website
                </button>
              </div>
              {error && (
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
                  {error}
                </div>
              )}
            </div>
          </form>
        </div>
        {result && (
          <div>
            <h2 className="text-xl font-bold text-green-900 mb-6">Analysis Results</h2>
            <div className="grid grid-cols-2 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <h3 className="text-red-700">Basic Informations</h3>
                <div className="bg-blue-50 rounded-lg p-4">
                  <div className="flex justify-between">
                    <span className="text-blue-700">HTML Version:</span>
                    <span className="font-medium text-green-700">{result.htmlVersion}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">Page Title:</span>
                    <span className="font-medium text-green-700">{result.pageTitle}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">Login Form Availability:</span>
                    <span className="font-medium text-green-700">{result.loginFormAvailability}</span>
                  </div>
                </div>
              </div>
              <div className="space-y-4">
                <h3 className="text-red-700">Headings Informations</h3>
                <div className="bg-blue-50 rounded-lg p-4">
                  <div className="flex justify-between">
                    <span className="text-blue-700">H1:</span>
                    <span className="font-medium text-green-700">{result.headings.h1}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">H2:</span>
                    <span className="font-medium text-green-700">{result.headings.h2}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">H3:</span>
                    <span className="font-medium text-green-700">{result.headings.h3}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">H4:</span>
                    <span className="font-medium text-green-700">{result.headings.h4}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">H5:</span>
                    <span className="font-medium text-green-700">{result.headings.h5}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">H6:</span>
                    <span className="font-medium text-green-700">{result.headings.h6}</span>
                  </div>
                </div>
              </div>
              <div className="space-y-4">
                <h3 className="text-red-700">Link Informations</h3>
                <div className="bg-blue-50 rounded-lg p-4">
                  <div className="flex justify-between">
                    <span className="text-blue-700">Internal Links:</span>
                    <span className="font-medium text-green-700">{result.linkData.internalLinks}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">External Links:</span>
                    <span className="font-medium text-green-700">{result.linkData.externalLinks}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-blue-700">Unaccessible Links:</span>
                    <span className="font-medium text-green-700">{result.linkData.unAccessibleLinks}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>

  );
}
