#!/bin/bash
#
# deploy.sh - Deploy Unfollow Tracker to production server
# Usage: ./deploy.sh [user@]hostname[:port]
#
# Environment variables:
#   DEPLOY_USER     - SSH user (default: root)
#   DEPLOY_HOST     - SSH host (required)
#   DEPLOY_PORT     - SSH port (default: 22)
#   DEPLOY_KEY      - SSH private key path (optional)
#   SKIP_BUILD      - Skip building binaries (set to "true")
#   SKIP_MIGRATIONS - Skip running migrations (set to "true")
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="unfollow-tracker"
LOCAL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REMOTE_DIR="/opt/unfollow-tracker"
BACKUP_DIR="/opt/unfollow-tracker-backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Parse arguments
SSH_TARGET="${1:-}"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Error handler
error_exit() {
    log_error "$1"
    exit 1
}

# Parse SSH target
parse_ssh_target() {
    if [[ -z "$SSH_TARGET" ]]; then
        if [[ -n "${DEPLOY_HOST:-}" ]]; then
            SSH_USER="${DEPLOY_USER:-root}"
            SSH_HOST="$DEPLOY_HOST"
            SSH_PORT="${DEPLOY_PORT:-22}"
        else
            error_exit "Usage: $0 [user@]hostname[:port] or set DEPLOY_HOST environment variable"
        fi
    else
        # Parse [user@]host[:port]
        if [[ "$SSH_TARGET" =~ ^([^@]+)@(.+)$ ]]; then
            SSH_USER="${BASH_REMATCH[1]}"
            SSH_HOST_PORT="${BASH_REMATCH[2]}"
        else
            SSH_USER="${DEPLOY_USER:-root}"
            SSH_HOST_PORT="$SSH_TARGET"
        fi

        if [[ "$SSH_HOST_PORT" =~ ^([^:]+):([0-9]+)$ ]]; then
            SSH_HOST="${BASH_REMATCH[1]}"
            SSH_PORT="${BASH_REMATCH[2]}"
        else
            SSH_HOST="$SSH_HOST_PORT"
            SSH_PORT="${DEPLOY_PORT:-22}"
        fi
    fi

    log_info "Deploy target: $SSH_USER@$SSH_HOST:$SSH_PORT"
}

# Build SSH options
build_ssh_opts() {
    SSH_OPTS="-o StrictHostKeyChecking=no -o ConnectTimeout=10 -p $SSH_PORT"
    SCP_OPTS="-o StrictHostKeyChecking=no -o ConnectTimeout=10 -P $SSH_PORT"

    if [[ -n "${DEPLOY_KEY:-}" ]]; then
        SSH_OPTS="$SSH_OPTS -i $DEPLOY_KEY"
        SCP_OPTS="$SCP_OPTS -i $DEPLOY_KEY"
    fi
}

# Test SSH connection
test_ssh() {
    log_info "Testing SSH connection..."
    if ! ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "echo 'SSH connection successful'" > /dev/null 2&amp;1; then
        error_exit "Cannot connect to $SSH_USER@$SSH_HOST:$SSH_PORT"
    fi
    log_success "SSH connection OK"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check Go installation
    if ! command -v go &> /dev/null; then
        error_exit "Go is not installed. Please install Go 1.23 or later."
    fi

    GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
    log_info "Go version: $GO_VERSION"

    # Check Node.js for web build
    if ! command -v npm &> /dev/null; then
        log_warning "npm is not installed. Web frontend will not be built."
    fi

    # Check local directories
    if [[ ! -d "$LOCAL_DIR/cmd/api" ]]; then
        error_exit "Cannot find api command directory. Are you in the right repository?"
    fi

    log_success "Prerequisites OK"
}

# Build Go binaries
build_binaries() {
    if [[ "${SKIP_BUILD:-}" == "true" ]]; then
        log_info "Skipping build (SKIP_BUILD=true)"
        return
    fi

    log_info "Building Go binaries for Linux AMD64..."

    cd "$LOCAL_DIR"

    # Create build directory
    mkdir -p "$LOCAL_DIR/build"

    # Build flags for production
    BUILD_FLAGS="-ldflags='-s -w' -trimpath"

    # Build API
    log_info "Building api binary..."
    GOOS=linux GOARCH=amd64 go build \
        -ldflags='-s -w' \
        -trimpath \
        -o "$LOCAL_DIR/build/api" \
        "$LOCAL_DIR/cmd/api" \
        || error_exit "Failed to build api binary"

    # Build Fetcher
    log_info "Building fetcher binary..."
    GOOS=linux GOARCH=amd64 go build \
        -ldflags='-s -w' \
        -trimpath \
        -o "$LOCAL_DIR/build/fetcher" \
        "$LOCAL_DIR/cmd/fetcher" \
        || error_exit "Failed to build fetcher binary"

    # Build Scheduler
    log_info "Building scheduler binary..."
    GOOS=linux GOARCH=amd64 go build \
        -ldflags='-s -w' \
        -trimpath \
        -o "$LOCAL_DIR/build/scheduler" \
        "$LOCAL_DIR/cmd/scheduler" \
        || error_exit "Failed to build scheduler binary"

    log_success "Binaries built successfully"

    # Build web frontend
    if command -v npm &> /dev/null && [[ -d "$LOCAL_DIR/web" ]]; then
        log_info "Building web frontend..."
        cd "$LOCAL_DIR/web"
        npm ci 2>/dev/null || npm install
        npm run build
        log_success "Web frontend built"
    fi
}

# Create remote backup
create_backup() {
    log_info "Creating remote backup..."

    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        if [[ -d $REMOTE_DIR/bin ]]; then
            mkdir -p $BACKUP_DIR
            cp -r $REMOTE_DIR $BACKUP_DIR/$APP_NAME-$TIMESTAMP
            echo 'Backup created at $BACKUP_DIR/$APP_NAME-$TIMESTAMP'
            # Keep only last 5 backups
            ls -t $BACKUP_DIR | tail -n +6 | xargs -I {} rm -rf $BACKUP_DIR/{}
        fi
    " || log_warning "Backup creation failed, continuing..."
}

# Deploy files to server
deploy_files() {
    log_info "Deploying files to server..."

    # Create temporary directory on remote
    REMOTE_TEMP="/tmp/deploy-$TIMESTAMP"
    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "mkdir -p $REMOTE_TEMP"

    # Copy binaries
    log_info "Copying binaries..."
    scp $SCP_OPTS \
        "$LOCAL_DIR/build/api" \
        "$LOCAL_DIR/build/fetcher" \
        "$LOCAL_DIR/build/scheduler" \
        "$SSH_USER@$SSH_HOST:$REMOTE_TEMP/" \
        || error_exit "Failed to copy binaries"

    # Copy migrations
    if [[ -d "$LOCAL_DIR/migrations" ]]; then
        log_info "Copying migrations..."
        scp $SCP_OPTS -r \
            "$LOCAL_DIR/migrations" \
            "$SSH_USER@$SSH_HOST:$REMOTE_TEMP/" \
            || error_exit "Failed to copy migrations"
    fi

    # Copy web dist
    if [[ -d "$LOCAL_DIR/web/dist" ]]; then
        log_info "Copying web dist..."
        # Remove old dist first
        ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "rm -rf $REMOTE_DIR/web/dist" 2>/dev/null || true
        scp $SCP_OPTS -r \
            "$LOCAL_DIR/web/dist" \
            "$SSH_USER@$SSH_HOST:$REMOTE_TEMP/web/" \
            || error_exit "Failed to copy web dist"
    fi

    # Copy .env file if it exists
    if [[ -f "$LOCAL_DIR/.env" ]]; then
        log_info "Copying .env file..."
        scp $SCP_OPTS \
            "$LOCAL_DIR/.env" \
            "$SSH_USER@$SSH_HOST:$REMOTE_TEMP/" \
            || error_exit "Failed to copy .env file"
    else
        log_warning "No .env file found locally. Make sure it exists on the server."
    fi

    # Move files to final location
    log_info "Moving files to $REMOTE_DIR..."
    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        mkdir -p $REMOTE_DIR/bin
        mkdir -p $REMOTE_DIR/migrations
        mkdir -p $REMOTE_DIR/web
        mkdir -p $REMOTE_DIR/logs

        # Move binaries
        mv $REMOTE_TEMP/api $REMOTE_DIR/bin/
        mv $REMOTE_TEMP/fetcher $REMOTE_DIR/bin/
        mv $REMOTE_TEMP/scheduler $REMOTE_DIR/bin/
        chmod +x $REMOTE_DIR/bin/*

        # Move migrations
        if [[ -d $REMOTE_TEMP/migrations ]]; then
            rm -rf $REMOTE_DIR/migrations
            mv $REMOTE_TEMP/migrations $REMOTE_DIR/
        fi

        # Move web dist
        if [[ -d $REMOTE_TEMP/web/dist ]]; then
            rm -rf $REMOTE_DIR/web/dist
            mv $REMOTE_TEMP/web/dist $REMOTE_DIR/web/
        fi

        # Move .env
        if [[ -f $REMOTE_TEMP/.env ]]; then
            mv $REMOTE_TEMP/.env $REMOTE_DIR/
            chmod 600 $REMOTE_DIR/.env
        fi

        # Set ownership
        chown -R unfollowtracker:unfollowtracker $REMOTE_DIR

        # Cleanup temp
        rm -rf $REMOTE_TEMP

        echo 'Files deployed successfully'
    " || error_exit "Failed to move files on server"

    log_success "Files deployed"
}

# Run database migrations
run_migrations() {
    if [[ "${SKIP_MIGRATIONS:-}" == "true" ]]; then
        log_info "Skipping migrations (SKIP_MIGRATIONS=true)"
        return
    fi

    log_info "Running database migrations..."

    # Check if golang-migrate is installed on server
    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        if ! command -v migrate &> /dev/null; then
            echo 'Installing golang-migrate...'
            curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz -C /tmp
            mv /tmp/migrate /usr/local/bin/migrate
            chmod +x /usr/local/bin/migrate
        fi
    " || log_warning "Could not ensure migrate is installed"

    # Run migrations
    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        cd $REMOTE_DIR
        # Load environment variables
        if [[ -f .env ]]; then
            export \$(grep -v '^#' .env | xargs)
        fi

        if [[ -z \"\${DATABASE_URL:-}\" ]]; then
            echo 'DATABASE_URL not set, skipping migrations'
            exit 0
        fi

        echo 'Running migrations...'
        migrate -path $REMOTE_DIR/migrations -database \"\$DATABASE_URL\" up
    " || error_exit "Migration failed"

    log_success "Migrations completed"
}

# Restart services
restart_services() {
    log_info "Restarting services..."

    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        # Stop services
        systemctl stop unfollow-tracker-api 2>/dev/null || true
        systemctl stop unfollow-tracker-fetcher 2>/dev/null || true
        systemctl stop unfollow-tracker-scheduler 2>/dev/null || true

        # Small delay to ensure clean shutdown
        sleep 2

        # Start services
        systemctl start unfollow-tracker-api
        systemctl start unfollow-tracker-fetcher
        systemctl start unfollow-tracker-scheduler

        # Enable services to start on boot
        systemctl enable unfollow-tracker-api
        systemctl enable unfollow-tracker-fetcher
        systemctl enable unfollow-tracker-scheduler

        echo 'Services restarted'
    " || error_exit "Failed to restart services"

    log_success "Services restarted"
}

# Verify deployment
verify_deployment() {
    log_info "Verifying deployment..."

    sleep 3

    # Check API health
    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        # Check if services are running
        for service in unfollow-tracker-api unfollow-tracker-fetcher unfollow-tracker-scheduler; do
            if systemctl is-active --quiet \$service; then
                echo \"\$service: running\"
            else
                echo \"\$service: NOT running\"
                systemctl status \$service --no-pager
                exit 1
            fi
        done

        # Test API health endpoint
        if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
            echo 'API health check: OK'
        else
            echo 'API health check: FAILED'
            exit 1
        fi
    " || error_exit "Deployment verification failed"

    log_success "Deployment verified successfully"
}

# Cleanup local build files
cleanup() {
    log_info "Cleaning up local build files..."
    rm -rf "$LOCAL_DIR/build"
    log_success "Cleanup complete"
}

# Rollback function
rollback() {
    log_warning "Rolling back deployment..."

    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        if [[ -d $BACKUP_DIR/$APP_NAME-$TIMESTAMP ]]; then
            # Stop services
            systemctl stop unfollow-tracker-api 2>/dev/null || true
            systemctl stop unfollow-tracker-fetcher 2>/dev/null || true
            systemctl stop unfollow-tracker-scheduler 2>/dev/null || true

            # Restore backup
            rm -rf $REMOTE_DIR
            cp -r $BACKUP_DIR/$APP_NAME-$TIMESTAMP $REMOTE_DIR

            # Restart services
            systemctl start unfollow-tracker-api
            systemctl start unfollow-tracker-fetcher
            systemctl start unfollow-tracker-scheduler

            echo 'Rollback completed'
        else
            echo 'No backup found for rollback'
            exit 1
        fi
    "
}

# Display deployment summary
display_summary() {
    echo ""
    echo "=========================================="
    echo "  Deployment Complete!"
    echo "=========================================="
    echo ""
    echo "Deployed to: $SSH_USER@$SSH_HOST:$SSH_PORT"
    echo "Remote directory: $REMOTE_DIR"
    echo "Timestamp: $TIMESTAMP"
    echo ""
    echo "Services status:"
    ssh $SSH_OPTS "$SSH_USER@$SSH_HOST" "
        systemctl is-active unfollow-tracker-api &>/dev/null && echo '  API: running' || echo '  API: stopped'
        systemctl is-active unfollow-tracker-fetcher &>/dev/null && echo '  Fetcher: running' || echo '  Fetcher: stopped'
        systemctl is-active unfollow-tracker-scheduler &>/dev/null && echo '  Scheduler: running' || echo '  Scheduler: stopped'
    "
    echo ""
    echo "View logs:"
    echo "  ssh $SSH_USER@$SSH_HOST 'journalctl -u unfollow-tracker-api -f'"
    echo ""
    echo "=========================================="
}

# Main execution
main() {
    echo "=========================================="
    echo "  Unfollow Tracker - Deploy"
    echo "=========================================="
    echo ""

    parse_ssh_target
    build_ssh_opts
    test_ssh
    check_prerequisites
    build_binaries
    create_backup
    deploy_files
    run_migrations
    restart_services
    verify_deployment
    cleanup
    display_summary
}

# Handle interrupts
trap 'log_error "Deployment interrupted"; exit 1' INT TERM

# Run main function
main "$@"