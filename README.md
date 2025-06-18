# PedimeApp Microservices Template

![CI](https://github.com/frcofilippi/ms-template/actions/workflows/ci.yml/badge.svg?branch=main)

This repository is a template for starting a microservices project using Go, RabbitMQ, and PostgreSQL. It includes a simple API service, a listener service, and a RabbitMQ message broker, all orchestrated with Docker Compose.

## Features

- **API Service**: Handles HTTP requests and business logic.
- **Listener Service**: Consumes messages from RabbitMQ.
- **RabbitMQ**: Message broker for asynchronous communication.
- **PostgreSQL**: Database for persistent storage.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)

### Environment Variables

Before starting the project, set up the following environment variables as needed (you can use a `.env` file or export them in your shell):

- `POSTGRES_USER`: Username for the PostgreSQL database
- `POSTGRES_PASSWORD`: Password for the PostgreSQL database
- `POSTGRES_DB`: Name of the PostgreSQL database
- `RABBITMQ_DEFAULT_USER`: Username for RabbitMQ
- `RABBITMQ_DEFAULT_PASS`: Password for RabbitMQ

You can customize these variables in the `project/docker-compose.yml` file or provide them via a `.env` file in the `project/` directory.

### Running the Project

1. **Start Dependencies** (PostgreSQL and RabbitMQ):

   ```sh
   cd project
   docker-compose up -d db rabbitmq
   ```

2. **Start API and Listener Services** (if not already included in docker-compose):

   ```sh
   docker-compose up -d api listener
   ```

   *(Or follow the instructions in the respective service folders to run them locally.)*

3. **Stop Services**:

   ```sh
   docker-compose stop
   ```

## TODO

- [ ] Create unit tests
- [ ] Refactor the RabbitMQ connection

---

Feel free to use this template as a starting point for your own microservices projects!
