# Design Document

## Overview

This design implements database migration management commands in a Makefile for a Go application that uses PostgreSQL with the golang-migrate/migrate library. The solution will create three main commands: `migrate-up`, `migrate-down`, and `migrate-create` that integrate with the existing database configuration and migration infrastructure.

## Architecture

The Makefile commands will leverage the existing migration infrastructure:
- **Existing Migration Code**: `conns/databases/migration.go` contains the `Migrate()` function that handles up migrations
- **Database Configuration**: Uses the existing configuration system in `conns/configs`
- **Migration Directory**: `resources/migrations` directory for storing migration files
- **golang-migrate Library**: Already included in go.mod for migration functionality

The design follows a command-line interface pattern where each Makefile target executes specific migration operations through the golang-migrate CLI tool or custom Go commands.

## Components and Interfaces

### 1. Makefile Targets

#### migrate-up
- **Purpose**: Apply all pending migrations
- **Implementation**: Uses golang-migrate CLI with database URL from config
- **Dependencies**: Database configuration, migration files in resources/migrations

#### migrate-down
- **Purpose**: Rollback the most recent migration
- **Implementation**: Uses golang-migrate CLI with -down flag
- **Dependencies**: Database configuration, existing migrations

#### migrate-create
- **Purpose**: Create new migration files with timestamp and name
- **Implementation**: Uses golang-migrate CLI create command
- **Parameters**: `name` parameter for migration name
- **Output**: Creates up/down migration files in resources/migrations

### 2. Database Connection

The commands will use the existing database configuration system:
- **Configuration Source**: Existing `configs.Configuration` structure
- **Connection String**: Generated using `GenerateDbString()` function
- **Database Type**: PostgreSQL (as evidenced by existing code)

### 3. Migration File Management

- **Location**: `resources/migrations` directory (SOURCE_DIR constant)
- **Naming Convention**: Timestamp prefix + provided name + .up/.down.sql
- **Template**: Basic SQL template for up/down operations

## Data Models

### Migration File Structure
```
resources/migrations/
├── YYYYMMDDHHMMSS_migration_name.up.sql
└── YYYYMMDDHHMMSS_migration_name.down.sql
```

### Configuration Integration
The commands will read database configuration from the existing system:
- Host, Port, User, Password, Database Name
- Connection pooling settings
- SSL mode configuration

## Error Handling

### Database Connection Errors
- **Detection**: Connection failures during migration attempts
- **Response**: Clear error messages with connection details (excluding sensitive info)
- **Recovery**: Suggest configuration verification steps

### Migration Execution Errors
- **Detection**: SQL syntax errors, constraint violations, missing dependencies
- **Response**: Display golang-migrate error messages with context
- **Recovery**: Provide rollback suggestions for failed up migrations

### File System Errors
- **Detection**: Permission issues, missing directories, file creation failures
- **Response**: Clear file system error messages
- **Recovery**: Suggest directory creation or permission fixes

### Parameter Validation
- **Detection**: Missing or invalid parameters (e.g., missing name for migrate-create)
- **Response**: Usage examples and parameter requirements
- **Recovery**: Show correct command syntax

## Testing Strategy

### Manual Testing Approach
Since these are Makefile commands, testing will be primarily manual:

1. **migrate-up Testing**
   - Test with no migrations (should show "no migrations to run")
   - Test with pending migrations (should apply successfully)
   - Test with database connection issues (should show clear errors)

2. **migrate-down Testing**
   - Test with migrations to rollback (should rollback successfully)
   - Test with no migrations to rollback (should show appropriate message)
   - Test rollback failure scenarios

3. **migrate-create Testing**
   - Test with valid migration name (should create files)
   - Test with missing name parameter (should show error)
   - Test with invalid characters in name (should handle gracefully)

### Integration Testing
- Verify commands work with existing database configuration
- Test migration file creation and execution flow
- Validate error handling with various failure scenarios

### Documentation Testing
- Verify help text is clear and accurate
- Test all provided examples work as documented
- Ensure error messages are helpful and actionable

## Implementation Notes

### golang-migrate CLI Integration
The design assumes golang-migrate CLI tool will be available. The Makefile will include installation instructions or automatic installation if needed.

### Environment Configuration
Commands will respect the existing environment configuration system, reading from config files or environment variables as currently implemented.

### Cross-Platform Compatibility
Makefile commands will be designed to work on Unix-like systems (Linux, macOS) where Make is commonly available.

### Security Considerations
- Database credentials will be handled through existing secure configuration methods
- Migration files will be created with appropriate file permissions
- Connection strings in error messages will exclude sensitive information