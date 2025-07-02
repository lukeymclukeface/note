import { render, screen } from '@testing-library/react'
import { describe, it, expect } from '@jest/globals'
import MeetingCard from '@/components/MeetingCard'
import { mockMeeting, mockMeetingMinimal } from '../mocks/data'

describe('MeetingCard', () => {
  it('renders meeting with all data', () => {
    render(<MeetingCard meeting={mockMeeting} />)
    
    // Check title (should be a link)
    const titleLink = screen.getByRole('link', { name: 'Weekly Team Standup' })
    expect(titleLink).toBeInTheDocument()
    expect(titleLink).toHaveAttribute('href', '/meetings/1')
    
    // Check meeting date
    expect(screen.getByText('Meeting Date:')).toBeInTheDocument()
    expect(screen.getByText('Dec 1, 2023')).toBeInTheDocument()
    
    // Check attendees
    expect(screen.getByText('Attendees:')).toBeInTheDocument()
    expect(screen.getByText('John Doe, Jane Smith, Bob Johnson')).toBeInTheDocument()
    
    // Check location
    expect(screen.getByText('Location:')).toBeInTheDocument()
    expect(screen.getByText('Conference Room A')).toBeInTheDocument()
    
    // Check summary
    expect(screen.getByText('Summary:')).toBeInTheDocument()
    expect(screen.getByText(/Team discussed current sprint progress/)).toBeInTheDocument()
    
    // Check content
    expect(screen.getByText('Content:')).toBeInTheDocument()
    expect(screen.getByText(/Discussion about project progress/)).toBeInTheDocument()
    
    // Check recording badge
    expect(screen.getByText('Recording')).toBeInTheDocument()
    
    // Check meeting badge
    expect(screen.getByText('Meeting')).toBeInTheDocument()
    
    // Check tags
    expect(screen.getByText('meeting')).toBeInTheDocument()
    expect(screen.getByText('standup')).toBeInTheDocument()
    expect(screen.getByText('team')).toBeInTheDocument()
    
    // Check created date and ID
    expect(screen.getByText(/Created:/)).toBeInTheDocument()
    expect(screen.getByText('ID: 1')).toBeInTheDocument()
  })
  
  it('renders minimal meeting without optional fields', () => {
    render(<MeetingCard meeting={mockMeetingMinimal} />)
    
    // Check title
    expect(screen.getByRole('link', { name: 'Quick Sync' })).toBeInTheDocument()
    
    // Should not show empty attendees or location
    expect(screen.queryByText('Attendees:')).not.toBeInTheDocument()
    expect(screen.queryByText('Location:')).not.toBeInTheDocument()
    expect(screen.queryByText('Meeting Date:')).not.toBeInTheDocument()
    
    // Should not have recording badge
    expect(screen.queryByText('🎵 Recording')).not.toBeInTheDocument()
    
    // Should still have meeting badge
    expect(screen.getByText('Meeting')).toBeInTheDocument()
  })
  
  it('formats dates correctly', () => {
    render(<MeetingCard meeting={mockMeeting} />)
    
    // Check meeting date format
    expect(screen.getByText('Dec 1, 2023')).toBeInTheDocument()
    
    // Check created date format (should include time)
    expect(screen.getByText(/Created:/)).toBeInTheDocument()
  })
  
  it('handles missing content and summary', () => {
    const meetingWithoutContent = {
      ...mockMeeting,
      content: '',
      summary: ''
    }
    
    render(<MeetingCard meeting={meetingWithoutContent} />)
    
    // Should not show content or summary sections
    expect(screen.queryByText('Content:')).not.toBeInTheDocument()
    expect(screen.queryByText('Summary:')).not.toBeInTheDocument()
    
    // Should still render title and other elements
    expect(screen.getByText('Weekly Team Standup')).toBeInTheDocument()
  })
  
  it('truncates long content with line-clamp', () => {
    const longSummary = 'This is a very long summary that should be truncated. '.repeat(20)
    const meetingWithLongSummary = {
      ...mockMeeting,
      summary: longSummary
    }
    
    render(<MeetingCard meeting={meetingWithLongSummary} />)
    
    // Check that summary is present but truncated
    const summaryElement = screen.getByText(/This is a very long summary/)
    expect(summaryElement).toHaveClass('line-clamp-3')
  })
  
  it('handles empty tags', () => {
    const meetingWithoutTags = { ...mockMeeting, tags: '' }
    render(<MeetingCard meeting={meetingWithoutTags} />)
    
    // Should not show any tag elements
    expect(screen.queryByText('meeting')).not.toBeInTheDocument()
    expect(screen.queryByText('standup')).not.toBeInTheDocument()
  })
  
  it('applies correct CSS classes for styling', () => {
    const { container } = render(<MeetingCard meeting={mockMeeting} />)
    
    // Check main container has proper classes
    const cardElement = container.firstChild as HTMLElement
    expect(cardElement).toHaveClass('bg-white', 'dark:bg-gray-800', 'rounded-lg', 'shadow-md')
  })
  
  it('has proper hover states for title link', () => {
    render(<MeetingCard meeting={mockMeeting} />)
    
    const titleLink = screen.getByRole('link', { name: 'Weekly Team Standup' })
    expect(titleLink).toHaveClass('hover:text-blue-600', 'dark:hover:text-blue-400')
  })
})
