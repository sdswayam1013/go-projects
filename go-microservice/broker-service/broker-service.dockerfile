# # -------------------------------
# # STAGE 1: BUILD THE GO BINARY
# # -------------------------------

# # Use official Go image (lightweight Alpine version)
# # We name this stage "builder" so we can refer to it later
# FROM golang:1.22-alpine AS builder

# # Create a directory inside the container to hold our app code
# # WHY: containers are isolated; we need a working directory
# RUN mkdir /app

# # Copy all project files from local machine → container
# # WHY: Docker doesn't automatically see your local files
# COPY . /app

# # Set working directory to /app
# # WHY: all next commands (like go build) will run inside this folder
# WORKDIR /app

# # Build the Go application
# # CGO_ENABLED=0 → ensures static binary (no OS dependencies)
# # CGO_ENABLED=0 disables C dependencies (cgo) during compilation and forces Go to build a
# # fully static, self-contained binary. This is important for Docker, especially when using
# # minimal base images like Alpine, which do not include standard system libraries.
# # Without this, the binary may fail at runtime due to missing dependencies.
# # In short: this ensures the compiled app is portable, lightweight, and runs reliably
# # across environments without requiring external libraries.

# # -o brokerApp → output binary name
# # ./cmd/api → path where your main.go exists
# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# # Make the binary executable
# # WHY: ensures the OS can run it inside container
# RUN chmod +x /app/brokerApp


# # -------------------------------
# # STAGE 2: CREATE LIGHTWEIGHT IMAGE
# # -------------------------------

# # Start a fresh minimal image (no Go, very small size)
# # WHY: we don't want to ship Go compiler in production image
# FROM alpine:latest

# # Create app directory in this new container
# RUN mkdir /app

# # Copy ONLY the compiled binary from builder stage
# # --from=builder → refers to first stage
# # WHY:
# #   - keeps image small
# #   - removes source code & build tools
# COPY --from=builder /app/brokerApp /app 
# #/app/brokerApp is the source path and destination of copy is /app in the new image.

# # Command to run when container starts
# # This runs your Go service
# CMD ["/app/brokerApp"]


# #I used a multi-stage Docker build where the first stage compiles the Go binary 
# #using a Go image, and the second stage uses a minimal Alpine image to run only 
# #the compiled binary. This keeps the final image small and production-ready.

FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD ["/app/brokerApp"]