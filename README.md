# Computer Inventory

A simple REST API service to manage computers assigned to employees, supporting create, read, and filtered queries by employee.

## API Endpoints

| Method | Path                         | Description                    |
| ------ | ---------------------------- | ------------------------------ |
| POST   | `/computers`                 | Add a new computer             |
| GET    | `/computers/{mac_address}`   | Get computer details by MAC    |
| GET    | `/computers`                 | List all computers             |
| GET    | `/employee/{abbr}/computers` | List computers for an employee |

## Environment Variables

- `NOTIFICATION_URL` — URL of the notification service (e.g., `http://admin-notification:8080`)

## Running Locally

```bash
make run
```

## Running go Tests

Run the API tests with:

```bash
make test
```

## Docker

Build the Docker image:

```bash
make docker
```

Run services with Docker Compose (will build the image automatically):

```bash
make docker-up
```

Stop and remove containers:

```bash
make docker-down
```

## Running integration tests on docker instances:

```bash
scripts/api-tests.sh
```

Should be done on empty db

## Docker Compose Services

- **computer-inventory** — The computer API service, exposed on port `3000`
- **admin-notification** — Notification service, exposed on port `8080`

The `computer-inventory` service depends on the notification service and uses the environment variable `NOTIFICATION_URL` to communicate with it.
