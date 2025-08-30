#!/bin/bash

# Open-Same Cloudflare Pages Deployment Script
# This script builds and deploys the frontend to Cloudflare Pages

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="open-same"
FRONTEND_DIR="frontend"
BUILD_DIR="dist"
CLOUDFLARE_ACCOUNT_ID=""
CLOUDFLARE_PROJECT_NAME=""

# Functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}  Open-Same Deployment Script  ${NC}"
    echo -e "${BLUE}================================${NC}"
}

print_step() {
    echo -e "${YELLOW}[STEP]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

check_requirements() {
    print_step "Checking requirements..."
    
    # Check if Node.js is installed
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js 18+ first."
        exit 1
    fi
    
    # Check Node.js version
    NODE_VERSION=$(node -v | cut -d'v' -f2 | cut -d'.' -f1)
    if [ "$NODE_VERSION" -lt 18 ]; then
        print_error "Node.js 18+ is required. Current version: $(node -v)"
        exit 1
    fi
    
    # Check if npm is installed
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install npm first."
        exit 1
    fi
    
    # Check if Wrangler is installed
    if ! command -v wrangler &> /dev/null; then
        print_error "Wrangler CLI is not installed. Please install it first: npm install -g wrangler"
        exit 1
    fi
    
    print_success "All requirements met"
}

setup_environment() {
    print_step "Setting up environment..."
    
    # Load environment variables
    if [ -f ".env" ]; then
        export $(cat .env | grep -v '^#' | xargs)
        print_info "Loaded environment variables from .env"
    fi
    
    # Set Cloudflare configuration
    if [ -z "$CLOUDFLARE_ACCOUNT_ID" ]; then
        CLOUDFLARE_ACCOUNT_ID=${CF_PAGES_ACCOUNT_ID:-""}
    fi
    
    if [ -z "$CLOUDFLARE_PROJECT_NAME" ]; then
        CLOUDFLARE_PROJECT_NAME=${CF_PAGES_PROJECT_NAME:-"open-same"}
    fi
    
    # Validate Cloudflare configuration
    if [ -z "$CLOUDFLARE_ACCOUNT_ID" ]; then
        print_error "Cloudflare Account ID not set. Please set CF_PAGES_ACCOUNT_ID or CLOUDFLARE_ACCOUNT_ID"
        exit 1
    fi
    
    print_info "Cloudflare Account ID: $CLOUDFLARE_ACCOUNT_ID"
    print_info "Cloudflare Project Name: $CLOUDFLARE_PROJECT_NAME"
}

install_dependencies() {
    print_step "Installing dependencies..."
    
    cd "$FRONTEND_DIR"
    
    if [ -f "package-lock.json" ]; then
        npm ci
    else
        npm install
    fi
    
    cd ..
    
    print_success "Dependencies installed"
}

build_frontend() {
    print_step "Building frontend..."
    
    cd "$FRONTEND_DIR"
    
    # Set build environment variables
    export REACT_APP_API_URL=${REACT_APP_API_URL:-"https://api.opensame.com"}
    export REACT_APP_WS_URL=${REACT_APP_WS_URL:-"wss://api.opensame.com/ws"}
    export REACT_APP_ENVIRONMENT=${REACT_APP_ENVIRONMENT:-"production"}
    
    print_info "Building with API URL: $REACT_APP_API_URL"
    print_info "Building with WebSocket URL: $REACT_APP_WS_URL"
    print_info "Building for environment: $REACT_APP_ENVIRONMENT"
    
    # Build the application
    npm run build
    
    # Check if build was successful
    if [ ! -d "$BUILD_DIR" ]; then
        print_error "Build failed: $BUILD_DIR directory not found"
        exit 1
    fi
    
    cd ..
    
    print_success "Frontend built successfully"
}

deploy_to_cloudflare() {
    print_step "Deploying to Cloudflare Pages..."
    
    # Check if wrangler.toml exists
    if [ ! -f "wrangler.toml" ]; then
        print_error "wrangler.toml not found. Please create it first."
        exit 1
    fi
    
    # Deploy using Wrangler
    print_info "Deploying to Cloudflare Pages..."
    
    if wrangler pages deploy "$FRONTEND_DIR/$BUILD_DIR" --project-name="$CLOUDFLARE_PROJECT_NAME" --branch=main; then
        print_success "Deployment successful!"
        print_info "Your application is now live on Cloudflare Pages"
    else
        print_error "Deployment failed"
        exit 1
    fi
}

setup_custom_domain() {
    print_step "Setting up custom domain (if configured)..."
    
    if [ -n "$CUSTOM_DOMAIN" ]; then
        print_info "Setting up custom domain: $CUSTOM_DOMAIN"
        
        # Add custom domain to Cloudflare Pages
        if wrangler pages domain add "$CUSTOM_DOMAIN" --project-name="$CLOUDFLARE_PROJECT_NAME"; then
            print_success "Custom domain added: $CUSTOM_DOMAIN"
        else
            print_warning "Failed to add custom domain. You may need to configure it manually in the Cloudflare dashboard."
        fi
    else
        print_info "No custom domain configured"
    fi
}

cleanup() {
    print_step "Cleaning up..."
    
    # Remove build artifacts
    if [ -d "$FRONTEND_DIR/$BUILD_DIR" ]; then
        rm -rf "$FRONTEND_DIR/$BUILD_DIR"
        print_info "Removed build directory"
    fi
    
    print_success "Cleanup completed"
}

main() {
    print_header
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --account-id)
                CLOUDFLARE_ACCOUNT_ID="$2"
                shift 2
                ;;
            --project-name)
                CLOUDFLARE_PROJECT_NAME="$2"
                shift 2
                ;;
            --custom-domain)
                CUSTOM_DOMAIN="$2"
                shift 2
                ;;
            --skip-build)
                SKIP_BUILD=true
                shift
                ;;
            --help)
                echo "Usage: $0 [OPTIONS]"
                echo "Options:"
                echo "  --account-id ID       Cloudflare Account ID"
                echo "  --project-name NAME   Cloudflare Project Name"
                echo "  --custom-domain DOMAIN Custom domain to configure"
                echo "  --skip-build          Skip building (use existing build)"
                echo "  --help                Show this help message"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    # Execute deployment steps
    check_requirements
    setup_environment
    
    if [ "$SKIP_BUILD" != "true" ]; then
        install_dependencies
        build_frontend
    else
        print_info "Skipping build step"
    fi
    
    deploy_to_cloudflare
    setup_custom_domain
    cleanup
    
    print_success "Deployment completed successfully!"
    print_info "Your Open-Same application is now live on Cloudflare Pages"
}

# Handle script interruption
trap 'print_error "Deployment interrupted"; exit 1' INT TERM

# Run main function
main "$@"