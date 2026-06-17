#!/bin/bash
# Campus AI Forum - Start all backend services

export MOCK_DATA_DIR=/root/campus-forum/shared/mock-data

# Kill any existing instances
pkill -f "user-service" 2>/dev/null
pkill -f "forum-service" 2>/dev/null
pkill -f "admin-service" 2>/dev/null
sleep 1

# Start services
nohup /root/campus-forum/user-service > /var/log/user-service.log 2>&1 &
echo "user-service PID: $!"

ADMIN_SERVICE_URL=http://127.0.0.1:8005 \
nohup /root/campus-forum/forum-service > /var/log/forum-service.log 2>&1 &
echo "forum-service PID: $!"

ADMIN_SERVICE_PORT=8005 \
nohup /root/campus-forum/admin-service > /var/log/admin-service.log 2>&1 &
echo "admin-service PID: $!"

echo "All services started."
