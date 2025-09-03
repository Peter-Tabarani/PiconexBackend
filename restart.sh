#!/bin/bash
# restart.sh inside piconexBackend repo
# Assumes backend listens for SIGTERM for graceful shutdown

BACKEND_DIR=$(dirname "$0")
LOG_FILE="$BACKEND_DIR/log/restart.log"
BACKEND_LOG="$BACKEND_DIR/log/backend.log"

cd $BACKEND_DIR

echo "[$(date)] 🔨 Building backend..." >> "$LOG_FILE"
go build -o accessify accessify.go >> "$LOG_FILE" 2>&1

# Find running backend PID
PID=$(pgrep -f accessify)

if [ -n "$PID" ]; then
    echo "[$(date)] 🛑 Sending SIGTERM to PID $PID..." >> "$LOG_FILE"
    kill -SIGTERM $PID

    # Wait for process to exit
    while kill -0 "$PID" 2>/dev/null; do
        sleep 1
    done
    echo "[$(date)] ✅ Backend stopped." >> "$LOG_FILE"
fi

# Wait until port 8080 is free
while lsof -i :8080 >/dev/null; do
    sleep 1
done

echo "[$(date)] 🚀 Starting backend..." >> "$LOG_FILE"
nohup $BACKEND_DIR/accessify >> "$BACKEND_LOG" 2>&1 < /dev/null &
echo "[$(date)] ✅ Backend started with PID $!" >> "$LOG_FILE"
