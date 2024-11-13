# Ublogging

## Uala Challenge

## Requirements

- [Golang 1.23](https://go.dev/doc/devel/release)
- [Docker](https://www.docker.com/)
- [GNU/Make](https://www.gnu.org/software/make/)

## How to run?

```bash
git clone git@github.com:wmb1207/ublogging.git
cd ublogging
go mod tidy

# Set up your .env file (check the wiki for an example)

# Setup the development database
make run-mongo

go run cmd/ublogging/main.go
```

## How to build?

```bash
git clone git@github.com:wmb1207/ublogging.git
cd ublogging
go mod tidy
go build cmd/ublogging/main.go -o ublogging
```

### How to create a user:

```bash
   curl -X POST \
   http://localhost:9099/api/users \
   -H "Content-Type: application/json" \
   -d '{ "username": "wenceslao1207", "email": "wmb1207@testing.com" }'
```

### How to follow an user
```bash
curl -X POST \
http://localhost:9099/api/users/{user_uuid}/follow \
-H "Content-Type: application/json" \
-H "Authorization: {user_uuid}"
```

### How to get the users feed
```bash
curl -X GET \
"http://localhost:9099/api/users/feed?page=1&limit=10" \
-H "Content-Type: application/json" \
-H "Authorization: {user_uuid}" 
```

### How to create a post
```bash
curl -X POST \
http://localhost:9099/api/posts \
-H "Content-Type: application/json" \
-H "Authorization: {user_uuid}" \
-d '{"Content": "testing a second user that Im following" }'
```

### How to comment on a post
```bash
curl -X POST http://localhost:9099/api/posts/{post_uuid} \
-H "Content-Type: application/json" \
-H "Authorization: {user_uuid}" \
-d '{"Content": "testing a second user that Im following" }'
```

### How to get the post with comments
```bash
curl -X GET \
http://localhost:9099/api/posts/{post_uuid} \
-H "Content-Type: application/json" \
-H "Authorization: {user_uuid}" 
```