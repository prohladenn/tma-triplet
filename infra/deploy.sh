#!/bin/bash

# Deployment script for TMA Triplet application
# This script helps deploy the application using Docker Compose

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    print_success "Docker is installed"
}

# Check if Docker Compose is installed
check_docker_compose() {
    if ! command -v docker compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    print_success "Docker Compose is installed"
}

# Check if .env file exists
check_env_file() {
    if [ ! -f .env ]; then
        print_error ".env file not found!"
        print_info "Creating .env from .env.example..."
        cp .env.example .env
        print_info "Please edit .env file and add your configuration"
        exit 1
    fi
    print_success ".env file found"
}

# Build images
build_images() {
    print_info "Building Docker images..."
    docker compose build --no-cache
    print_success "Images built successfully"
}

# Start services
start_services() {
    print_info "Starting services..."
    docker compose up -d
    print_success "Services started successfully"
}

# Stop services
stop_services() {
    print_info "Stopping services..."
    docker compose down
    print_success "Services stopped successfully"
}

# Restart services
restart_services() {
    print_info "Restarting services..."
    docker compose restart
    print_success "Services restarted successfully"
}

# Show logs
show_logs() {
    docker compose logs -f
}

# Show status
show_status() {
    print_info "Service Status:"
    docker compose ps
}

# Health check
health_check() {
    print_info "Checking service health..."
    
    # Check backend
    if curl -s http://localhost:3000/api/notes > /dev/null; then
        print_success "Backend is healthy"
    else
        print_error "Backend is not responding"
    fi
    
    # Check frontend
    if curl -s http://localhost/ > /dev/null; then
        print_success "Frontend is healthy"
    else
        print_error "Frontend is not responding"
    fi
}

# Clean up
cleanup() {
    print_info "Cleaning up..."
    docker compose down -v
    docker system prune -f
    print_success "Cleanup completed"
}

# Main script
main() {
    echo "======================================"
    echo "  TMA Triplet Deployment Script"
    echo "======================================"
    echo ""

    case "$1" in
        check)
            print_info "Running pre-deployment checks..."
            check_docker
            check_docker_compose
            check_env_file
            print_success "All checks passed!"
            ;;
        build)
            check_docker
            check_docker_compose
            check_env_file
            build_images
            ;;
        start)
            check_docker
            check_docker_compose
            check_env_file
            start_services
            echo ""
            print_success "Application is running!"
            print_info "Frontend: http://localhost"
            print_info "Backend API: http://localhost:3000/api"
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        logs)
            show_logs
            ;;
        status)
            show_status
            ;;
        health)
            health_check
            ;;
        deploy)
            check_docker
            check_docker_compose
            check_env_file
            build_images
            start_services
            echo ""
            print_success "Deployment completed!"
            print_info "Waiting for services to be ready..."
            sleep 5
            health_check
            ;;
        cleanup)
            cleanup
            ;;
        *)
            echo "Usage: $0 {check|build|start|stop|restart|logs|status|health|deploy|cleanup}"
            echo ""
            echo "Commands:"
            echo "  check   - Check if all dependencies are installed"
            echo "  build   - Build Docker images"
            echo "  start   - Start all services"
            echo "  stop    - Stop all services"
            echo "  restart - Restart all services"
            echo "  logs    - Show service logs (follow mode)"
            echo "  status  - Show service status"
            echo "  health  - Check service health"
            echo "  deploy  - Full deployment (build + start + health check)"
            echo "  cleanup - Stop services and clean up Docker resources"
            exit 1
            ;;
    esac
}

main "$@"
