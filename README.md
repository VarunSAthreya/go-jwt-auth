# JWT Authentication using GO

This is a basic server written in Go for the JWT Authentication and Authorization based on user type.

## Tech Stack

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin)
- [JWT](https://jwt.io/)
- [MongoDB](https://www.mongodb.com/)

## Routes

- `/users/signup`:
  - Sign Up for a user.
  - Requires data in the body:

```json
            {
                "first_name": "John",
                "last_name": "Doe",
                "email": "john@doe.com",
                "password": "password",
                "phone": "1234567890",
                "user_type": "ADMIN" // Can only be "ADMIN" or "USER"
            }
```

- Returns user ID.

- `/users/signin`:
  - Sign In for a user.
  - Requires data in the body:

```json
                {
                    "email": "john@doe.com",
                    "password": "password"
                }
```

- Returns User object.

- `/users/:id`:
  - Get a user by ID.
  - Requires `Authorization` token in request Header.
  - Returns User object.

- `/users`:
  - Get all users, only available for `ADMIN` user type.
  - Requires `Authorization` token and `uid` (User ID) in request Header.
  - Returns total count and all users.

## Installation

### Prerequisites

- [Docker](https://www.docker.com/)

### Environment Variables

Create a `.env` file in the root directory of the project.

```.env
PORT = 8000
MONGO_URL = mongodb://mongo:27017/<db name>
SECRET_KEY = <JWT SECRET KEY>
```

- In `MONGO_URL` it is "mongodb://`mongo`:27017/\<db name>" where `mongo` is the name of the container.
- And change in `PORT` number should also be done in `Dockerfile` and `docker-compose.yml` files.

### Build

```bash
docker-compose build
```

### Run

```bash
docker-compose up
```
