#!/bin/bash

REPO_URL="https://github.com/yourusername/yourrepository.git"
REPO_DIR="/opt/virtbrowser"
SERVICE_NAME="virtbrowser"
GO_BINARY="$REPO_DIR/virtbrowser"
REQUIRED_GO_VERSION="1.22.5"

function check_go_installed() {
    if ! command -v go &> /dev/null; then
        echo "Go $REQUIRED_GO_VERSION is required as a dependency. Please install it and try again."
        exit 1
    fi

    INSTALLED_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$INSTALLED_GO_VERSION" != "$REQUIRED_GO_VERSION" ]; then
        echo "Go $REQUIRED_GO_VERSION is required as a dependency. You have Go $INSTALLED_GO_VERSION installed. Please update it and try again."
        exit 1
    fi
}

function clone_repo() {
    if [ -d "$REPO_DIR" ]; then
        echo "Repository already exists at $REPO_DIR"
    else
        git clone "$REPO_URL" "$REPO_DIR"
    fi
}

function pull_repo() {
    cd "$REPO_DIR" || exit
    git pull --rebase
}

function rebase_repo() {
    cd "$REPO_DIR" || exit
    git fetch
    git rebase
}

function build_service() {
    cd "$REPO_DIR" || exit
    go mod download
    go build -o virtbrowser .
}

function start_service() {
    systemctl start $SERVICE_NAME
}

function stop_service() {
    systemctl stop $SERVICE_NAME
}

function restart_service() {
    systemctl restart $SERVICE_NAME
}

function setup_service() {
    cat <<EOF >/etc/systemd/system/$SERVICE_NAME.service
[Unit]
Description=VirtBrowser Service
After=network.target

[Service]
WorkingDirectory=$REPO_DIR
ExecStart=$GO_BINARY
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable $SERVICE_NAME
    systemctl start $SERVICE_NAME
}

function show_help() {
    echo "Usage: $0 {clone|pull|rebase|build|start|stop|restart|setup}"
    echo "  clone   - Clone the repository"
    echo "  pull    - Pull the latest changes from the repository"
    echo "  rebase  - Rebase the repository"
    echo "  build   - Build the Go application"
    echo "  start   - Start the service"
    echo "  stop    - Stop the service"
    echo "  restart - Restart the service"
    echo "  setup   - Setup the service to run on system startup"
}

check_go_installed

case "$1" in
    clone)
        clone_repo
        ;;
    pull)
        pull_repo
        ;;
    rebase)
        rebase_repo
        ;;
    build)
        build_service
        ;;
    start)
        start_service
        ;;
    stop)
        stop_service
        ;;
    restart)
        restart_service
        ;;
    setup)
        setup_service
        ;;
    *)
        show_help
        ;;
esac
