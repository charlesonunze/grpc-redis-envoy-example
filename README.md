# grpc-redis-envoy-example

json to rpc example with envoy, go, grpc, redis

## Run

Make sure you have docker installed locally

Run the services

```bash
  docker-compose up --build
```

## Testing

POST http://localhost:1337/user/login

body {"name": "John" }

.

POST http://localhost:1337/user/balance

body {"token": "token gotten from the login" }

.

POST http://localhost:1337/transactions/up

body {"amount": 200, "token": "token gotten from the login" }

.

POST http://localhost:1337/transactions/down

body {"amount": 700, "token": "token gotten from the login" }
