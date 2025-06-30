import { render, screen, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, jest, beforeEach } from '@jest/globals'
import { ThemeSelector } from '@/components/ThemeSelector'
import { ThemeProvider } from '@/providers/ThemeProvider'

// Create a wrapper component for tests
const ThemeSelectorWrapper = () => (
  <ThemeProvider>
    <ThemeSelector />
  </ThemeProvider>
)

describe('ThemeSelector', () => {
  it('renders theme selector button', () => {
    render(<ThemeSelectorWrapper />)
    
    // Should have a button that's accessible
    const button = screen.getByRole('button')
    expect(button).toBeInTheDocument()
    expect(button).toHaveAttribute('aria-expanded', 'false')
    expect(button).toHaveAttribute('aria-haspopup', 'true')
  })

  it('opens dropdown when clicked', async () => {
    const user = userEvent.setup()
    render(<ThemeSelectorWrapper />)
    
    const button = screen.getByRole('button')
    await user.click(button)
    
    // Should show all theme options
    expect(screen.getByRole('menuitem', { name: /Light/ })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /Dark/ })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /System/ })).toBeInTheDocument()
    
    // Should show icons for each option in the dropdown
    expect(screen.getByText('☀️')).toBeInTheDocument()
    expect(screen.getByText('🌙')).toBeInTheDocument()
    expect(screen.getAllByText('💻')).toHaveLength(2) // One in button, one in dropdown
  })

  it('shows checkmark for current theme', async () => {
    const user = userEvent.setup()
    render(<ThemeSelectorWrapper />)
    
    const button = screen.getByRole('button')
    await user.click(button)
    
    // One of the theme options should have a checkmark (the current one)
    const menuItems = screen.getAllByRole('menuitem')
    const checkedItem = menuItems.find(item => item.querySelector('svg'))
    expect(checkedItem).toBeInTheDocument()
  })

  it('changes theme when option is selected', async () => {
    const user = userEvent.setup()
    render(<ThemeSelectorWrapper />)
    
    const button = screen.getByRole('button')
    await user.click(button)
    
    // Click on Light theme
    const lightOption = screen.getByRole('menuitem', { name: /Light/ })
    await user.click(lightOption)
    
    // The theme should change (we can't easily test the actual theme change without more complex mocking)
    // But we can verify the dropdown closes
    expect(screen.queryByRole('menuitem', { name: /System/ })).not.toBeInTheDocument()
  })

  it('closes dropdown after selecting theme', async () => {
    const user = userEvent.setup()
    render(<ThemeSelectorWrapper />)
    
    const button = screen.getByRole('button')
    await user.click(button)
    
    // Select dark theme
    const darkOption = screen.getByRole('menuitem', { name: /Dark/ })
    await user.click(darkOption)
    
    // Dropdown should close (options should not be visible)
    expect(screen.queryByRole('menuitem', { name: /Light/ })).not.toBeInTheDocument()
  })

  it('closes dropdown when clicking outside', async () => {
    const user = userEvent.setup()
    render(
      <div>
        <ThemeSelectorWrapper />
        <div data-testid="outside">Outside element</div>
      </div>
    )
    
    const button = screen.getByRole('button')
    await user.click(button)
    
    // Dropdown should be open
    expect(screen.getByRole('menuitem', { name: /Light/ })).toBeInTheDocument()
    
    // Click outside
    const outsideElement = screen.getByTestId('outside')
    fireEvent.mouseDown(outsideElement)
    
    // Dropdown should close
    expect(screen.queryByRole('menuitem', { name: /Light/ })).not.toBeInTheDocument()
  })

  it('displays theme correctly', () => {
    render(<ThemeSelectorWrapper />)
    
    // Should show some theme (could be any of the three)
    const button = screen.getByRole('button')
    expect(button).toBeInTheDocument()
    
    // Should contain one of the theme icons
    const hasThemeIcon = 
      screen.queryByText('💻') || 
      screen.queryByText('☀️') || 
      screen.queryByText('🌙')
    expect(hasThemeIcon).toBeInTheDocument()
  })

  it('has proper ARIA attributes', () => {
    render(<ThemeSelectorWrapper />)
    
    const button = screen.getByRole('button')
    expect(button).toHaveAttribute('aria-expanded', 'false')
    expect(button).toHaveAttribute('aria-haspopup', 'true')
  })

  it('applies correct CSS classes', () => {
    render(<ThemeSelectorWrapper />)
    
    const button = screen.getByRole('button')
    expect(button).toHaveClass(
      'inline-flex',
      'items-center',
      'border',
      'border-gray-300',
      'dark:border-gray-600',
      'rounded-md'
    )
  })

  it('hides label on small screens', () => {
    render(<ThemeSelectorWrapper />)
    
    // Find the span with the theme label (could be any theme)
    const hiddenSpan = screen.getByRole('button').querySelector('span.hidden.sm\\:inline')
    expect(hiddenSpan).toBeInTheDocument()
    expect(hiddenSpan).toHaveClass('hidden', 'sm:inline')
  })
})
