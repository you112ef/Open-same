# Open-Same Project Structure

This document provides a comprehensive overview of the Open-Same project structure, explaining the organization and purpose of each component.

## 📁 Root Directory Structure

```
open-same/
├── .github/                    # GitHub Actions CI/CD workflows
├── backend/                    # Go backend API service
├── frontend/                   # React frontend application
├── sdk/                       # JavaScript/TypeScript SDK
├── monitoring/                # Monitoring and observability configs
├── scripts/                   # Database and utility scripts
├── docs/                      # Comprehensive documentation
├── docker-compose.yml         # Complete development environment
├── .env.example              # Environment configuration template
├── README.md                  # Project overview and quick start
├── CONTRIBUTING.md            # Contribution guidelines
├── LICENSE                    # MIT license
└── PROJECT_STRUCTURE.md       # This file
```

## 🏗️ Backend Service (`backend/`)

The Go backend provides the core API functionality, authentication, and business logic.

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/                  # HTTP handlers and routes
│   ├── auth/                 # Authentication and authorization
│   ├── config/               # Configuration management
│   ├── database/             # Database connection and models
│   ├── middleware/           # HTTP middleware (CORS, auth, rate limiting)
│   ├── models/               # Data models and GORM definitions
│   ├── redis/                # Redis client and operations
│   ├── services/             # Business logic services
│   ├── utils/                # Utility functions
│   └── websocket/            # Real-time collaboration hub
├── pkg/                      # Public packages
├── go.mod                    # Go module dependencies
├── go.sum                    # Go module checksums
└── Dockerfile                # Multi-stage Docker build
```

### Key Backend Components

- **API Layer**: RESTful endpoints with GraphQL support
- **Authentication**: JWT-based auth with OAuth 2.0 support
- **Database**: PostgreSQL with GORM ORM
- **Caching**: Redis for session and data caching
- **Real-time**: WebSocket hub for live collaboration
- **Security**: Rate limiting, CORS, input validation

## 🎨 Frontend Application (`frontend/`)

The React frontend provides a modern, responsive user interface with PWA capabilities.

```
frontend/
├── public/                   # Static assets and PWA manifest
├── src/
│   ├── components/           # Reusable UI components
│   │   ├── ui/              # Base UI components (Radix UI)
│   │   ├── layout/          # Layout components
│   │   ├── forms/           # Form components
│   │   └── editor/          # Content editor components
│   ├── pages/               # Page components
│   ├── hooks/               # Custom React hooks
│   ├── services/            # API client and external services
│   ├── store/               # State management (Zustand)
│   ├── types/               # TypeScript type definitions
│   ├── utils/               # Utility functions
│   ├── styles/              # Global styles and Tailwind config
│   ├── App.tsx              # Main application component
│   └── main.tsx             # Application entry point
├── package.json              # Node.js dependencies
├── vite.config.ts            # Vite build configuration
├── tailwind.config.js        # Tailwind CSS configuration
└── Dockerfile                # Frontend Docker build
```

### Key Frontend Features

- **Modern UI**: Built with Radix UI and Tailwind CSS
- **PWA Support**: Offline capabilities and app-like experience
- **Real-time**: WebSocket integration for live updates
- **Responsive**: Mobile-first design approach
- **Accessibility**: WCAG 2.1 compliant components
- **Type Safety**: Full TypeScript implementation

## 🔧 SDK Package (`sdk/`)

The JavaScript/TypeScript SDK provides client libraries for easy integration.

```
sdk/
├── src/
│   ├── client/               # API client implementations
│   │   ├── OpenSameClient.ts # Main client class
│   │   ├── AuthClient.ts     # Authentication client
│   │   ├── ContentClient.ts  # Content management client
│   │   ├── UserClient.ts     # User management client
│   │   └── CollaborationClient.ts # Real-time collaboration
│   ├── types/                # TypeScript type definitions
│   │   ├── config.ts         # Configuration types
│   │   ├── models.ts         # Data model types
│   │   ├── api.ts            # API response types
│   │   ├── collaboration.ts  # Collaboration types
│   │   └── websocket.ts      # WebSocket types
│   └── index.ts              # Main export file
├── package.json              # Package configuration
├── rollup.config.js          # Rollup build configuration
└── Dockerfile                # SDK development environment
```

### SDK Capabilities

- **API Client**: HTTP client for all API endpoints
- **Real-time**: WebSocket client for live collaboration
- **Type Safety**: Full TypeScript support
- **Multiple Formats**: CommonJS and ES modules
- **Browser & Node**: Universal compatibility

## 📊 Monitoring & Observability (`monitoring/`)

Comprehensive monitoring setup for production environments.

```
monitoring/
├── prometheus.yml            # Prometheus configuration
├── grafana/                  # Grafana dashboards and datasources
│   ├── dashboards/          # Pre-configured dashboards
│   └── datasources/         # Data source configurations
├── alertmanager/             # Alert management configuration
└── rules/                    # Prometheus alerting rules
```

### Monitoring Stack

- **Metrics**: Prometheus for time-series data
- **Visualization**: Grafana dashboards
- **Alerting**: Configurable alert rules
- **Logging**: Centralized log aggregation
- **Tracing**: Distributed tracing with Jaeger

## 🐳 Docker & Infrastructure

Complete containerized development and production environment.

```
docker-compose.yml            # Multi-service development environment
├── postgres                 # PostgreSQL 15 database
├── redis                    # Redis 7 cache
├── rabbitmq                 # RabbitMQ message queue
├── kong                     # API gateway and routing
├── backend                  # Go backend service
├── frontend                 # React frontend
├── sdk                      # SDK development server
├── prometheus               # Metrics collection
├── grafana                  # Monitoring dashboards
└── adminer                  # Database administration
```

### Infrastructure Features

- **Microservices**: Service-oriented architecture
- **API Gateway**: Kong for routing and middleware
- **Message Queue**: RabbitMQ for async processing
- **Caching**: Redis for performance optimization
- **Health Checks**: Comprehensive service monitoring

## 📚 Documentation (`docs/`)

Comprehensive documentation covering all aspects of the platform.

```
docs/
├── README.md                 # Documentation overview
├── quickstart.md             # Quick start guide
├── installation.md           # Installation instructions
├── configuration.md          # Configuration guide
├── deployment.md             # Production deployment
├── architecture.md           # System architecture
├── api-reference.md          # API documentation
├── sdk.md                    # SDK usage guide
├── frontend.md               # Frontend development
├── backend.md                # Backend development
├── plugins.md                # Plugin development
├── monitoring.md             # Monitoring setup
├── troubleshooting.md        # Common issues and solutions
└── development-setup.md      # Development environment
```

## 🔄 CI/CD Pipeline (`.github/workflows/`)

Automated testing, building, and deployment pipeline.

```
.github/workflows/
├── ci.yml                    # Main CI/CD pipeline
├── security.yml              # Security scanning
├── release.yml               # Release automation
└── deploy.yml                # Deployment automation
```

### Pipeline Stages

1. **Testing**: Backend, frontend, and SDK tests
2. **Building**: Docker image creation
3. **Security**: Vulnerability scanning
4. **Deployment**: Staging and production deployment
5. **Release**: Automated release management

## 🛠️ Scripts (`scripts/`)

Utility scripts for development and deployment.

```
scripts/
├── init-db.sql              # Database initialization
├── deploy.sh                 # Deployment script
├── backup.sh                 # Database backup
└── health-check.sh           # Service health monitoring
```

## 🌐 Key Technologies & Dependencies

### Backend Stack
- **Language**: Go 1.21+
- **Framework**: Gin (HTTP), GORM (ORM)
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Message Queue**: RabbitMQ 3
- **Authentication**: JWT, OAuth 2.0

### Frontend Stack
- **Framework**: React 18, TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **UI Components**: Radix UI
- **State Management**: Zustand
- **Real-time**: WebSocket

### DevOps & Infrastructure
- **Containerization**: Docker, Docker Compose
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus, Grafana
- **API Gateway**: Kong
- **Orchestration**: Kubernetes (production)

## 🚀 Getting Started

1. **Clone the repository**
   ```bash
   git clone https://github.com/open-same/open-same.git
   cd open-same
   ```

2. **Start with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Access the platform**
   - Frontend: http://localhost:3000
   - API: http://localhost:8000
   - Admin: http://localhost:8081

4. **Development setup**
   ```bash
   # Backend
   cd backend && go mod download
   
   # Frontend
   cd frontend && npm install
   
   # SDK
   cd sdk && npm install
   ```

## 📈 Architecture Principles

- **Microservices**: Loosely coupled, independently deployable services
- **API-First**: RESTful APIs with GraphQL support
- **Real-time**: WebSocket-based live collaboration
- **Scalable**: Horizontal scaling with load balancing
- **Secure**: Comprehensive security measures
- **Observable**: Full monitoring and logging
- **Developer Experience**: Comprehensive SDKs and documentation

## 🔮 Future Enhancements

- **Plugin System**: Extensible architecture for custom functionality
- **Multi-tenancy**: Support for multiple organizations
- **Advanced Analytics**: User behavior and content insights
- **Mobile Apps**: Native iOS and Android applications
- **AI Integration**: Smart content suggestions and automation
- **Enterprise Features**: SSO, advanced permissions, compliance

---

This structure provides a solid foundation for a production-ready, scalable platform that can grow with your needs while maintaining code quality and developer experience.