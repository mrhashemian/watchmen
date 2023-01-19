# watchmen

## Build & Run

### Manual Build and Run

Golang 1.16 is required for building the project.

1. Build the project with the following command:
    ```bash
    make build-static-vendor
    ```

2. Launch the project with:
    ```bash
    ./watchmen serve --config config.yml
    ```

### Docker using Makefile

You can build and run the project with docker compose. [Docker](https://docs.docker.com/) and
[Docker Compose](https://docs.docker.com/compose/) should be installed on the host which you are running the project.

1. Run this command to run the application with docker compose:
    ```bash
    make up
    ```

2. For shutting down the project you can run this command:
    ```bash
    make down
    ```

## Development Notes

Before starting development on this project it's better to consider these notes:

### Makefile Rules

```bash
make                    # Format, lint and static build
make build              # Build the project
make docker             # Build docker image
make format             # Format the project
make lint               # Run lintter
make up                 # Prepare and run the project with docker compose
make down               # Shut down the project running in the docker compose
```

### Database

you can find about database in cmd->migrate.go

### Endpoints

check postman collection