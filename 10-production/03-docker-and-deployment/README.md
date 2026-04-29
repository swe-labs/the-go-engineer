# Track DOCKER: Containerization & Deployment

## Mission

Master the "Packaging" of Go services. Learn how to turn your Go source code into a minimal, secure, and high-performance **Docker Image**. Understand the power of **Multi-stage Builds**, **Distroless** images, and **Docker Compose** for local development. Finally, learn the principles of **CI/CD Pipelines** and **Blue/Green Deployments** for zero-downtime releases.

## Stage Ownership

This track belongs to [10 Production Operations](../README.md).

## Track Map

| ID | Type | Surface | Mission | Requires |
| --- | --- | --- | --- | --- |
| **DOCKER.1** | Lesson | [Docker Basics](./1-docker-basics) | Master the `Dockerfile` for Go. | entry |
| **DOCKER.2** | Lesson | [Multi-stage Builds](./2-multi-stage-builds) | Build in one image, run in a minimal one. | **DOCKER.1** |
| **DOCKER.3** | Lesson | [Docker Compose](./3-docker-compose) | Orchestrate a Service + Database locally. | **DOCKER.2** |
| **DEPLOY.1** | Lesson | [CI/CD Pipelines](./4-cicd-pipelines) | Automate testing and building on every push. | **DOCKER.3** |
| **DEPLOY.2** | Lesson | [Blue/Green & Rollback](./5-blue-green-and-rollback) | Learn zero-downtime deployment patterns. | **DEPLOY.1** |
| **DEPLOY.3** | Exercise | [Dockerised Service](./6-dockerised-service) | Containerize a complete Go API with DB. | **DOCKER/DEPLOY** |

## Why This Track Matters

In the modern cloud, "It works on my machine" is not enough.

1. **Isolation**: Containers ensure that your application has the exact same dependencies and environment in Production as it did in Development.
2. **Security**: Multi-stage builds and "Distroless" images remove the shell and package manager from your production image, drastically reducing the attack surface.
3. **Velocity**: Automated CI/CD pipelines allow your team to move fast by catching errors early and deploying automatically.

## Next Step

After mastering containerization, learn how to automate your development workflow. Continue to [Track CG: Code Generation](../06-code-generation).
