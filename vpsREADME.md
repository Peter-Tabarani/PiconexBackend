README to help understand the VPS

ssh code = ssh piconex@178.156.189.138

4 Components That Make Up The Virtual Private Server

    1. Database
        Lives on the VPS and stores all persistent data.
        Logs backend actions to:
            log/backend.log — shows live backend output (API logs, DB errors, etc.)

    2. API / Backend
        Go server handles all API requests and communicates with the database.

    3. updateBackend.sh
        Broken down into 3 steps:
            1. Checks for updates
                Runs every 60 seconds and compares the local Git repository with the remote GitHub repository.
                If a new update is detected, it pulls the changes and runs go mod tidy to ensure dependencies are up-to-date.
                All output and errors during the pull are logged to log/updateBackend.log.
            2. Trigger restart
                After pulling updates, it calls restart.sh to rebuild and restart the backend automatically.
                The restart process logs to both log/restart.log (progress of the restart) and log/backend.log (live backend output).
            3. Continuous loop
                Always runs indefinitely in the background, so any new commits pushed to GitHub will automatically trigger the update and restart sequence.

    4. restart.sh
        This is called by updateBackend.sh whenever a new update is pulled from GitHub.
        Broken down into 3 steps:
            1. Shutdown
                Detects the current backend process (main) and sends a SIGTERM signal to gracefully shutdown.
                Waits until the process fully stops and port 8080 is freed.
            2. Rebuild
                Compiles the Go project (main.go) into a new binary (main).
                Errors, warnings, and informational messages are logged to log/restart.log for debugging.
            3. Restart
                Launches the new backend in the background using nohup, ensuring it continues running even after SSH logout.
                Logs restart actions to:
                    log/restart.log — shows restart progress and status messages