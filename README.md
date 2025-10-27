# Backend Starter with Auth 🚀

A production-ready Go backend starter template with Google OAuth authentication built in. Get your authenticated backend up and running in minutes!

## Features

✨ **Google OAuth Integration** - Pre-configured authentication flow using `goth` and `gothic`  
🔒 **Session Management** - Secure cookie-based sessions  
🌐 **CORS Support** - Ready for frontend integration  
📦 **Chi Router** - Fast and lightweight HTTP router  
🐳 **Docker Ready** - Includes Docker Compose configuration  
🧪 **Testing Setup** - Unit tests included  
📁 **Clean Architecture** - Organized project structure with separation of concerns

## Tech Stack

- **Go** - Backend language
- **Chi** - HTTP router
- **Goth/Gothic** - Multi-provider OAuth authentication
- **Gorilla Sessions** - Session management
- **PostgreSQL** - Database (optional, template ready)
- **Docker** - Containerization

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for database)
- Google OAuth credentials ([Get them here](https://console.cloud.google.com/))

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/GRACENOBLE/backend-starter-with-auth.git
   cd backend-starter-with-auth
   ```

2. **Set up environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   GOOGLE_CLIENT_ID=your_google_client_id
   GOOGLE_CLIENT_SECRET=your_google_client_secret
   ```

3. **Configure Google OAuth**
   
   In your [Google Cloud Console](https://console.cloud.google.com/):
   - Create a new project or select an existing one
   - Enable the Google+ API
   - Create OAuth 2.0 credentials
   - Add authorized redirect URI: `http://localhost:3000/auth/google/callback`

4. **Run the application**
   ```bash
   make run
   ```

   The server will start on `http://localhost:3000`

5. **Test the authentication**
   
   Navigate to `http://localhost:3000/auth/google` to initiate the OAuth flow.

### Frontend Integration

This backend works seamlessly with any frontend. Here's a React example:

```jsx
const handleLogin = () => {
  window.location.href = "http://localhost:3000/auth/google";
};

<button onClick={handleLogin}>Log in with Google</button>
```

After successful authentication, users will be redirected back with session cookies.

## Available Make Commands

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/
│   │   ├── auth.go              # Auth configuration
│   │   └── providers/
│   │       └── google/          # Google OAuth provider
│   ├── database/
│   │   └── database.go          # Database setup
│   └── server/
│       ├── routes.go            # API routes
│       └── server.go            # Server configuration
├── docker-compose.yml           # Docker configuration
├── go.mod                       # Go dependencies
├── Makefile                     # Build commands
└── README.md
```

## API Endpoints

- `GET /` - Hello World endpoint
- `GET /health` - Health check
- `GET /auth/{provider}` - Initiate OAuth flow (e.g., `/auth/google`)
- `GET /auth/{provider}/callback` - OAuth callback handler

## Customization

### Adding More OAuth Providers

This template uses `goth` which supports 40+ providers. To add more:

1. Install the provider package
2. Import it in `internal/auth/auth.go`
3. Add configuration in `NewAuth()` function

Example for GitHub:
```go
import "github.com/markbates/goth/providers/github"

// In NewAuth()
goth.UseProviders(
    google.New(...),
    github.New(githubKey, githubSecret, callbackURL),
)
```

### Updating Callback URLs

Update the callback URL in `internal/auth/auth.go`:
```go
google.New(googleClientId, googleClientSecret, "YOUR_CALLBACK_URL")
```

### Session Configuration

Modify session settings in `internal/auth/auth.go`:
```go
const (
    key    = "randomString"  // Use env variable in production
    MaxAge = 86400 * 30      // Session duration in seconds
    IsProd = false           // Set to true in production
)
```

## Environment Variables

Create a `.env` file with the following variables:

```env
# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# Add more providers as needed
# GITHUB_CLIENT_ID=your_github_client_id
# GITHUB_CLIENT_SECRET=your_github_client_secret
```

## Production Deployment

Before deploying to production:

1. **Security**
   - Set `IsProd = true` in `internal/auth/auth.go`
   - Use a strong, random session key from environment variables
   - Enable HTTPS

2. **Update Redirect URIs**
   - Update OAuth provider settings with production URLs
   - Update callback URLs in your code

3. **Environment Variables**
   - Use secure secret management (not `.env` files)
   - Set all required OAuth credentials

## Contributing

Contributions are welcome! Feel free to:
- Open issues for bugs or feature requests
- Submit pull requests
- Improve documentation

## License

This project is open source and available under the [MIT License](LICENSE).

## Acknowledgments

- [Goth](https://github.com/markbates/goth) - Multi-provider authentication
- [Chi](https://github.com/go-chi/chi) - Lightweight router
- [Gorilla Sessions](https://github.com/gorilla/sessions) - Session management

---

**Built with ❤️ by [GRACENOBLE](https://github.com/GRACENOBLE)**

Happy coding! 🎉
