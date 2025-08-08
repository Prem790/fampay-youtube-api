import { useState, useRef, useEffect } from 'react'
import { Search, X, Sparkles, Clock, Globe, Database } from 'lucide-react'

export default function SearchBar({ onSearch, onClear, isSearchMode, searchQuery, searchType }) {
  const [query, setQuery] = useState(searchQuery || '')
  const [recentSearches, setRecentSearches] = useState([])
  const inputRef = useRef(null)

  useEffect(() => {
    // Load recent searches from localStorage
    const saved = localStorage.getItem('recentSearches')
    if (saved) {
      setRecentSearches(JSON.parse(saved).slice(0, 5))
    }
  }, [])

  const handleSubmit = (e) => {
    e.preventDefault()
    if (query.trim()) {
      const trimmedQuery = query.trim()
      onSearch(trimmedQuery)
      
      // Save to recent searches
      const updated = [trimmedQuery, ...recentSearches.filter(s => s !== trimmedQuery)].slice(0, 5)
      setRecentSearches(updated)
      localStorage.setItem('recentSearches', JSON.stringify(updated))
    }
  }

  const handleClear = () => {
    setQuery('')
    onClear()
    inputRef.current?.focus()
  }

  const handleKeyDown = (e) => {
    if (e.key === 'Escape') {
      handleClear()
    }
  }

  const handleQuickSearch = (term) => {
    setQuery(term)
    onSearch(term)
    
    // Save to recent searches
    const updated = [term, ...recentSearches.filter(s => s !== term)].slice(0, 5)
    setRecentSearches(updated)
    localStorage.setItem('recentSearches', JSON.stringify(updated))
  }

  // Updated search suggestions for more diverse content
  const popularSearches = [
    'talkfootballhd', 'how to make tea', 'cooking recipes', 'tech reviews', 
    'gaming highlights', 'travel vlogs', 'music videos', 'movie trailers',
    'programming tutorials', 'fitness workouts', 'comedy sketches', 'science experiments'
  ]

  return (
    <div className="space-y-4">
      <form onSubmit={handleSubmit} className="relative">
        <div className="relative">
          <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="Search anything on YouTube... (channels, topics, tutorials, etc.)"
            className="w-full pl-12 pr-12 py-4 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-900 placeholder-gray-500 text-lg"
          />
          {(query || isSearchMode) && (
            <button
              type="button"
              onClick={handleClear}
              className="absolute right-4 top-1/2 transform -translate-y-1/2 h-6 w-6 text-gray-400 hover:text-gray-600 transition-colors duration-200 bg-gray-100 rounded-full flex items-center justify-center"
            >
              <X className="h-4 w-4" />
            </button>
          )}
        </div>
        
        {/* Search button for mobile */}
        <button
          type="submit"
          className="md:hidden mt-4 w-full bg-blue-600 text-white py-3 rounded-xl hover:bg-blue-700 transition-colors duration-200 font-medium"
        >
          Search YouTube
        </button>
      </form>

      {/* Search Mode Info */}
      {isSearchMode && (
        <div className={`border rounded-lg p-3 ${
          searchType === 'live' 
            ? 'bg-red-50 border-red-200' 
            : 'bg-blue-50 border-blue-200'
        }`}>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              {searchType === 'live' ? (
                <>
                  <Globe className="h-4 w-4 text-red-600" />
                  <span className="text-sm text-red-800">
                    <span className="font-semibold">Live YouTube Search:</span> "{searchQuery}"
                  </span>
                </>
              ) : (
                <>
                  <Database className="h-4 w-4 text-blue-600" />
                  <span className="text-sm text-blue-800">
                    <span className="font-semibold">Stored Videos Search:</span> "{searchQuery}"
                  </span>
                </>
              )}
            </div>
            <button
              onClick={handleClear}
              className={`text-sm font-medium transition-colors duration-200 ${
                searchType === 'live' 
                  ? 'text-red-600 hover:text-red-700' 
                  : 'text-blue-600 hover:text-blue-700'
              }`}
            >
              Clear
            </button>
          </div>
        </div>
      )}

      {/* Search Suggestions */}
      {!isSearchMode && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Popular Searches */}
          <div>
            <div className="flex items-center space-x-2 mb-3">
              <Sparkles className="h-4 w-4 text-yellow-500" />
              <span className="text-sm font-medium text-gray-700">Try These Searches</span>
            </div>
            <div className="flex flex-wrap gap-2">
              {popularSearches.slice(0, 6).map((term) => (
                <button
                  key={term}
                  onClick={() => handleQuickSearch(term)}
                  className="px-3 py-1.5 text-sm bg-gradient-to-r from-blue-50 to-purple-50 text-blue-700 rounded-full hover:from-blue-100 hover:to-purple-100 transition-all duration-200 border border-blue-200"
                >
                  {term}
                </button>
              ))}
            </div>
          </div>

          {/* Recent Searches */}
          {recentSearches.length > 0 && (
            <div>
              <div className="flex items-center space-x-2 mb-3">
                <Clock className="h-4 w-4 text-gray-500" />
                <span className="text-sm font-medium text-gray-700">Recent Searches</span>
              </div>
              <div className="flex flex-wrap gap-2">
                {recentSearches.slice(0, 4).map((term, index) => (
                  <button
                    key={index}
                    onClick={() => handleQuickSearch(term)}
                    className="px-3 py-1.5 text-sm bg-gray-100 text-gray-700 rounded-full hover:bg-gray-200 transition-colors duration-200"
                  >
                    {term}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Feature Explanation */}
      {!isSearchMode && (
        <div className="bg-gradient-to-r from-blue-50 to-purple-50 border border-blue-200 rounded-lg p-4">
          <h3 className="text-sm font-semibold text-gray-900 mb-2">🚀 Two Search Modes Available:</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            <div className="flex items-start space-x-2">
              <Globe className="h-4 w-4 text-red-500 mt-0.5" />
              <div>
                <div className="font-medium text-red-700">Live YouTube Search</div>
                <div className="text-gray-600">Search anything on YouTube in real-time</div>
              </div>
            </div>
            <div className="flex items-start space-x-2">
              <Database className="h-4 w-4 text-blue-500 mt-0.5" />
              <div>
                <div className="font-medium text-blue-700">Stored Videos Search</div>
                <div className="text-gray-600">Search through our curated collection</div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}