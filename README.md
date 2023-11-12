# Little chat API

## Overview
LittleChat API is a simple chat service that allows users to create chat sessions, send messages, and retrieve messages from a chat session. It is built using Go and the Gin web framework, and it uses a PostgreSQL database to store chat data.

## Features
- Create new chat sessions
- Send messages to a chat session
- Retrieve messages from a chat session

## Prerequisites
- Docker
- Docker Compose

## Installation
To set up the LittleChat on your local machine using Docker, follow these steps:

1. Clone the repository to your local machine.
2. Navigate to the directory containing the `docker-compose.yml` file.
3. Run `docker-compose up app` to build and start the services.
4. The application should now be accessible on `http://localhost:8080`.

## Usage
Once the server is running, you can interact with the API using the following endpoints:

- `GET /chat?id={chat_id}`: Retrieve messages from a given chat ID.
  ```sh
  curl -X GET "http://localhost:8080/chat?id=1"
  ```
- `GET /new`: Create a new chat session and receive a unique chat ID.
  ```sh
  curl -X GET "http://localhost:8080/new"
  ```
- `POST /send`: Send a message to a chat session.
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{"chat_id":1,"sender_id":1,"text":"Hello, World!"}' "http://localhost:8080/send"
  ```

## Contributing
Contributions to the LittleChat API are welcome. Please feel free to fork the repository, make changes, and submit a pull request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.
