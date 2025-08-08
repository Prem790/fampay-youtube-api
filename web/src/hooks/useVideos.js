import { useQuery } from '@tanstack/react-query'
import { fetchVideos, searchVideos } from '../services/api'

export const useVideos = (page = 1, pageSize = 12) => {
  return useQuery({
    queryKey: ['videos', page, pageSize],
    queryFn: () => fetchVideos(page, pageSize),
    keepPreviousData: true,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export const useVideoSearch = (query, page = 1, pageSize = 12, enabled = true) => {
  return useQuery({
    queryKey: ['videos', 'search', query, page, pageSize],
    queryFn: () => searchVideos(query, page, pageSize),
    enabled: enabled && !!query,
    keepPreviousData: true,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export default { useVideos, useVideoSearch }