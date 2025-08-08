import VideoCard from './VideoCard'
import { Search, Video, RefreshCw, Filter } from 'lucide-react'

export default function VideoGrid({ videos, isLoading, isEmpty, searchQuery, isSearchMode, sortBy }) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {[...Array(12)].map((_, i) => (
          <div key={i} className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="aspect-video bg-gray-200 loading-shimmer"></div>
            <div className="p-4 space-y-3">
              <div className="h-4 bg-gray-200 rounded loading-shimmer"></div>
              <div className="h-4 bg-gray-200 rounded w-3/4 loading-shimmer"></div>
              <div className="h-3 bg-gray-200 rounded w-1/2 loading-shimmer"></div>
              <div className="flex justify-between items-center">
                <div className="h-3 bg-gray-200 rounded w-1/4 loading-shimmer"></div>
                <div className="h-3 bg-gray-200 rounded w-1/4 loading-shimmer"></div>
              </div>
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (isEmpty) {
    return (
      <div className="text-center py-20">
        <div className="max-w-md mx-auto">
          {isSearchMode ? (
            <>
              <div className="w-20 h-20 mx-auto mb-6 bg-gray-100 rounded-full flex items-center justify-center">
                <Search className="h-10 w-10 text-gray-400" />
              </div>
              <h3 className="text-2xl font-semibold text-gray-900 mb-3">
                No videos found
              </h3>
              <p className="text-gray-500 mb-8 leading-relaxed">
                We couldn't find any videos matching "<span className="font-medium text-gray-700">{searchQuery}</span>".
                <br />Try different keywords or browse all videos.
              </p>
            </>
          ) : (
            <>
              <div className="w-20 h-20 mx-auto mb-6 bg-blue-100 rounded-full flex items-center justify-center">
                <Video className="h-10 w-10 text-blue-500" />
              </div>
              <h3 className="text-2xl font-semibold text-gray-900 mb-3">
                Loading videos...
              </h3>
              <p className="text-gray-500 mb-8 leading-relaxed">
                Our system is fetching the latest videos from YouTube.
                <br />Videos will appear here shortly.
              </p>
            </>
          )}
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <button
              onClick={() => window.location.reload()}
              className="flex items-center justify-center space-x-2 bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors duration-200"
            >
              <RefreshCw className="h-4 w-4" />
              <span>Refresh Page</span>
            </button>
            {isSearchMode && (
              <button
                onClick={() => {
                  window.location.href = window.location.pathname
                }}
                className="flex items-center justify-center space-x-2 bg-gray-200 text-gray-700 px-6 py-3 rounded-lg hover:bg-gray-300 transition-colors duration-200"
              >
                <span>Browse All Videos</span>
              </button>
            )}
          </div>
        </div>
      </div>
    )
  }

  const getSortDisplayText = (sortBy) => {
    const sortMap = {
      'latest': 'Latest First',
      'oldest': 'Oldest First', 
      'title': 'Title A-Z',
      'channel': 'Channel Name'
    }
    return sortMap[sortBy] || 'Latest First'
  }

  return (
    <div className="space-y-6">
      {/* Results Header */}
      <div className="flex items-center justify-between">
        <div className="text-sm text-gray-600">
          Sorted by <span className="font-medium text-gray-900">{getSortDisplayText(sortBy)}</span>
        </div>
        {isSearchMode && (
          <div className="text-sm text-blue-600">
            <Search className="h-4 w-4 inline mr-1" />
            Search results
          </div>
        )}
      </div>

      {/* Video Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {videos.map((video) => (
          <VideoCard key={video.id} video={video} />
        ))}
      </div>
    </div>
  )
}