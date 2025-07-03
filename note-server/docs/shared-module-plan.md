# Note Shared Module Plan

## Overview
This document outlines the plan for creating a separate `note-shared` module that can be used across multiple note application components.

## Future Structure

```
note-shared/
├── go.mod
├── go.sum
├── README.md
├── pkg/
│   ├── models/           # Shared data models
│   │   ├── note.go
│   │   ├── user.go
│   │   └── session.go
│   ├── events/           # Event definitions for pub/sub
│   │   ├── note_events.go
│   │   └── user_events.go
│   ├── validation/       # Input validation utilities
│   │   ├── note_validator.go
│   │   └── user_validator.go
│   ├── constants/        # Shared constants
│   │   ├── status.go
│   │   └── errors.go
│   └── proto/           # Protocol buffer definitions (if using gRPC)
│       ├── note.proto
│       └── user.proto
├── internal/
│   └── testutils/       # Test utilities for consumers
└── examples/           # Usage examples
    └── basic/
```

## Migration Strategy

### Phase 1: Identify Shared Code
- [ ] Audit current `pkg/` directories in all note services
- [ ] Identify common models, validation logic, constants
- [ ] Document dependencies between services

### Phase 2: Create Shared Module
- [ ] Create new `note-shared` repository
- [ ] Move identified shared code to appropriate packages
- [ ] Set up CI/CD for the shared module
- [ ] Create comprehensive tests and documentation

### Phase 3: Update Services
- [ ] Update `note-server` to use `note-shared`
- [ ] Update `note-web` backend components (if any)
- [ ] Update any other services that emerge

### Phase 4: Versioning Strategy
- [ ] Implement semantic versioning for `note-shared`
- [ ] Set up automated release process
- [ ] Document breaking change policy

## Benefits

1. **Code Reuse**: Avoid duplicating models and utilities across services
2. **Consistency**: Ensure all services use the same data structures
3. **Maintainability**: Single source of truth for shared logic
4. **Type Safety**: Shared interfaces prevent integration errors
5. **Testing**: Centralized test utilities

## When to Create

Create the shared module when:
- [ ] At least 2 services need shared code
- [ ] Shared code is stable and well-defined
- [ ] Team has bandwidth to maintain separate module
- [ ] Benefits outweigh the added complexity

## Current Shared Code Candidates

From the current codebase, these could be moved to `note-shared`:

1. **pkg/response/**: HTTP response utilities
2. **pkg/timeutil/**: Time formatting utilities
3. **Future models**: Note, User, Session data structures
4. **Future validation**: Input validation logic
5. **Future constants**: Status codes, error messages

## Usage Example

```go
// In note-server
import (
    "github.com/your-org/note-shared/pkg/models"
    "github.com/your-org/note-shared/pkg/response"
)

func handleCreateNote(w http.ResponseWriter, r *http.Request) {
    var note models.Note
    // ... handle request
    response.WriteJSONSuccess(w, note)
}
```

## Considerations

- **Versioning**: Use semantic versioning carefully
- **Breaking Changes**: Plan migration strategy for breaking changes
- **Documentation**: Maintain clear usage documentation
- **Testing**: Comprehensive test coverage
- **Dependencies**: Minimize external dependencies in shared module
