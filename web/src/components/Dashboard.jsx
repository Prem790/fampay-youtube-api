import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Search, RefreshCw, Video, Globe, Database } from 'lucide-react'
import VideoGrid from './VideoGrid'
import SearchBar from './SearchBar'
import Pagination from './Pagination'
import LoadingSpinner from './LoadingSpinner'
import { fetchVideos, searchVideos, searchYouTubeLive } from '../services/api'

export default function Dashboard() {
  const [currentPage, setCurrentPage] = useState(1)
  const [searchQuery, setSearchQuery] = useState('')
  const [isSearchMode, setIsSearchMode] = useState(false)
  const [searchType, setSearchType] = useState('live') // 'live' or 'stored'
  const [sortBy, setSortBy] = useState('latest')
  const pageSize = 12

  // Videos query with proper dependency tracking
  const {
    data: videosData,
    isLoading,
    error,
    refetch,
    isFetching
  } = useQuery({
    queryKey: ['videos', currentPage, searchQuery, isSearchMode, searchType, sortBy],
    queryFn: () => {
      if (isSearchMode && searchQuery) {
        if (searchType === 'live') {
          return searchYouTubeLive(searchQuery, currentPage, pageSize, sortBy)
        } else {
          return searchVideos(searchQuery, currentPage, pageSize, sortBy)
        }
      }
      return fetchVideos(currentPage, pageSize, sortBy)
    },
    keepPreviousData: true,
    staleTime: 2 * 60 * 1000, // 2 minutes
  })

  const handleSearch = (query) => {
    setSearchQuery(query)
    setIsSearchMode(true)
    setCurrentPage(1)
  }

  const clearSearch = () => {
    setSearchQuery('')
    setIsSearchMode(false)
    setCurrentPage(1)
  }

  const handlePageChange = (page) => {
    setCurrentPage(page)
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const handleSortChange = (newSort) => {
    setSortBy(newSort)
    setCurrentPage(1)
  }

  const handleSearchTypeChange = (type) => {
    setSearchType(type)
    setCurrentPage(1)
    if (isSearchMode) {
      refetch() // Refresh results with new search type
    }
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center bg-red-50 border border-red-200 rounded-lg p-8">
          <div className="text-red-600 text-lg mb-4">
            ⚠️ Error loading videos
          </div>
          <p className="text-red-500 mb-6">{error.message}</p>
          <button 
            onClick={() => refetch()}
            className="bg-red-600 text-white px-6 py-2 rounded-lg hover:bg-red-700 transition-colors duration-200"
          >
            <RefreshCw className="h-4 w-4 mr-2 inline" />
            Try Again
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Enhanced Header with Search Type Toggle */}
      <div className="mb-8">
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center space-x-4">
              <div className="p-3 bg-gradient-to-br from-red-500 to-purple-600 rounded-lg">
                <Video className="h-6 w-6 text-white" />
              </div>
              <div>
                <h2 className="text-2xl font-bold text-gray-900">
                  YouTube Video Dashboard
                </h2>
                <p className="text-gray-500">
                  {isSearchMode 
                    ? `Live search results for "${searchQuery}"` 
                    : `${videosData?.count ? videosData.count.toLocaleString() : 'Loading'} stored videos`
                  }
                </p>
              </div>
            </div>
            
            <div className="flex items-center space-x-4">
              {isFetching && (
                <div className="flex items-center space-x-2 text-blue-600">
                  <RefreshCw className="h-4 w-4 animate-spin" />
                  <span className="text-sm">Loading...</span>
                </div>
              )}
            </div>
          </div>

          {/* Search Type Toggle */}
          {isSearchMode && (
            <div className="flex items-center space-x-4 p-3 bg-gray-50 rounded-lg">
              <span className="text-sm font-medium text-gray-700">Search Mode:</span>
              <div className="flex space-x-2">
                <button
                  onClick={() => handleSearchTypeChange('live')}
                  className={`flex items-center space-x-2 px-3 py-1.5 rounded-lg text-sm transition-colors duration-200 ${
                    searchType === 'live' 
                      ? 'bg-red-600 text-white' 
                      : 'bg-white text-gray-700 hover:bg-gray-100'
                  }`}
                >
                  <Globe className="h-4 w-4" />
                  <span>Live YouTube Search</span>
                </button>
                <button
                  onClick={() => handleSearchTypeChange('stored')}
                  className={`flex items-center space-x-2 px-3 py-1.5 rounded-lg text-sm transition-colors duration-200 ${
                    searchType === 'stored' 
                      ? 'bg-blue-600 text-white' 
                      : 'bg-white text-gray-700 hover:bg-gray-100'
                  }`}
                >
                  <Database className="h-4 w-4" />
                  <span>Search Stored Videos</span>
                </button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Enhanced Search and Filters */}
      <div className="mb-8 space-y-6">
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <SearchBar 
            onSearch={handleSearch}
            onClear={clearSearch}
            isSearchMode={isSearchMode}
            searchQuery={searchQuery}
            searchType={searchType}
          />
        </div>
        
        {/* Filters and Results Info */}
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between space-y-4 lg:space-y-0">
          <div className="flex flex-col sm:flex-row sm:items-center space-y-2 sm:space-y-0 sm:space-x-4">
            <div className="text-sm text-gray-600">
              {videosData?.count ? (
                <>
                  Showing <span className="font-semibold">{videosData.results?.length || 0}</span>
                  {!isSearchMode && ` of ${videosData.count.toLocaleString()}`} videos
                  {isSearchMode && searchType === 'live' && (
                    <span className="text-green-600 ml-1">(Live from YouTube)</span>
                  )}
                  {isSearchMode && searchType === 'stored' && (
                    <span className="text-blue-600 ml-1">(From stored videos)</span>
                  )}
                </>
              ) : (
                'Loading videos...'
              )}
            </div>
            
            {isSearchMode && (
              <button
                onClick={clearSearch}
                className="text-sm text-red-600 hover:text-red-700 transition-colors duration-200"
              >
                Clear search
              </button>
            )}
          </div>
          
          <div className="flex items-center space-x-3">
            <span className="text-sm text-gray-500">Sort by:</span>
            <select
              value={sortBy}
              onChange={(e) => handleSortChange(e.target.value)}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
            >
              <option value="latest">Latest First</option>
              <option value="oldest">Oldest First</option>
              <option value="title">Title A-Z</option>
              <option value="channel">Channel Name</option>
            </select>
          </div>
        </div>
      </div>

      {/* Video Grid */}
      {isLoading ? (
        <LoadingSpinner message={isSearchMode ? "Searching YouTube..." : "Loading videos..."} />
      ) : (
        <VideoGrid 
          videos={videosData?.results || []} 
          isLoading={isLoading}
          isEmpty={!isLoading && videosData?.results?.length === 0}
          searchQuery={searchQuery}
          isSearchMode={isSearchMode}
          searchType={searchType}
          sortBy={sortBy}
        />
      )}

      {/* Pagination */}
      {videosData && videosData.count > pageSize && (
        <div className="mt-12">
          <Pagination
            currentPage={currentPage}
            totalItems={videosData.count}
            pageSize={pageSize}
            onPageChange={handlePageChange}
          />
        </div>
      )}
    </div>
  )
}