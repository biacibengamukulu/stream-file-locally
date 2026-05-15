# Stream File Locally - Project Scope

## Goal

Build a Go service that accepts files through a JSON POST request, stores the decoded file content, and exposes a public streaming endpoint so clients can retrieve the uploaded file by URL.

The project should follow the existing DDD-style structure in the repository and keep the API, storage, configuration, documentation, and deployment workflow production-friendly.

## Functional Requirements

### File Upload

- Provide an HTTP POST endpoint that accepts a JSON request body.
- The request body must include the file content as a base64 encoded string.
- The request should also support useful file metadata such as file name and content type.
- On successful upload, the server must respond with a JSON object containing the stored file information.
- The response must include a URL that can be used to stream the file publicly.

### Storage Options

The service must support two storage modes, selected by configuration:

1. Disk storage
   - Save decoded files to the local filesystem.
   - The service will run in Docker, so the storage path must be backed by a Docker volume.
   - Provide a sensible default local/container directory.
   - The path must be configurable through environment variables.

2. Cassandra storage
   - Save file metadata and file content to Cassandra.
   - Cassandra defaults to `safer.easipath.com`, but all Cassandra settings must come from environment variables with safe defaults.
   - Add application-managed migrations because direct `cqlsh` access is not available.
   - The service should create the required keyspace/table schema automatically when Cassandra storage is enabled.

### File Streaming

- Create a public stream handler.
- The handler must receive a file URL or file identifier and stream the file bytes to the user.
- Streaming responses should set appropriate content headers such as `Content-Type` and `Content-Disposition`.

## API Documentation

- Create an OpenAPI/Swagger specification for the service.
- Create integration documentation that explains:
  - required environment variables
  - how to run with Docker Compose
  - how to upload a file
  - how to stream a file
  - how to switch between disk and Cassandra storage

## Docker And Deployment

- Use Docker Compose to run the project.
- Build/publish the container as:

```text
010309/stream-file-locally:latest
```

- The Docker Compose file must include a volume for disk-based storage.
- For deployment, the Compose file should be suitable for this server path:

```text
/apps/docker-compose-script/stream-file-locally
```

- Deployment target is accessible through:

```text
ssh safer
```

Actual deployment should only happen after explicit approval.

## Engineering Requirements

- Complete the current scaffold instead of replacing it.
- Follow the repository's existing DDD-style package layout.
- Keep configuration environment-driven.
- Use best practices for validation, error handling, and response structure.
- Keep the implementation simple and maintainable.
- Verify the work with formatting and tests/builds where possible.
