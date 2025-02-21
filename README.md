# Log Microservice

This microservice is responsible for monitoring a RabbitMQ queue and recording access and performance logs of an API in MongoDB, providing a centralized and efficient way to manage information.

- [Log Microservice](#log-microservice)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [How to Run](#how-to-run)
    - [Docker Compose Example](#docker-compose-example)
  - [Log Data (EventPayload)](#log-data-eventpayload)
    - [Payload Example](#payload-example)

## Prerequisites

- **Docker** and **Docker Compose** installed.
- Proper access configuration to RabbitMQ and MongoDB according to the environment variables.

## Environment Variables

Make sure to define the following environment variables for the service to function correctly:

- **QUEUE_NAME**: Name of the RabbitMQ queue where messages will be consumed.
- **RABBITHOST**: RabbitMQ connection URL (e.g., `amqp://admin:password@localhost:5672/`).
- **MONGODB_URI**: MongoDB connection URI.
- **MONGODB_DATABASE**: Name of the database to store logs.
- **MONGODB_COLLECTION**: Name of the collection where logs will be inserted.
- **MONGODB_USERNAME**: Username for MongoDB authentication.
- **MONGODB_PASSWORD**: Password for MongoDB authentication.

## How to Run

To start the service using Docker Compose, run the following command in the project's root directory:

```bash
docker-compose up --build
```

This command will build the images and start the necessary containers for the microservice to operate.

### Docker Compose Example

Below is an example of a `docker-compose.yml` file for this service:

```yaml
consumer:
  container_name: consumer
  build:
    context: .
    dockerfile: Dockerfile
  restart: unless-stopped
  environment:
    QUEUE_NAME: 'logger'
    RABBITHOST: 'amqp://user:password@rabbitmq:5672/'
    MONGODB_URI: 'mongodb://mongo:27017'
    MONGODB_DATABASE: 'logs'
    MONGODB_COLLECTION: 'logs'
    MONGODB_USERNAME: 'user'
    MONGODB_PASSWORD: 'password'
  depends_on:
    mongo:
      condition: service_healthy
    rabbitmq:
      condition: service_healthy
```

## Log Data (EventPayload)

The data sent to the queue and later stored in MongoDB must follow the format below:

- **userAgent**: string  
- **userId**: string  
- **ip**: string  
- **requestMethod**: string  
- **url**: string  
- **timeToProcess**: int  
- **data**: generic type (interface{})  

### Payload Example

```json
{
  "userAgent": "Mozilla/5.0",
  "userId": "12345",
  "ip": "192.168.1.1",
  "requestMethod": "GET",
  "url": "/api/data",
  "timeToProcess": 100,
  "data": {
    "key": "value"
  }
}
```



