# Section 20: Docker & Deployment

Welcome to the **Production Engineering Block** of The Go Engineer!

Until now, we have executed everything locally using `go run`. While this is great for development, how do you release a Go application to the world? Enter Docker.

## Multi-Stage Builds (The Go Superpower)

Because Go compiles to a single statically linked binary, it has an unfair advantage in the container world. Other languages (like Node.js or Python) require massive Docker images (500MB+) full of runtime dependencies.

Go allows you to use a **Multi-Stage Build**:

1. **Stage 1 (Builder):** Uses a heavy 800MB `golang` image to compile your code.
2. **Stage 2 (Runner):** Copies *only the compiled binary* into an empty 5MB `alpine` image.

The final Docker image you deploy to production will literally only be ~15MB!

## Exploration

1. Go into `1-dockerfile` to see exactly how a multi-stage `Dockerfile` is constructed.
2. We will use this exact Dockerfile in the final Capstone project in Section 22!
