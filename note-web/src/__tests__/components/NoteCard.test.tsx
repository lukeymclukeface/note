import { render, screen } from '@testing-library/react'
import { describe, it, expect } from '@jest/globals'
import NoteCard from '@/components/NoteCard'
import { mockNote, mockNoteWithoutRecording } from '../mocks/data'

describe('NoteCard', () => {
  it('renders note with all data', () => {
    render(<NoteCard note={mockNote} />)
    
    // Check title
    expect(screen.getByText('Test Note Title')).toBeInTheDocument()
    
    // Check content preview
    expect(screen.getByText(/This is a test note content/)).toBeInTheDocument()
    
    // Check summary
    expect(screen.getByText('Summary:')).toBeInTheDocument()
    expect(screen.getByText(/This is a test summary/)).toBeInTheDocument()
    
    // Check recording badge
    expect(screen.getByText('Recording')).toBeInTheDocument()
    
    // Check note badge
    expect(screen.getByText('Note')).toBeInTheDocument()
    
    // Check tags
    expect(screen.getByText('test')).toBeInTheDocument()
    expect(screen.getByText('mock')).toBeInTheDocument()
    expect(screen.getByText('unit-test')).toBeInTheDocument()
    
    // Check ID
    expect(screen.getByText('ID: 1')).toBeInTheDocument()
  })
  
  it('renders note without recording', () => {
    render(<NoteCard note={mockNoteWithoutRecording} />)
    
    // Check title
    expect(screen.getByText('Note Without Recording')).toBeInTheDocument()
    
    // Should not have recording badge
    expect(screen.queryByText('🎵 Recording')).not.toBeInTheDocument()
    
    // Should still have note badge
    expect(screen.getByText('Note')).toBeInTheDocument()
  })
  
  it('formats date correctly', () => {
    render(<NoteCard note={mockNote} />)
    
    // Check that date is formatted (exact format may vary by locale)
    expect(screen.getByText(/Dec 1, 2023/)).toBeInTheDocument()
  })
  
  it('handles empty tags gracefully', () => {
    const noteWithoutTags = { ...mockNote, tags: '' }
    render(<NoteCard note={noteWithoutTags} />)
    
    // Should still render without error
    expect(screen.getByText('Test Note Title')).toBeInTheDocument()
    
    // No tags should be displayed
    expect(screen.queryByText('test')).not.toBeInTheDocument()
  })
  
  it('truncates long content', () => {
    const longContent = 'A'.repeat(400) // More than 300 characters
    const noteWithLongContent = { ...mockNote, content: longContent }
    
    render(<NoteCard note={noteWithLongContent} />)
    
    // Should show truncated content with ellipsis
    const contentElement = screen.getByText(/A+\.\.\./)
    expect(contentElement).toBeInTheDocument()
  })
  
  it('shows updated date when different from created date', () => {
    const updatedNote = {
      ...mockNote,
      updated_at: '2023-12-02T15:00:00.000Z'
    }
    
    render(<NoteCard note={updatedNote} />)
    
    // Should show updated date
    expect(screen.getByText(/Updated:/)).toBeInTheDocument()
  })
  
  it('applies correct CSS classes for styling', () => {
    const { container } = render(<NoteCard note={mockNote} />)
    
    // Check main container has proper shadcn/ui Card classes
    const cardElement = container.firstChild as HTMLElement
    expect(cardElement).toHaveClass('rounded-lg', 'border', 'bg-card', 'text-card-foreground', 'shadow-sm')
  })
})
