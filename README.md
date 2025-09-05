# Gjallarhorn - Uptime Monitoring App

A full-stack uptime monitoring application built with Go and React. Gjallarhorn monitors your services and sends Pushover notifications when they go down.

## Features

- **Service Monitoring**: Add, edit, and delete services to monitor
- **Real-time Status**: Live status indicators (online/offline/unknown)
- **Pushover Notifications**: Get instant alerts when services go down
- **Modern UI**: Clean, responsive interface built with React and TailwindCSS
- **Single Binary**: Deploy as a single executable with embedded frontend
- **Docker Support**: Multi-architecture Docker images
- **REST API**: Full REST API for service management

## Architecture

- **Backend**: Go with Echo framework
- **Frontend**: React with Vite and TailwindCSS
- **Database**: In-memory (services stored in memory)
- **Notifications**: Pushover API integration
- **Deployment**: Single binary with embedded frontend

## Quick Start

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Docker (optional)

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd gjallarhorn
   ```

2. **Install dependencies**:
   ```bash
   # Install Go dependencies
   go mod download
   
   # Install frontend dependencies
   npm install
   ```

3. **Configure environment** (optional):
   ```bash
   cp env.example .env
   # Edit .env with your Pushover credentials
   ```

4. **Start development servers**:
   ```bash
   # Option 1: Use the development script
   ./scripts/dev.sh
   
   # Option 2: Start manually
   # Terminal 1 - Backend
   go run .
   
   # Terminal 2 - Frontend
   npm run dev
   ```

5. **Access the application**:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080/api

### Production Build

1. **Build the application**:
   ```bash
   ./scripts/build.sh
   ```

2. **Run the binary**:
   ```bash
   ./gjallarhorn
   ```

3. **Access the application**:
   - Application: http://localhost:8080

### Docker Deployment

1. **Build Docker image**:
   ```bash
   # Single architecture
   docker build -t gjallarhorn .
   
   # Multi-architecture (requires docker buildx)
   ./scripts/docker-build.sh
   ```

2. **Run with Docker Compose**:
   ```bash
   # Copy environment file
   cp env.example .env
   # Edit .env with your settings
   
   # Start the application
   docker-compose up -d
   ```

3. **Data Persistence**:
   - Services and configuration are stored in `/data` inside the container
   - Docker volume `gjallarhorn_data` is automatically created
   - Data persists across container restarts and updates
   - To backup: `docker run --rm -v gjallarhorn_data:/data -v $(pwd):/backup alpine tar czf /backup/gjallarhorn-backup.tar.gz -C /data .`

4. **Access the application**:
   - Application: http://localhost:8080

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `PUSHOVER_USER_KEY` | Your Pushover user key | - |
| `PUSHOVER_APP_TOKEN` | Your Pushover app token | - |
| `PUSHOVER_ENABLED` | Enable Pushover notifications | `false` |

### Pushover Setup

1. Sign up for a free account at [pushover.net](https://pushover.net/)
2. Find your User Key on the main page
3. Create a new application or use the default token
4. Set the environment variables or configure through the UI

## API Endpoints

### Services

- `GET /api/services` - List all services
- `POST /api/services` - Create a new service
- `PUT /api/services/:id` - Update a service
- `DELETE /api/services/:id` - Delete a service
- `GET /api/services/:id/status` - Get service status

### Notifications

- `GET /api/notifications/config` - Get notification configuration
- `POST /api/notifications/config` - Update notification configuration

### Service Object

```json
{
  "id": "uuid",
  "name": "My Website",
  "url": "https://example.com",
  "interval": 60,
  "status": "online",
  "lastChecked": "2024-01-01T12:00:00Z",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```

## Development

### Project Structure

```
gjallarhorn/
├── main.go                 # Application entry point
├── models.go              # Data models
├── monitor.go             # Service monitoring logic
├── notifications.go       # Pushover integration
├── go.mod                 # Go dependencies
├── package.json           # Frontend dependencies
├── vite.config.js         # Vite configuration
├── tailwind.config.js     # TailwindCSS configuration
├── src/                   # Frontend source code
│   ├── components/        # React components
│   ├── context/          # React context providers
│   ├── services/         # API service layer
│   └── main.jsx          # Frontend entry point
├── scripts/              # Build and development scripts
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Docker Compose configuration
└── README.md             # This file
```

### Adding New Features

1. **Backend**: Add new endpoints in `main.go` and implement logic in appropriate files
2. **Frontend**: Add new components in `src/components/` and update routing in `App.jsx`
3. **API**: Update the API service layer in `src/services/api.js`

### Building for Different Platforms

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o gjallarhorn-linux-amd64 .

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o gjallarhorn-windows-amd64.exe .

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o gjallarhorn-darwin-amd64 .

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o gjallarhorn-darwin-arm64 .
```

## Monitoring

The application includes health checks and monitoring:

- **Health Check**: `GET /api/services` (used by Docker health check)
- **Service Status**: Real-time status updates every 10 seconds
- **Error Handling**: Comprehensive error handling and logging
- **Graceful Shutdown**: Proper cleanup on application shutdown

## Troubleshooting

### Common Issues

1. **Frontend not loading**: Ensure the frontend is built and the `dist/` folder exists
2. **CORS errors**: The development server includes CORS middleware
3. **Pushover not working**: Verify your credentials and check the application logs
4. **Services not updating**: Check the console for errors and ensure the backend is running

### Logs

- **Development**: Logs are printed to stdout
- **Production**: Use your system's logging mechanism
- **Docker**: Use `docker logs <container-name>` to view logs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

## Support

For issues and questions:
1. Check the troubleshooting section
2. Search existing issues
3. Create a new issue with detailed information

---

**Gjallarhorn** - Named after the mythical horn that signals the beginning of Ragnarök, this tool signals when your services are down! 🔥
