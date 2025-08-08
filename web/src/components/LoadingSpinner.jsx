import { Loader2, Youtube } from 'lucide-react'

export default function LoadingSpinner({ message = "Loading videos..." }) {
  return (
    <div className="flex flex-col items-center justify-center py-16">
      <div className="relative">
        {/* Animated YouTube icon */}
        <Youtube className="h-16 w-16 text-red-500 animate-pulse" />
        <Loader2 className="h-6 w-6 text-blue-500 animate-spin absolute -top-1 -right-1" />
      </div>
      
      <h3 className="mt-4 text-lg font-medium text-gray-900">{message}</h3>
      <p className="mt-2 text-sm text-gray-500 max-w-md text-center">
        Fetching the latest cricket videos from YouTube API...
      </p>
      
      {/* Loading dots */}
      <div className="flex space-x-1 mt-4">
        <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce"></div>
        <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
        <div className="w-2 h-2 bg-blue-500 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
      </div>
    </div>
  )
}
