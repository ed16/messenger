# Messenger Project

This project implements a simple messenger system with authentication and user management.

## Project Structure

The project consists of several packages and files:

- `cmd`: Contains main applications for different services.
  - `auth_service`: Authentication service.
  - `user_service`: User management service.
  - `message_service`: Message exchangement service.
- `internal`: Contains internal packages and libraries.
  - `handlers`: HTTP request handlers for different endpoints.
  - `lib`: Libraries and utilities.
  - `repository`: Data access layer.
- `domain`: Contains domain models and error definitions.
- `services`: Business logic layer.

## Setup and Installation

### Installation

1. Clone the repository:

```bash
git clone https://github.com/ed16/messenger.git
```

2. Navigate to the project directory:

```bash
cd messenger
```

3. Build and run the services using docker-compose:
```bash
colima start
cp .env_template .env
nano .env
docker-compose up
```

## Usage

The API endpoints provided by the services can be accessed through HTTP requests. Here are the available endpoints:

### Authentication Service
- `/auth/login`: Login endpoint.
- `/auth/validate-token`: Token validation endpoint.

### User Service
- `/users`: User management endpoints.
- `/users/contacts`: Add or retrieve user contacts.
- `/users/profile`: Update user profile details.

### Messages
- `/messages/user/{userId}`: Send a personal message to a user.
- `/messages/group/{groupId}`: Send a message to a group.
- `/messages/user/`: Retrieve personal messages between users.
- `/messages/group/`: Retrieve messages from a group.
- `/messages/user/{userId}/file`: Send a media file to a user.
- `/messages/group/{groupId}/file`: Send a media file to a group.
- `/messages/file`: Retrieve media files from messages.