'use client'
import { useState } from "react";

export default function Home() {
  const [url, setUrl] = useState('')
  const handleSubmit = async (e: React.FormEvent) => {

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
            </div>
          </form>
        </div>
      </div>
    </div>

  );
}
