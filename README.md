# Usage
## Run with `go run main.go`
- CLI will request confirmation on whether p8 key is included locally (local to where the program is running)
- CLI will then request Key ID, Team ID, Bundle ID and it will use the local p8 and a default push token automatically without needing this from the end user
- The program will make a request and if a 400 is returned (file is valid but something went wrong), it will try again using the sandbox environment
- Errors have been updated to give guidance on what could be wrong and what next steps should be
