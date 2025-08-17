---
description: Plan to fix schema violations in PROJECTS_USER_GUIDE.md
status: planning
priority: high
estimated_effort: 4-6 hours
dependencies: []
---

# Projects User Guide Schema Fix Plan

## Overview

The PROJECTS_USER_GUIDE.md file contains extensive schema violations that prevent it from being a reliable reference for users. This plan outlines the systematic approach to fix all schema violations and align the documentation with the actual project schema.

## Current Issues Summary

### Critical Schema Violations

1. **Non-existent blocks**: `components`, `environment`, `security`, `settings`
2. **Invalid metadata structure**: Using blocks instead of objects
3. **Missing required fields**: `version`, `author`, `email`, `url`, `facts`, `logging`
4. **Invalid field names**: `parallel_workers`, `timeout_seconds`, `log_level`
5. **Invalid data types**: Boolean values for complex configurations
6. **Invalid examples**: All examples fail schema validation

### Impact

- Users following the guide will create invalid configurations
- Validation will fail with confusing error messages
- Documentation is not trustworthy
- Poor user experience and support burden

## Fix Strategy

### Phase 1: Schema Analysis and Planning (1 hour)

#### 1.1 Review Actual Schema
- [ ] Analyze `internal/schemas/schemas/project.schema.hcl`
- [ ] Document all valid fields and their requirements
- [ ] Create reference table of correct field names and types
- [ ] Identify all required vs optional fields

#### 1.2 Inventory Violations
- [ ] Create comprehensive list of all schema violations
- [ ] Categorize violations by severity (critical, major, minor)
- [ ] Map violations to specific line numbers in the file
- [ ] Estimate effort for each category

### Phase 2: Core Structure Fixes (2 hours)

#### 2.1 Remove Non-existent Blocks
- [ ] Remove all `components` blocks (lines 67, 95, 133, 175, 207, 239, 271, 303, 335, 367, 399, 431, 463, 495, 527, 559, 591, 623, 655, 687, 719, 751, 783)
- [ ] Remove all `environment` blocks (lines 115, 147, 179, 211, 243, 275, 307, 339, 371, 403, 435, 467, 499, 531, 563, 595, 627, 659, 691, 723, 755, 787)
- [ ] Remove all `security` blocks (lines 125, 157, 189, 221, 253, 285, 317, 349, 381, 413, 445, 477, 509, 541, 573, 605, 637, 669, 701, 733, 765, 797)
- [ ] Remove all `settings` blocks (lines 159, 191, 223, 255, 287, 319, 351, 383, 415, 447, 479, 511, 543, 575, 607, 639, 671, 703, 735, 767, 799)

#### 2.2 Fix Metadata Structure
- [ ] Convert all `metadata` blocks to `metadata = {}` objects
- [ ] Update field names to match schema requirements
- [ ] Fix data types (strings instead of arrays for tags)
- [ ] Ensure all metadata examples are schema-compliant

#### 2.3 Add Required Fields
- [ ] Add `version` field to all project examples
- [ ] Add `author` field to all project examples
- [ ] Add `email` field to all project examples
- [ ] Add `url` field to all project examples
- [ ] Add proper `facts` configuration blocks
- [ ] Add proper `logging` configuration blocks

### Phase 3: Field Name and Type Fixes (1 hour)

#### 3.1 Fix Field Names
- [ ] Change `parallel_workers` to `max_parallel` (lines 161, 193, 225, 257, 289, 321, 353, 385, 417, 449, 481, 513, 545, 577, 609, 641, 673, 705, 737, 769, 801)
- [ ] Change `timeout_seconds` to `default_timeout` (lines 162, 194, 226, 258, 290, 322, 354, 386, 418, 450, 482, 514, 546, 578, 610, 642, 674, 706, 738, 770, 802)
- [ ] Change `log_level` to `level` in logging blocks (lines 163, 195, 227, 259, 291, 323, 355, 387, 419, 451, 483, 515, 547, 579, 611, 643, 675, 707, 739, 771, 803)

#### 3.2 Fix Data Types
- [ ] Ensure all boolean values are `true`/`false`
- [ ] Ensure all enum values match schema requirements
- [ ] Ensure all numeric values are within valid ranges
- [ ] Fix string formatting and escaping

### Phase 4: Example Rewrites (1.5 hours)

#### 4.1 Basic Project Configuration
- [ ] Rewrite "Basic Project Configuration" section (lines 55-73)
- [ ] Ensure schema compliance
- [ ] Add proper validation
- [ ] Test with schema validator

#### 4.2 Advanced Project Configuration
- [ ] Rewrite "Advanced Project Configuration" section (lines 75-133)
- [ ] Remove invalid blocks
- [ ] Add proper facts and logging configuration
- [ ] Ensure all examples are valid

#### 4.3 Project Types
- [ ] Rewrite "Basic Project" example (lines 135-155)
- [ ] Rewrite "Web Application Project" example (lines 157-207)
- [ ] Rewrite "Database Project" example (lines 209-271)
- [ ] Ensure all examples follow schema

#### 4.4 Environment Configuration
- [ ] Remove "Environment Configuration" section entirely
- [ ] Update references to environment configuration
- [ ] Replace with project-scoped configuration examples

#### 4.5 Security Configuration
- [ ] Remove "Security Configuration" section entirely
- [ ] Update references to security configuration
- [ ] Replace with project-scoped security examples

#### 4.6 Common Use Cases
- [ ] Rewrite "Web Application Deployment" example (lines 495-527)
- [ ] Rewrite "Infrastructure Management" example (lines 529-591)
- [ ] Rewrite "Database Management" example (lines 593-655)
- [ ] Ensure all examples are schema-compliant

### Phase 5: Validation and Testing (0.5 hours)

#### 5.1 Schema Validation
- [ ] Run all examples through schema validator
- [ ] Fix any remaining validation errors
- [ ] Ensure all examples pass validation
- [ ] Document any edge cases or limitations

#### 5.2 Documentation Review
- [ ] Review all text for accuracy
- [ ] Update any references to removed sections
- [ ] Ensure CLI commands are correct
- [ ] Verify troubleshooting section is accurate

## Corrected Schema Structure

### Basic Project Configuration (Corrected)
```hcl
project {
  name = "my-web-application"
  description = "A web application deployment project"
  version = "1.0.0"
  author = "spooky-user"
  email = "user@example.com"
  url = "https://example.com/project"
  
  run {
    default_timeout = 300
    max_parallel = 4
    dry_run_default = false
    validate_before_run = true
    backup_before_changes = false
  }
  
  facts {
    timeout = 30
    auto_collect = false
    parallel_collection = 10
    retry_attempts = 3
    retry_delay = 5
    storage_format = "memory"
    compression = true
    encryption = false
  }
  
  logging {
    level = "info"
    format = "json"
    output = "stdout"
  }
  
  metadata = {
    license = "MIT"
    tags = "web,deployment,production"
  }
}
```

### Advanced Project Configuration (Corrected)
```hcl
project {
  name = "enterprise-infrastructure"
  description = "Enterprise infrastructure management project"
  version = "2.0.0"
  author = "enterprise-team"
  email = "team@enterprise.com"
  url = "https://enterprise.com/infrastructure"
  
  run {
    default_timeout = 1800
    max_parallel = 16
    dry_run_default = false
    validate_before_run = true
    backup_before_changes = true
  }
  
  facts {
    timeout = 60
    auto_collect = true
    parallel_collection = 20
    retry_attempts = 5
    retry_delay = 10
    storage_format = "memory"
    compression = true
    encryption = true
  }
  
  logging {
    level = "debug"
    format = "structured"
    output = "file"
    file {
      path = "logs/project.log"
      permissions = "0644"
      append = true
    }
    rotation {
      enabled = true
      max_size = "100MB"
      max_age = "7d"
      max_backups = 10
      compress = true
    }
  }
  
  metadata = {
    license = "Proprietary"
    tags = "enterprise,infrastructure,security"
  }
}
```

## Success Criteria

### Functional Requirements
- [ ] All examples pass schema validation
- [ ] No non-existent blocks in documentation
- [ ] All required fields are present in examples
- [ ] All field names match schema requirements
- [ ] All data types are correct

### Quality Requirements
- [ ] Documentation is accurate and trustworthy
- [ ] Examples are clear and educational
- [ ] No confusing or misleading information
- [ ] All CLI commands are correct
- [ ] Troubleshooting section is accurate

### Technical Requirements
- [ ] All HCL examples are valid
- [ ] Schema validation passes for all examples
- [ ] No syntax errors in configuration examples
- [ ] Proper field validation rules are followed
- [ ] Enum values are within allowed ranges

## Risk Mitigation

### Risks
1. **Breaking existing user configurations**: Users may have copied invalid examples
2. **Documentation inconsistency**: Other docs may reference removed sections
3. **User confusion**: Major changes may confuse existing users

### Mitigation Strategies
1. **Clear migration guide**: Provide migration path from invalid to valid configs
2. **Cross-reference updates**: Update all related documentation
3. **Gradual rollout**: Consider deprecation warnings before removal
4. **User communication**: Clear messaging about schema compliance

## Dependencies

### Required
- Access to actual project schema (`internal/schemas/schemas/project.schema.hcl`)
- Schema validation tools
- Understanding of current user configurations

### Optional
- User feedback on current documentation issues
- Examples of real user configurations
- Performance impact analysis

## Timeline

### Week 1
- [ ] Complete Phase 1 (Schema Analysis and Planning)
- [ ] Complete Phase 2 (Core Structure Fixes)
- [ ] Initial validation of basic examples

### Week 2
- [ ] Complete Phase 3 (Field Name and Type Fixes)
- [ ] Complete Phase 4 (Example Rewrites)
- [ ] Complete Phase 5 (Validation and Testing)

### Week 3
- [ ] Review and final validation
- [ ] Update related documentation
- [ ] Create migration guide
- [ ] Final review and approval

## Post-Implementation Tasks

### Documentation Updates
- [ ] Update all cross-references to PROJECTS_USER_GUIDE.md
- [ ] Update CLI reference documentation
- [ ] Update troubleshooting guides
- [ ] Update API reference documentation

### User Communication
- [ ] Create migration guide for existing users
- [ ] Update release notes
- [ ] Communicate changes to user community
- [ ] Provide support for migration questions

### Monitoring
- [ ] Monitor for user confusion or issues
- [ ] Track validation error patterns
- [ ] Collect feedback on new examples
- [ ] Plan future improvements

## Conclusion

This plan provides a systematic approach to fixing all schema violations in the PROJECTS_USER_GUIDE.md file. The phased approach ensures that changes are made incrementally and validated at each step. The end result will be documentation that users can trust and rely on for creating valid project configurations.

The estimated effort of 4-6 hours reflects the comprehensive nature of the changes required, but the benefits in terms of user experience and reduced support burden will be significant.
