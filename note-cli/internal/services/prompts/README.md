# OpenAI Prompts

This directory contains embedded prompt templates used by the note-cli application for OpenAI API interactions. All prompts are loaded using Go's `embed` feature and processed as templates.

## Content Analysis Prompts

### `content_analysis_system.md`
System prompt that instructs the AI to analyze content and determine its type (meeting, interview, lecture, conversation, presentation, or other) and generate an appropriate title.

### `content_analysis_user.md`
User prompt template that provides content for analysis. Uses `{{.Content}}` template variable.

## Specialized Summary Prompts

Each content type has dedicated system and user prompts for creating structured summaries:

### Meeting Prompts (`meeting_system.md`, `meeting_user.md`)
- **Structure**: Overview → Key Topics → Decisions → Action Items → Next Steps
- **Focus**: Professional meeting minutes with decisions and action items

### Interview Prompts (`interview_system.md`, `interview_user.md`)
- **Structure**: Overview → Key Q&A → Insights → Notable Quotes → Conclusion
- **Focus**: Structured interview summaries highlighting key exchanges

### Lecture Prompts (`lecture_system.md`, `lecture_user.md`)
- **Structure**: Topic Overview → Key Concepts → Main Points → Examples → Conclusion
- **Focus**: Academic content with learning objectives and concepts

### Presentation Prompts (`presentation_system.md`, `presentation_user.md`)
- **Structure**: Overview → Key Points → Data & Evidence → Conclusions → Call to Action
- **Focus**: Business presentations with supporting data and recommendations

### Conversation Prompts (`conversation_system.md`, `conversation_user.md`)
- **Structure**: Participants → Main Topics → Key Points → Agreements/Disagreements → Outcomes
- **Focus**: General conversations with participant perspectives

### General Prompts (`general_system.md`, `general_user.md`)
- **Structure**: Overview → Key Points → Details → Conclusion
- **Focus**: Fallback for unclassified or "other" content types

## Legacy Prompts

### `legacy_summary_system.md`, `legacy_summary_user.md`
Backward-compatible prompts used by the `SummarizeText` method to maintain consistency with existing functionality.

## Template Variables

All user prompt templates support the following variables:
- `{{.Content}}`: The text content to be processed

## Customization

To modify prompts:
1. Edit the relevant `.md` files
2. Rebuild the application (prompts are embedded at compile time)
3. Test with various content types to ensure proper formatting

## Usage in Code

Prompts are accessed through the `PromptService`:

```go
promptService := NewPromptService()

// For content analysis
systemPrompt, userPrompt, err := promptService.GetContentAnalysisPrompts(content)

// For specialized summaries
systemPrompt, userPrompt, err := promptService.GetSummaryPrompts(contentType, content)

// For legacy summaries
systemPrompt, userPrompt, err := promptService.GetLegacySummaryPrompts(content)
```
