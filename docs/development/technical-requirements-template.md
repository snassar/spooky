# Technical Requirements Template

## What We're Building
[One paragraph description of the feature/system]

## User Stories
- As a [user], I want [feature] so that [benefit]
- As a [user], I want [feature] so that [benefit]

## Technical Approach
- [High-level architecture decision]
- [Key technology choices]
- [Integration points]

## Functional Requirements
- [ ] [Specific thing the system must do]
- [ ] [Another specific thing]
- [ ] [Another specific thing]

## Non-Functional Requirements
- **Performance**: [Response time, throughput]
- **Security**: [Authentication, data protection]
- **Reliability**: [Uptime, error handling]

## API Design
```
POST /api/resource
{
  "field": "value"
}
```

## Data Model
```sql
CREATE TABLE resource (
  id UUID PRIMARY KEY,
  field VARCHAR(255)
);
```

## Testing Strategy
- [ ] Unit tests for core logic
- [ ] Integration tests for APIs
- [ ] Manual testing checklist

## Deployment
- [Where it runs]
- [How it gets deployed]
- [Environment requirements]

## Open Questions
- [ ] [Decision we need to make]
- [ ] [Risk we need to assess]
- [ ] [Alternative we need to evaluate]
