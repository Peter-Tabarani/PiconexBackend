#!/bin/bash
# updateBackend.sh - watches for updates and restarts backend automatically

BACKEND_DIR=/home/piconex/backend
RESTART_SCRIPT=$BACKEND_DIR/restart.sh
LOG_FILE=$BACKEND_DIR/log/updateBackend.log

cd $BACKEND_DIR || exit

# If already running in background, exit
if [[ $PPID != 1 ]]; then
    # Relaunch itself in background
    nohup "$0" >> "$LOG_FILE" 2>&1 < /dev/null &
    exit 0
fi

echo "[$(date)] 🟢 Update watcher started." >> "$LOG_FILE"

while true; do
    git remote update
    LOCAL=$(git rev-parse @)
    REMOTE=$(git rev-parse @{u})

    if [ "$LOCAL" != "$REMOTE" ]; then
        echo "[$(date)] 🔄 Update detected. Pulling changes..." >> "$LOG_FILE"
        git pull origin main >> "$LOG_FILE" 2>&1
        go mod tidy >> "$LOG_FILE" 2>&1
        echo "[$(date)] 🔧 Restarting backend..." >> "$LOG_FILE"
        $RESTART_SCRIPT >> "$LOG_FILE" 2>&1
    fi

    sleep 60
done

