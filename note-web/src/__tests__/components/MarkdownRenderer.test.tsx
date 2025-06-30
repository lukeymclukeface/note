import { render, screen } from '@testing-library/react'
import { describe, it, expect } from '@jest/globals'
import { MarkdownRenderer } from '@/components/MarkdownRenderer'

describe('MarkdownRenderer', () => {
  it('renders plain text correctly', () => {
    const content = 'This is plain text'
    render(<MarkdownRenderer content={content} />)
    
    expect(screen.getByText(content)).toBeInTheDocument()
  })

  it('renders basic markdown elements', () => {
    const content = '# Heading\n\nSome content with **bold** text.'
    render(<MarkdownRenderer content={content} />)
    
    expect(screen.getByText('Heading')).toBeInTheDocument()
    expect(screen.getByText('bold')).toBeInTheDocument()
  })

  it('applies custom className when provided', () => {
    const content = 'Test content'
    const customClass = 'custom-test-class'
    
    const { container } = render(<MarkdownRenderer content={content} className={customClass} />)
    const wrapper = container.firstChild as HTMLElement
    
    expect(wrapper).toHaveClass(customClass)
    expect(wrapper).toHaveClass('prose', 'prose-gray', 'dark:prose-invert', 'max-w-none')
  })

  it('handles empty content gracefully', () => {
    const { container } = render(<MarkdownRenderer content="" />)
    
    // Should not crash and render the wrapper div
    expect(container.firstChild).toBeInTheDocument()
    expect(container.firstChild).toHaveClass('prose')
  })

  it('renders the component structure correctly', () => {
    const content = 'Test markdown content'
    const { container } = render(<MarkdownRenderer content={content} />)
    
    // Check that the component renders with the expected structure
    const wrapper = container.firstChild as HTMLElement
    expect(wrapper).toHaveClass('prose', 'prose-gray', 'dark:prose-invert', 'max-w-none')
    
    // Content should be present
    expect(wrapper).toHaveTextContent('Test markdown content')
  })

  it('accepts and applies additional CSS classes', () => {
    const content = 'Test content'
    const additionalClasses = 'my-custom-class another-class'
    
    const { container } = render(
      <MarkdownRenderer content={content} className={additionalClasses} />
    )
    
    const wrapper = container.firstChild as HTMLElement
    expect(wrapper).toHaveClass('my-custom-class', 'another-class')
  })
})
