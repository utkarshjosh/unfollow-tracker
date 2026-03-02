#!/bin/bash
#
# setup-server.sh - One-time server setup for fresh Ubuntu 22.04/24.04
# Usage: ./setup-server.sh [domain_name]
#

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
DB_NAME="unfollow_tracker"
DB_USER="unfollowtracker"
DB_PASSWORD="${DB_PASSWORD:-$(openssl rand -base64 32)}"
APP_USER="unfollowtracker"
APP_DIR="/opt/unfollow-tracker"
DOMAIN="${1:-}"

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

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        error_exit "This script must be run as root (use sudo)"
    fi
}

# Detect Ubuntu version
detect_ubuntu_version() {
    if [[ -f /etc/os-release ]]; then
        # shellcheck source=/dev/null
        source /etc/os-release
        if [[ "$ID" != "ubuntu" ]]; then
            error_exit "This script is designed for Ubuntu only. Detected: $ID"
        fi
        UBUNTU_VERSION="$VERSION_ID"
        log_info "Detected Ubuntu $UBUNTU_VERSION"
    else
        error_exit "Cannot detect OS version"
    fi
}

# Update system packages
update_system() {
    log_info "Updating system packages..."
    apt-get update || error_exit "Failed to update package list"
    apt-get upgrade -y || error_exit "Failed to upgrade packages"
    log_success "System packages updated"
}

# Install required packages
install_base_packages() {
    log_info "Installing base packages..."
    apt-get install -y \
        curl \
        wget \
        git \
        build-essential \
        ca-certificates \
        gnupg \
        lsb-release \
        software-properties-common \
        apt-transport-https \
        nginx \
        certbot \
        python3-certbot-nginx \
        ufw \
        htop \
        jq \
        || error_exit "Failed to install base packages"
    log_success "Base packages installed"
}

# Install PostgreSQL 16
install_postgresql() {
    log_info "Installing PostgreSQL 16..."

    # Add PostgreSQL repository
    if ! command -v psql &> /dev/null; then
        log_info "Adding PostgreSQL repository..."
        curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | gpg --dearmor -o /usr/share/keyrings/postgresql-keyring.gpg
        echo "deb [signed-by=/usr/share/keyrings/postgresql-keyring.gpg] http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/postgresql.list
        apt-get update
    fi

    apt-get install -y postgresql-16 postgresql-contrib || error_exit "Failed to install PostgreSQL"

    # Start and enable PostgreSQL
    systemctl start postgresql
    systemctl enable postgresql

    log_success "PostgreSQL 16 installed"
}

# Setup PostgreSQL database and user
setup_database() {
    log_info "Setting up PostgreSQL database and user..."

    # Check if database already exists
    if sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
        log_warning "Database '$DB_NAME' already exists, skipping creation"
    else
        sudo -u postgres psql -c "CREATE DATABASE $DB_NAME;" || error_exit "Failed to create database"
        log_success "Database '$DB_NAME' created"
    fi

    # Check if user already exists
    if sudo -u postgres psql -t -c "\du" | grep -qw "$DB_USER"; then
        log_warning "User '$DB_USER' already exists, skipping creation"
    else
        sudo -u postgres psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';" || error_exit "Failed to create user"
        log_success "User '$DB_USER' created"
    fi

    # Grant privileges
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;" || error_exit "Failed to grant privileges"
    sudo -u postgres psql -c "ALTER DATABASE $DB_NAME OWNER TO $DB_USER;" || error_exit "Failed to set database owner"

    # Allow local connections with password
    PG_HBA="/etc/postgresql/16/main/pg_hba.conf"
    if grep -q "^local.*$DB_NAME.*$DB_USER.*md5" "$PG_HBA"; then
        log_info "pg_hba.conf already configured for $DB_USER"
    else
        # Backup original
        cp "$PG_HBA" "${PG_HBA}.backup"
        # Update authentication method
        sed -i 's/^local\s\+all\s\+all\s\+peer/local   all             all                                     md5/' "$PG_HBA"
        sed -i 's/^host\s\+all\s\+all\s\+127.0.0.1\/32\s\+scram-sha-256/host    all             all             127.0.0.1\/32            md5/' "$PG_HBA"
        sed -i 's/^host\s\+all\s\+all\s\+::1\/128\s\+scram-sha-256/host    all             all             ::1\/128                 md5/' "$PG_HBA"
        systemctl restart postgresql
    fi

    log_success "Database setup complete"
    log_info "Database password: $DB_PASSWORD"
    log_info "Save this password in your .env file as DATABASE_URL"
}

# Install Redis
install_redis() {
    log_info "Installing Redis..."

    apt-get install -y redis-server || error_exit "Failed to install Redis"

    # Configure Redis for production
    REDIS_CONF="/etc/redis/redis.conf"
    cp "$REDIS_CONF" "${REDIS_CONF}.backup"

    # Enable persistence
    sed -i 's/^#\?save 900 1/save 900 1/' "$REDIS_CONF"
    sed -i 's/^#\?save 300 10/save 300 10/' "$REDIS_CONF"
    sed -i 's/^#\?save 60 10000/save 60 10000/' "$REDIS_CONF"

    # Set memory limit (adjust based on server RAM)
    sed -i 's/^#\?maxmemory .*/maxmemory 256mb/' "$REDIS_CONF"
    sed -i 's/^#\?maxmemory-policy .*/maxmemory-policy allkeys-lru/' "$REDIS_CONF"

    # Start and enable Redis
    systemctl restart redis-server
    systemctl enable redis-server

    log_success "Redis installed and configured"
}

# Create application user
create_app_user() {
    log_info "Creating application user '$APP_USER'..."

    if id "$APP_USER" &>/dev/null; then
        log_warning "User '$APP_USER' already exists, skipping creation"
    else
        useradd -r -s /bin/false -d "$APP_DIR" -m "$APP_USER" || error_exit "Failed to create user"
        log_success "User '$APP_USER' created"
    fi
}

# Create application directory
setup_app_directory() {
    log_info "Setting up application directory..."

    mkdir -p "$APP_DIR" || error_exit "Failed to create app directory"
    mkdir -p "$APP_DIR/bin"
    mkdir -p "$APP_DIR/web"
    mkdir -p "$APP_DIR/migrations"
    mkdir -p "$APP_DIR/logs"

    chown -R "$APP_USER:$APP_USER" "$APP_DIR"
    chmod 755 "$APP_DIR"

    log_success "Application directory created at $APP_DIR"
}

# Setup Nginx configuration
setup_nginx() {
    log_info "Setting up Nginx..."

    # Remove default site
    rm -f /etc/nginx/sites-enabled/default

    # Create app nginx config
    if [[ -n "$DOMAIN" ]]; then
        cat > /etc/nginx/sites-available/unfollow-tracker << 'EOF'
upstream api_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    listen [::]:80;
    server_name DOMAIN_PLACEHOLDER;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml application/json application/javascript application/rss+xml application/atom+xml image/svg+xml;

    # Static files
    location / {
        root /opt/unfollow-tracker/web/dist;
        try_files $uri $uri/ /index.html;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # API proxy
    location /api/ {
        proxy_pass http://api_backend/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://api_backend/health;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        access_log off;
    }
}
EOF
        sed -i "s/DOMAIN_PLACEHOLDER/$DOMAIN/g" /etc/nginx/sites-available/unfollow-tracker
        log_info "Nginx config created for domain: $DOMAIN"
    else
        # Generic config without domain
        cat > /etc/nginx/sites-available/unfollow-tracker << 'EOF'
upstream api_backend {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80 default_server;
    listen [::]:80 default_server;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml application/json application/javascript application/rss+xml application/atom+xml image/svg+xml;

    # Static files
    location / {
        root /opt/unfollow-tracker/web/dist;
        try_files $uri $uri/ /index.html;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # API proxy
    location /api/ {
        proxy_pass http://api_backend/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 300s;
        proxy_connect_timeout 75s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://api_backend/health;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        access_log off;
    }
}
EOF
        log_warning "No domain specified, using default server block"
    fi

    # Enable site
    ln -sf /etc/nginx/sites-available/unfollow-tracker /etc/nginx/sites-enabled/

    # Test nginx config
    nginx -t || error_exit "Nginx configuration test failed"

    # Reload nginx
    systemctl reload nginx
    systemctl enable nginx

    log_success "Nginx configured"
}

# Setup SSL with Certbot
setup_ssl() {
    if [[ -z "$DOMAIN" ]]; then
        log_warning "No domain specified, skipping SSL setup"
        return
    fi

    log_info "Setting up SSL with Let's Encrypt..."

    # Check if certbot is installed
    if ! command -v certbot &> /dev/null; then
        apt-get install -y certbot python3-certbot-nginx
    fi

    # Obtain certificate
    certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos --email "admin@$DOMAIN" || {
        log_warning "SSL certificate generation failed. You can run 'certbot --nginx' manually later."
        return
    }

    # Setup auto-renewal cron job
    if ! crontab -l 2>/dev/null | grep -q "certbot renew"; then
        (crontab -l 2>/dev/null; echo "0 12 * * * /usr/bin/certbot renew --quiet") | crontab -
        log_info "Certbot auto-renewal cron job added"
    fi

    log_success "SSL certificate installed for $DOMAIN"
}

# Setup firewall
setup_firewall() {
    log_info "Configuring firewall..."

    # Allow SSH
    ufw allow OpenSSH || true

    # Allow HTTP and HTTPS
    ufw allow 'Nginx Full' || true

    # Allow PostgreSQL only from localhost (already default)
    # ufw allow from 127.0.0.1 to any port 5432

    # Enable firewall
    ufw --force enable

    log_success "Firewall configured"
}

# Create systemd service files
create_systemd_services() {
    log_info "Creating systemd service files..."

    # API Service
    cat > /etc/systemd/system/unfollow-tracker-api.service << EOF
[Unit]
Description=Unfollow Tracker API Server
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=$APP_USER
Group=$APP_USER
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/bin/api
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=unfollow-tracker-api
Environment="PATH=/usr/local/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
EOF

    # Fetcher Service
    cat > /etc/systemd/system/unfollow-tracker-fetcher.service << EOF
[Unit]
Description=Unfollow Tracker Fetcher Worker
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=$APP_USER
Group=$APP_USER
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/bin/fetcher
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=unfollow-tracker-fetcher
Environment="PATH=/usr/local/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
EOF

    # Scheduler Service
    cat > /etc/systemd/system/unfollow-tracker-scheduler.service << EOF
[Unit]
Description=Unfollow Tracker Scheduler
After=network.target postgresql.service redis.service
Wants=postgresql.service redis.service

[Service]
Type=simple
User=$APP_USER
Group=$APP_USER
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/bin/scheduler
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=unfollow-tracker-scheduler
Environment="PATH=/usr/local/bin:/usr/bin:/bin"

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload

    log_success "Systemd services created"
}

# Install Go (for running migrations locally)
install_go() {
    log_info "Installing Go..."

    if command -v go &> /dev/null; then
        log_info "Go is already installed: $(go version)"
        return
    fi

    GO_VERSION="1.23.4"
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz

    # Add to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh
    # shellcheck source=/dev/null
    export PATH=$PATH:/usr/local/go/bin

    log_success "Go $GO_VERSION installed"
}

# Install golang-migrate for database migrations
install_migrate() {
    log_info "Installing golang-migrate..."

    if command -v migrate &> /dev/null; then
        log_info "golang-migrate is already installed"
        return
    fi

    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz -C /tmp
    mv /tmp/migrate /usr/local/bin/migrate
    chmod +x /usr/local/bin/migrate

    log_success "golang-migrate installed"
}

# Display summary
display_summary() {
    echo ""
    echo "=========================================="
    echo "  Server Setup Complete!"
    echo "=========================================="
    echo ""
    echo "Database Configuration:"
    echo "  Database: $DB_NAME"
    echo "  Username: $DB_USER"
    echo "  Password: $DB_PASSWORD"
    echo ""
    echo "Application Directory: $APP_DIR"
    echo "Application User: $APP_USER"
    echo ""
    echo "Services installed:"
    echo "  - PostgreSQL 16"
    echo "  - Redis"
    echo "  - Nginx"
    echo "  - Certbot (SSL)"
    echo ""
    echo "Next steps:"
    echo "  1. Copy your application files to $APP_DIR"
    echo "  2. Create .env file with DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@localhost/$DB_NAME?sslmode=disable"
    echo "  3. Run migrations: migrate -path $APP_DIR/migrations -database \"postgres://$DB_USER:$DB_PASSWORD@localhost/$DB_NAME?sslmode=disable\" up"
    echo "  4. Start services: systemctl start unfollow-tracker-api unfollow-tracker-fetcher unfollow-tracker-scheduler"
    echo ""
    echo "Service management:"
    echo "  systemctl start|stop|restart|status unfollow-tracker-api"
    echo "  systemctl start|stop|restart|status unfollow-tracker-fetcher"
    echo "  systemctl start|stop|restart|status unfollow-tracker-scheduler"
    echo ""
    echo "View logs:"
    echo "  journalctl -u unfollow-tracker-api -f"
    echo "  journalctl -u unfollow-tracker-fetcher -f"
    echo "  journalctl -u unfollow-tracker-scheduler -f"
    echo ""

    if [[ -n "$DOMAIN" ]]; then
        echo "Your application will be available at:"
        echo "  https://$DOMAIN"
    else
        echo "Your application will be available at:"
        echo "  http://$(hostname -I | awk '{print $1}')"
    fi
    echo ""
    echo "=========================================="

    # Save credentials to file
    CREDENTIALS_FILE="$APP_DIR/.credentials"
    cat > "$CREDENTIALS_FILE" << EOF
# Unfollow Tracker Server Credentials
# Generated on $(date)
# KEEP THIS FILE SECURE!

DATABASE_NAME=$DB_NAME
DATABASE_USER=$DB_USER
DATABASE_PASSWORD=$DB_PASSWORD
DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@localhost/$DB_NAME?sslmode=disable

APP_USER=$APP_USER
APP_DIR=$APP_DIR
EOF
    chmod 600 "$CREDENTIALS_FILE"
    chown "$APP_USER:$APP_USER" "$CREDENTIALS_FILE"

    log_info "Credentials saved to $CREDENTIALS_FILE (restricted access)"
}

# Main execution
main() {
    echo "=========================================="
    echo "  Unfollow Tracker - Server Setup"
    echo "=========================================="
    echo ""

    check_root
    detect_ubuntu_version

    if [[ -z "$DOMAIN" ]]; then
        log_warning "No domain specified. Usage: $0 [domain_name]"
        log_warning "Continuing without SSL configuration..."
    fi

    update_system
    install_base_packages
    install_postgresql
    setup_database
    install_redis
    create_app_user
    setup_app_directory
    setup_nginx
    setup_ssl
    setup_firewall
    create_systemd_services
    install_go
    install_migrate

    display_summary
}

# Run main function
main "$@"