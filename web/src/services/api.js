import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 20000, // Increased timeout for live searches
  headers: {
    'Content-Type': 'application/json',
  },
})

// Live YouTube search (searches YouTube directly)
export const searchYouTubeLive = async (query, page = 1, pageSize = 12, sortBy = 'latest') => {
  try {
    const params = new URLSearchParams({
      q: query,
      page: page.toString(),
      page_size: pageSize.toString(),
      sort: sortBy
    })
    const response = await api.get(`/api/videos/youtube-search?${params}`)
    return response.data
  } catch (error) {
    throw new Error(`Failed to search YouTube: ${error.response?.data?.error || error.message}`)
  }
}

// Search stored videos (searches database)
export const searchVideos = async (query, page = 1, pageSize = 12, sortBy = 'latest') => {
  try {
    const params = new URLSearchParams({
      q: query,
      page: page.toString(),
      page_size: pageSize.toString(),
      sort: sortBy
    })
    const response = await api.get(`/api/videos/search?${params}`)
    return response.data
  } catch (error) {
    throw new Error(`Failed to search stored videos: ${error.response?.data?.error || error.message}`)
  }
}

// Fetch stored videos
export const fetchVideos = async (page = 1, pageSize = 12, sortBy = 'latest') => {
  try {
    const params = new URLSearchParams({
      page: page.toString(),
      page_size: pageSize.toString(),
      sort: sortBy
    })
    const response = await api.get(`/api/videos?${params}`)
    return response.data
  } catch (error) {
    throw new Error(`Failed to fetch videos: ${error.response?.data?.error || error.message}`)
  }
}

export const getHealthStatus = async () => {
  try {
    const response = await api.get('/health')
    return response.data
  } catch (error) {
    throw new Error(`Health check failed: ${error.message}`)
  }
}

export default api