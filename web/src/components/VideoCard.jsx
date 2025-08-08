import { Play, ExternalLink, Clock, User } from 'lucide-react'
import { formatTimeAgo, truncateText, getYouTubeVideoUrl } from '../utils/formatters'

export default function VideoCard({ video }) {
  const openYouTubeVideo = () => {
    window.open(getYouTubeVideoUrl(video.video_id), '_blank')
  }

  return (
    <div className="video-card group">
      {/* Thumbnail */}
      <div 
        className="relative aspect-video bg-gray-200 rounded-t-lg overflow-hidden cursor-pointer"
        onClick={openYouTubeVideo}
      >
        <img
          src={video.thumbnails?.high || video.thumbnails?.medium || video.thumbnails?.default}
          alt={video.title}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          loading="lazy"
          onError={(e) => {
            e.target.src = 'https://via.placeholder.com/480x360/f3f4f6/9ca3af?text=Video+Thumbnail'
          }}
        />
        
        {/* Play Overlay */}
        <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-30 transition-all duration-300 flex items-center justify-center">
          <div className="transform scale-0 group-hover:scale-100 transition-transform duration-300">
            <div className="bg-red-600 rounded-full p-3 shadow-lg">
              <Play className="h-6 w-6 text-white fill-current" />
            </div>
          </div>
        </div>
        
        {/* Duration Badge (if available) */}
        <div className="absolute bottom-2 right-2 bg-black bg-opacity-75 text-white text-xs px-2 py-1 rounded">
          <Clock className="h-3 w-3 inline mr-1" />
          Video
        </div>
      </div>
      
      {/* Content */}
      <div className="p-4">
        {/* Title */}
        <h3 
          className="font-medium text-gray-900 line-clamp-2 cursor-pointer hover:text-blue-600 transition-colors duration-200 leading-tight mb-2"
          onClick={openYouTubeVideo}
          title={video.title}
        >
          {video.title}
        </h3>
        
        {/* Channel Info */}
        <div className="flex items-center space-x-2 mb-3">
          <User className="h-4 w-4 text-gray-400" />
          <span className="text-sm text-gray-600 truncate">
            {video.channel_title}
          </span>
        </div>
        
        {/* Description */}
        {video.description && (
          <p className="text-sm text-gray-500 line-clamp-2 mb-3">
            {truncateText(video.description, 120)}
          </p>
        )}
        
        {/* Footer */}
        <div className="flex items-center justify-between text-xs text-gray-500">
          <div className="flex items-center space-x-1">
            <Clock className="h-3 w-3" />
            <span>{formatTimeAgo(video.published_at)}</span>
          </div>
          
          <button
            onClick={openYouTubeVideo}
            className="flex items-center space-x-1 text-blue-600 hover:text-blue-700 transition-colors duration-200"
          >
            <ExternalLink className="h-3 w-3" />
            <span>Watch</span>
          </button>
        </div>
      </div>
    </div>
  )
}