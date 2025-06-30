import React from 'react';

// Simple mock that renders markdown content as plain text with basic HTML conversion
const ReactMarkdown = ({ children }) => {
  // Basic markdown parsing for testing
  const content = children || '';
  
  // Handle basic markdown patterns for testing
  const processContent = (text) => {
    // Replace headings
    text = text.replace(/^# (.+)$/gm, '<h1>$1</h1>');
    text = text.replace(/^## (.+)$/gm, '<h2>$1</h2>');
    text = text.replace(/^### (.+)$/gm, '<h3>$1</h3>');
    
    // Replace bold and italic
    text = text.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
    text = text.replace(/\*(.+?)\*/g, '<em>$1</em>');
    
    // Replace inline code
    text = text.replace(/`(.+?)`/g, '<code>$1</code>');
    
    // Replace links
    text = text.replace(/\[(.+?)\]\((.+?)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer">$1</a>');
    
    // Replace blockquotes
    text = text.replace(/^> (.+)$/gm, '<blockquote>$1</blockquote>');
    
    // Replace lists
    text = text.replace(/^- (.+)$/gm, '<li>$1</li>');
    text = text.replace(/^1\. (.+)$/gm, '<li>$1</li>');
    
    // Replace horizontal rules
    text = text.replace(/^---$/gm, '<hr />');
    
    // Handle strikethrough
    text = text.replace(/~~(.+?)~~/g, '<del>$1</del>');
    
    // Handle checkboxes
    text = text.replace(/- \[x\] (.+)/g, '<li>$1</li>');
    text = text.replace(/- \[ \] (.+)/g, '<li>$1</li>');
    
    return text;
  };
  
  const processedContent = processContent(content);
  
  return React.createElement('div', {
    dangerouslySetInnerHTML: { __html: processedContent }
  });
};

export default ReactMarkdown;
