import { formatDistanceToNow, format, parseISO } from 'date-fns'

export const formatTimeAgo = (dateString) => {
  try {
    const date = parseISO(dateString)
    return formatDistanceToNow(date, { addSuffix: true })
  } catch (error) {
    return 'Unknown time'
  }
}

export const formatFullDate = (dateString) => {
  try {
    const date = parseISO(dateString)
    return format(date, 'PPP at p')
  } catch (error) {
    return 'Unknown date'
  }
}

export const truncateText = (text, maxLength = 100) => {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

export const getYouTubeVideoUrl = (videoId) => {
  return `https://www.youtube.com/watch?v=${videoId}`
}

export const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

export default {
  formatTimeAgo,
  formatFullDate,
  truncateText,
  getYouTubeVideoUrl,
  formatNumber,
}
