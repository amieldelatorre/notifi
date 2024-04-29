# notifi

## Running via Docker Compose
```bash
docker compose -f compose-dev.yaml up --build -d
docker compose -f compose-dev.yaml watch # To enable rebuild on file changes
```

## Notes
Known [issue](https://github.com/testcontainers/testcontainers-go/issues/605) causing flaky tests - https://github.com/testcontainers/testcontainers-go/issues/605