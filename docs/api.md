# Stream File Locally API

## Run With Docker Compose

```bash
docker compose -f deploy/docker-compose.yml --env-file deploy/.env up -d
```

Default base URL:

```text
https://cloudcalls.easipath.com/stream-file-locally
```

Health check:

```bash
curl https://cloudcalls.easipath.com/stream-file-locally/health
```

## OpenAPI And Swagger

Swagger UI:

```text
https://cloudcalls.easipath.com/stream-file-locally/swagger
```

Raw OpenAPI YAML:

```text
https://cloudcalls.easipath.com/stream-file-locally/openapi.yaml
```

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `HTTP_PORT` | `8080` | Container HTTP port. |
| `ROUTE_PREFIX` | `/stream-file-locally` | Prefix for health and API routes. |
| `PUBLIC_BASE_URL` | empty | External base URL used in upload responses. |
| `STORAGE_DRIVER` | `disk` | Use `disk` or `cassandra`. |
| `DISK_STORAGE_PATH` | `/data/stream-file-locally` | Container path for disk storage. |
| `HOST_STORAGE_PATH` | `./data` | Host folder bind-mounted into the container for disk storage. |
| `CASSANDRA_HOSTS` | `safer.easipath.com` | Comma-separated Cassandra hosts. |
| `CASSANDRA_KEYSPACE` | `biatechwallet_stream_file_locally` | Cassandra keyspace. |
| `CASSANDRA_REPLICATION_FACTOR` | `1` | Replication factor for app-managed keyspace migration. |

## Upload

```bash
curl -X POST https://cloudcalls.easipath.com/stream-file-locally/api/v1/files/upload \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "hello.txt",
    "content_type": "text/plain",
    "base64_content": "SGVsbG8sIHdvcmxkIQ=="
  }'
```

Example response:

```json
{
  "data": {
    "id": "a7b8d8d4-bf8a-43cf-8a36-5c71c62c60d9",
    "name": "hello.txt",
    "content_type": "text/plain",
    "url": "https://cloudcalls.easipath.com/stream-file-locally/api/v1/files/a7b8d8d4-bf8a-43cf-8a36-5c71c62c60d9/stream",
    "extension": ".txt",
    "size": 13,
    "created_at": "2026-05-15T19:00:00Z"
  }
}
```

## Metadata

```bash
curl https://cloudcalls.easipath.com/stream-file-locally/api/v1/files/{id}
```

## Stream

```bash
curl -L https://cloudcalls.easipath.com/stream-file-locally/api/v1/files/{id}/stream
```

The stream endpoint returns the stored file bytes with `Content-Type`, `Content-Length`, and `Content-Disposition` headers.

## Delete

```bash
curl -X DELETE https://cloudcalls.easipath.com/stream-file-locally/api/v1/files/{id}
```

Successful deletion returns `204 No Content`. In disk mode this removes the file folder from the host storage path. In Cassandra mode this removes the row from the `files` table.

## Disk Storage Path

The deployment Compose file uses a host bind mount instead of a Docker named volume:

```yaml
volumes:
  - "${HOST_STORAGE_PATH:-./data}:${DISK_STORAGE_PATH:-/data/stream-file-locally}"
```

On the server, with the deployment stored in `/apps/docker-compose-script/stream-file-locally`, files are physically available at:

```text
/apps/docker-compose-script/stream-file-locally/data
```

This path is not removed by `docker compose down -v` or `docker volume rm`.

## Cassandra Mode

Set:

```bash
STORAGE_DRIVER=cassandra
CASSANDRA_HOSTS=safer.easipath.com
CASSANDRA_KEYSPACE=biatechwallet_stream_file_locally
```

When Cassandra mode starts, the app creates the keyspace and `files` table if they do not exist.

## Deployment Note

The deployment server path is:

```text
/apps/docker-compose-script/stream-file-locally
```

Copying files to `ssh safer` should only be done after explicit approval.
