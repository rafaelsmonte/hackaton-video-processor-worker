## Development

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

**Run Locally**

- **Requirements**

```
 Docker version > 20.10.21
 Docker compose version > 2.12.2

 See the instructions to install on your OS at Docker's official website: https://docs.docker.com/engine/install
```

- **Setup Environment Variables**

```
You can edit the environment variables for testing environment at: docker/dev-local.env
```

- **Run the application on docker environment**

```
./scripts/run-dev.sh up // it is recommended to run ./scripts/run-dev.sh down after stopping the container execution.

# Daemon mode
./scripts/run-dev.sh upd
```

- **Stop the application**

```
./scripts/run-dev.sh down

# Stop the application and destroy database volume (clear data)
./scripts/run-dev.sh down-volume
```

- **Restart the application**

```
./scripts/run-dev.sh restart
```

- **Watch the application logs**

```
./scripts/run-dev.sh logs
```

- **Run unit tests**

```
go test ./...
```
