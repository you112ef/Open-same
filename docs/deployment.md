# Deployment Guide

This guide covers deploying Open-Same to various platforms, with a focus on Cloudflare Pages for the frontend and cloud providers for the backend.

## 🚀 Quick Deployment Options

### Option 1: Cloudflare Pages (Frontend) + Cloudflare Workers (Backend)
- **Frontend**: Deploy to Cloudflare Pages for global CDN
- **Backend**: Use Cloudflare Workers for serverless backend
- **Database**: Cloudflare D1 (SQLite) or external PostgreSQL
- **Storage**: Cloudflare R2 for file storage

### Option 2: Vercel (Frontend) + Railway (Backend)
- **Frontend**: Deploy to Vercel for fast builds
- **Backend**: Deploy to Railway for managed infrastructure
- **Database**: Railway PostgreSQL
- **Storage**: Railway file storage or external S3

### Option 3: Self-Hosted (Full Control)
- **Frontend**: Deploy to your own server/CDN
- **Backend**: Deploy to VPS, cloud VM, or Kubernetes
- **Database**: Self-managed PostgreSQL
- **Storage**: Local storage or external S3

## 🌐 Cloudflare Pages Deployment (Recommended)

### Prerequisites

1. **Cloudflare Account**
   - Sign up at [cloudflare.com](https://cloudflare.com)
   - Get your Account ID from the dashboard

2. **Wrangler CLI**
   ```bash
   npm install -g wrangler
   wrangler login
   ```

3. **Domain Name** (optional)
   - Cloudflare provides free subdomains
   - Custom domains can be configured

### Step 1: Configure Environment

1. **Copy environment template**
   ```bash
   cp .env.example .env
   ```

2. **Edit environment variables**
   ```bash
   # Frontend Configuration
   REACT_APP_API_URL=https://your-backend-api.com
   REACT_APP_WS_URL=wss://your-backend-api.com/ws
   REACT_APP_ENVIRONMENT=production
   
   # Cloudflare Configuration
   CF_PAGES_ACCOUNT_ID=your-cloudflare-account-id
   CF_PAGES_PROJECT_NAME=open-same
   ```

### Step 2: Update Cloudflare Configuration

1. **Edit wrangler.toml**
   ```toml
   name = "open-same"
   compatibility_date = "2024-01-01"
   compatibility_flags = ["nodejs_compat"]
   
   [build]
   command = "npm run build"
   cwd = "frontend"
   
   [build.upload]
   format = "directory"
   dir = "frontend/dist"
   
   [site]
   bucket = "frontend/dist"
   
   [env.production.vars]
   ENVIRONMENT = "production"
   API_URL = "https://your-backend-api.com"
   ```

### Step 3: Deploy Frontend

1. **Use the deployment script**
   ```bash
   chmod +x scripts/deploy-cloudflare.sh
   ./scripts/deploy-cloudflare.sh --account-id YOUR_ACCOUNT_ID
   ```

2. **Or deploy manually**
   ```bash
   cd frontend
   npm install
   npm run build
   wrangler pages deploy dist --project-name=open-same
   ```

### Step 4: Configure Custom Domain (Optional)

```bash
wrangler pages domain add yourdomain.com --project-name=open-same
```

## ☁️ Backend Deployment Options

### Option A: Cloudflare Workers (Serverless)

1. **Create worker script**
   ```typescript
   // workers/api.ts
   export default {
     async fetch(request: Request): Promise<Response> {
       // Your API logic here
       return new Response('Hello from Open-Same!')
     }
   }
   ```

2. **Deploy worker**
   ```bash
   wrangler deploy
   ```

### Option B: Railway (Managed Infrastructure)

1. **Install Railway CLI**
   ```bash
   npm install -g @railway/cli
   railway login
   ```

2. **Deploy backend**
   ```bash
   cd backend
   railway init
   railway up
   ```

3. **Configure environment**
   ```bash
   railway variables set DATABASE_URL=your-postgres-url
   railway variables set REDIS_URL=your-redis-url
   railway variables set JWT_SECRET=your-jwt-secret
   ```

### Option C: DigitalOcean App Platform

1. **Create app specification**
   ```yaml
   # .do/app.yaml
   name: open-same-backend
   services:
   - name: api
     source_dir: /backend
     github:
       repo: your-org/open-same
       branch: main
     run_command: go run cmd/server/main.go
     environment_slug: go
     instance_count: 1
     instance_size_slug: basic-xxs
   ```

2. **Deploy using doctl**
   ```bash
   doctl apps create --spec .do/app.yaml
   ```

### Option D: Kubernetes (Production)

1. **Create deployment manifests**
   ```yaml
   # k8s/deployment.yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: open-same-backend
   spec:
     replicas: 3
     selector:
       matchLabels:
         app: open-same-backend
     template:
       metadata:
         labels:
           app: open-same-backend
       spec:
         containers:
         - name: backend
           image: your-registry/open-same-backend:latest
           ports:
           - containerPort: 8080
           env:
           - name: DATABASE_URL
             valueFrom:
               secretKeyRef:
                 name: db-secret
                 key: url
   ```

2. **Apply manifests**
   ```bash
   kubectl apply -f k8s/
   ```

## 🗄️ Database Setup

### PostgreSQL (Recommended)

1. **Cloudflare D1 (SQLite)**
   ```bash
   # Create database
   wrangler d1 create open-same-db
   
   # Run migrations
   wrangler d1 execute open-same-db --file=./backend/migrations/001_init.sql
   ```

2. **External PostgreSQL**
   - Use managed services: Supabase, Neon, Railway
   - Or self-hosted: DigitalOcean, AWS RDS, Google Cloud SQL

3. **Run migrations**
   ```bash
   cd backend
   go run cmd/migrate/main.go
   ```

### Redis Setup

1. **Cloudflare KV (Key-Value)**
   ```bash
   wrangler kv:namespace create "OPEN_SAME_CACHE"
   ```

2. **External Redis**
   - Managed: Redis Cloud, Upstash
   - Self-hosted: DigitalOcean, AWS ElastiCache

## 🔐 Environment Configuration

### Required Environment Variables

```bash
# Application
ENVIRONMENT=production
VERSION=1.0.0

# Server
API_PORT=8080
API_HOST=0.0.0.0

# Database
DB_HOST=your-db-host
DB_PORT=5432
DB_NAME=opensame
DB_USER=opensame
DB_PASSWORD=your-secure-password
DB_SSLMODE=require

# Redis
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRATION_HOURS=24
JWT_REFRESH_HOURS=168

# AI Services
OPENAI_API_KEY=your-openai-key
OPENAI_MODEL=gpt-4
ANTHROPIC_API_KEY=your-anthropic-key
ANTHROPIC_MODEL=claude-3-sonnet-20240229

# Rate Limiting
RATE_LIMIT=100.0
```

### Security Considerations

1. **Secrets Management**
   - Use environment variables for sensitive data
   - Never commit secrets to version control
   - Use secret management services in production

2. **SSL/TLS**
   - Always use HTTPS in production
   - Configure proper SSL certificates
   - Use HSTS headers

3. **CORS Configuration**
   - Restrict origins to your domains
   - Configure proper CORS policies
   - Test CORS in different environments

## 📊 Monitoring & Observability

### Health Checks

1. **Create health check endpoint**
   ```go
   // backend/internal/api/health.go
   func HealthCheck(c *gin.Context) {
       c.JSON(http.StatusOK, gin.H{
           "status": "healthy",
           "timestamp": time.Now().ISO8601(),
           "version": "1.0.0",
       })
   }
   ```

2. **Configure monitoring**
   - Set up uptime monitoring (UptimeRobot, Pingdom)
   - Configure alerting for downtime
   - Monitor response times and error rates

### Logging

1. **Structured logging**
   ```go
   log.WithFields(log.Fields{
       "user_id": userID,
       "action": "content_create",
       "content_type": contentType,
   }).Info("Content created successfully")
   ```

2. **Log aggregation**
   - Use services like DataDog, New Relic, or ELK stack
   - Configure log retention policies
   - Set up log-based alerting

## 🚀 CI/CD Pipeline

### GitHub Actions

1. **Create workflow file**
   ```yaml
   # .github/workflows/deploy.yml
   name: Deploy to Cloudflare Pages
   
   on:
     push:
       branches: [main]
   
   jobs:
     deploy:
       runs-on: ubuntu-latest
       steps:
       - uses: actions/checkout@v3
       
       - name: Setup Node.js
         uses: actions/setup-node@v3
         with:
           node-version: '18'
           
       - name: Install dependencies
         run: |
           cd frontend
           npm ci
           
       - name: Build
         run: |
           cd frontend
           npm run build
           
       - name: Deploy to Cloudflare Pages
         uses: cloudflare/pages-action@v1
         with:
           apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
           accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
           projectName: open-same
           directory: frontend/dist
           gitHubToken: ${{ secrets.GITHUB_TOKEN }}
   ```

2. **Configure secrets**
   - `CLOUDFLARE_API_TOKEN`: Your Cloudflare API token
   - `CLOUDFLARE_ACCOUNT_ID`: Your Cloudflare account ID

### Automated Testing

1. **Run tests before deployment**
   ```yaml
   - name: Run tests
     run: |
       cd frontend
       npm run test
       cd ../backend
       go test ./...
   ```

2. **Security scanning**
   ```yaml
   - name: Security scan
     uses: snyk/actions/node@master
     with:
       args: --severity-threshold=high
   ```

## 🔄 Rollback Strategy

### Cloudflare Pages Rollback

1. **Automatic rollback**
   - Cloudflare Pages automatically detects failed deployments
   - Previous version remains active during deployment

2. **Manual rollback**
   ```bash
   # List deployments
   wrangler pages deployment list --project-name=open-same
   
   # Rollback to specific deployment
   wrangler pages deployment rollback DEPLOYMENT_ID --project-name=open-same
   ```

### Database Rollback

1. **Migration rollback**
   ```bash
   cd backend
   go run cmd/migrate/main.go --rollback
   ```

2. **Backup restoration**
   - Maintain regular database backups
   - Test backup restoration procedures
   - Document rollback procedures

## 📈 Performance Optimization

### Frontend Optimization

1. **Code splitting**
   ```typescript
   // Lazy load components
   const Editor = lazy(() => import('./components/Editor'))
   ```

2. **Bundle optimization**
   - Use dynamic imports
   - Implement tree shaking
   - Optimize images and assets

### Backend Optimization

1. **Database optimization**
   - Add proper indexes
   - Use connection pooling
   - Implement query caching

2. **Caching strategy**
   - Redis for session storage
   - CDN for static assets
   - Browser caching headers

## 🆘 Troubleshooting

### Common Issues

1. **Build failures**
   - Check Node.js version compatibility
   - Verify all dependencies are installed
   - Check for TypeScript compilation errors

2. **Deployment failures**
   - Verify Cloudflare credentials
   - Check account permissions
   - Review build output for errors

3. **Runtime errors**
   - Check environment variables
   - Verify database connectivity
   - Review application logs

### Debug Commands

```bash
# Check Cloudflare status
wrangler whoami

# Test database connection
cd backend && go run cmd/test-db/main.go

# Check Redis connection
cd backend && go run cmd/test-redis/main.go

# Verify environment
cd frontend && npm run build:check
```

## 📚 Additional Resources

- [Cloudflare Pages Documentation](https://developers.cloudflare.com/pages/)
- [Wrangler CLI Reference](https://developers.cloudflare.com/workers/wrangler/)
- [Open-Same Architecture Guide](docs/architecture.md)
- [API Reference](docs/api.md)
- [Contributing Guidelines](CONTRIBUTING.md)

---

**Need help with deployment?** Join our [Discord community](https://discord.gg/open-same) or open an issue on GitHub!