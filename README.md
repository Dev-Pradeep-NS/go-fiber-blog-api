# API Documentation

This README provides examples of how to interact with the API using curl commands. Make sure to replace `http://localhost:8000` with the appropriate base URL if your API is hosted elsewhere.

## Authentication

### Register a new user
```bash
curl -X POST http://localhost:8000/register -H "Content-Type: application/json" -d '{"username": "newuser", "email": "newuser@example.com", "password": "password123"}'
```

### Login
```bash
curl -X POST http://localhost:8000/login -H "Content-Type: application/json" -d '{"email": "newuser@example.com", "password": "password123"}'
```

After successful login, you'll receive a JWT token. Use this token in the `Authorization` header for authenticated requests.

## Posts

### Get all posts
```bash
curl -X GET http://localhost:8000/posts -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Create a new post
```bash
curl -X POST http://localhost:8000/posts -H "Authorization: Bearer YOUR_JWT_TOKEN" -H "Content-Type: application/json" -d '{"title": "New Post", "content": "Post content", "category": "General", "tags": ["tag1", "tag2"]}' -F "image=@/path/to/image.j```pg"


### Get a post by ID
```bash
curl -X GET http://localhost:8000/posts/1 -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update a post
```bash
curl -X PUT http://localhost:8000/posts/1 -H "Authorization: Bearer YOUR_JWT_TOKEN" -H "Content-Type: application/json" -d '{"title": "Updated Post Title", "content": "Updated content"}'
```

### Delete a post
```bash
curl -X DELETE http://localhost:8000/posts/1 -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get posts by user
```bash
curl -X GET http://localhost:8000/users/1/posts -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Comments

### Add a comment to a post
```bash
curl -X POST http://localhost:8000/posts/2/comments -H "Authorization: Bearer YOUR_JWT_TOKEN" -H "Content-Type: application/json" -d '{"content": "This is a comment"}'
```

### Update a comment
```bash
curl -X PUT http://localhost:8000/comments/1 -H "Authorization: Bearer YOUR_JWT_TOKEN" -H "Content-Type: application/json" -d '{"content": "Updated comment"}'
```

### Delete a comment
```bash
curl -X DELETE http://localhost:8000/comments/1 -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Reactions

### Add or update a reaction to a post
```bash
curl -X POST http://localhost:8000/posts/2/reactions -H "Authorization: Bearer YOUR_JWT_TOKEN" -H "Content-Type: application/json" -d '{"reaction_type": "like"}'
```

### Remove a reaction from a post
```bash
curl -X DELETE http://localhost:8000/posts/2/reactions -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Note: Replace `YOUR_JWT_TOKEN` with the actual JWT token received after login.
