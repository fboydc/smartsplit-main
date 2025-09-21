# Go Project Structure

This Go project has been refactored for better organization and maintainability. Here's the new structure:

## Project Structure

```
go/
├── main.go                     # Application entry point
├── config/
│   └── config.go              # Configuration management
├── models/
│   └── models.go              # Data structures and types
├── client/
│   └── plaid.go               # Plaid API client wrapper
├── handlers/
│   ├── auth.go                # Authentication handlers
│   ├── plaid.go               # Plaid API handlers
│   ├── budget.go              # Budget-related handlers
│   └── dummy.go               # Dummy data handlers
├── services/
│   └── database.go            # Database operations
└── utils/                     # Utility functions (for future use)
```

## Key Components

### 1. Configuration (`config/config.go`)
- Handles environment variable loading
- Validates configuration
- Provides default values
- Uses `godotenv` for `.env` file support

### 2. Models (`models/models.go`)
- All struct definitions for data types
- Request/response structures
- Database models
- API response structures

### 3. Plaid Client (`client/plaid.go`)
- Complete wrapper around Plaid API
- Interface-based design for testability
- Organized by functionality (auth, transactions, etc.)
- Helper functions for data conversion

### 4. Handlers
- **Auth Handlers** (`handlers/auth.go`): Login, JWT validation, middleware
- **Plaid Handlers** (`handlers/plaid.go`): All Plaid API endpoints
- **Budget Handlers** (`handlers/budget.go`): Budget and category management
- **Dummy Handlers** (`handlers/dummy.go`): Test data endpoints

### 5. Services (`services/database.go`)
- Database connection management
- Query execution helpers
- Database service abstraction

## Features

### Authentication
- JWT-based authentication
- Password hashing with bcrypt
- Middleware for protected routes
- User session management

### Plaid Integration
- Complete Plaid API wrapper
- Support for all major Plaid products:
  - Transactions
  - Auth
  - Identity
  - Investments
  - Assets
  - Payment Initiation (UK/EU)
  - Transfer (ACH)
  - Signal
  - CRA Reports

### Database
- PostgreSQL integration
- Prepared statements for security
- Connection pooling support
- Error handling

## Running the Application

1. **Set up environment variables** (copy `.env.example` to `.env`):
   ```
   PLAID_CLIENT_ID=your_client_id
   PLAID_SECRET=your_secret
   PLAID_ENV=sandbox
   PLAID_PRODUCTS=transactions
   PLAID_COUNTRY_CODES=US
   APP_PORT=8000
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build and run**:
   ```bash
   go build .
   ./quickstart
   ```

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login

### Plaid
- `POST /api/create_link_token` - Create Plaid Link token
- `POST /api/set_access_token` - Exchange public token for access token
- `GET /api/accounts` - Get account information
- `GET /api/transactions` - Get transactions
- `GET /api/balance` - Get account balances
- `GET /api/identity` - Get identity information
- `GET /api/auth` - Get auth information

### Budget
- `GET /api/budget` - Get user budget
- `POST /api/save_budget` - Save budget information
- `GET /api/categories` - Get categories

### Testing
- `GET /api/dummy/transactions` - Get dummy transaction data

## Benefits of This Structure

1. **Separation of Concerns**: Each package has a specific responsibility
2. **Testability**: Interface-based design makes testing easier
3. **Maintainability**: Code is organized by functionality
4. **Scalability**: Easy to add new features and handlers
5. **Reusability**: Components can be reused across different parts of the application
6. **Error Handling**: Consistent error handling patterns
7. **Configuration Management**: Centralized configuration with validation

## Future Improvements

1. Add comprehensive unit tests
2. Implement proper logging
3. Add database migrations
4. Implement request/response validation
5. Add API documentation (Swagger)
6. Implement rate limiting
7. Add health check endpoints
8. Implement graceful shutdown

## Migration Notes

If you're migrating from the old `server.go` structure:

1. Update any direct function calls to use the new handler methods
2. Update import paths to use the new package structure
3. Configuration is now managed through the `config` package
4. Database operations should use the `services` package
5. All Plaid operations should use the `client` package wrapper