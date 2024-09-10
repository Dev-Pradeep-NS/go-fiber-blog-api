# Social Media API Documentation

This README provides documentation for the Blog App API, detailing the available endpoints and how to use them.

## Table of Contents
1. [Authentication](#authentication)
2. [User Management](#user-management)
3. [Post Management](#post-management)
4. [Comment Management](#comment-management)
5. [Reaction Management](#reaction-management)

## Authentication

### Register a new user

```bash
curl -X POST http://localhost:8000/register \
-H "Content-Type: application/json" \
-d '{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "strongpassword123"
}'
```

### Login

```bash
curl -X POST http://localhost:8000/login \
-H "Content-Type: application/json" \
-d '{
  "email": "john@example.com",
  "password": "strongpassword123"
}'
```

## User Management

### Get user details

```bash
curl -X GET http://localhost:8000/users/1 \
-H "Authorization: Bearer <your-jwt-token>"
```

### Update user details

```bash
curl -X PUT http://localhost:8000/users/1 \
-H "Authorization: Bearer <your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{
  "username": "john_new",
  "email": "john_new@example.com",
  "bio": "Updated bio"
}'
```

### Upload user avatar

```bash
curl -X POST http://localhost:8000/users/1/avatar \
-H "Authorization: Bearer <your-jwt-token>" \
-F "avatar=@/path/to/avatar.png"
```

### Follow a user

```bash
curl -X POST http://localhost:8000/users/1/follow/2 \
-H "Authorization: Bearer <your-jwt-token>"
```

### Unfollow a user

```bash
curl -X DELETE http://localhost:8000/users/1/unfollow/2 \
-H "Authorization: Bearer <your-jwt-token>"
```

### Get user followers

```bash
curl -X GET http://localhost:8000/users/1/followers \
-H "Authorization: Bearer <your-jwt-token>"
```

### Get users being followed

```bash
curl -X GET http://localhost:8000/users/1/following \
-H "Authorization: Bearer <your-jwt-token>"
```

## Post Management

### Get all posts

```bash
curl -X GET http://localhost:8000/posts \
-H "Authorization: Bearer <your-jwt-token>"
```

### Create a new post

```bash
curl -X POST http://localhost:8000/posts \
-H "Content-Type: multipart/form-data" \
-H "Authorization: Bearer <your-jwt-token>" \
-F "title=Sample Post" \
-F "content=This is the content of the post" \
-F "category=tech" \
-F "tags=tech,programming" \
-F "image=@/path/to/image.jpg"
```

### Update a post

```bash
curl -X PUT http://localhost:8000/posts/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <your-jwt-token>" \
-d '{
  "title": "Updated Post Title",
  "content": "Updated content",
  "category": "tech"
}'
```

### Delete a post

```bash
curl -X DELETE http://localhost:8000/posts/1 \
-H "Authorization: Bearer <your-jwt-token>"
```

### Get posts by user

```bash
curl -X GET http://localhost:8000/users/1/posts \
-H "Authorization: Bearer <your-jwt-token>"
```

### Get post by username and slug

```bash
curl -X GET http://localhost:8000/posts/username/sample-post-slug \
-H "Authorization: Bearer <your-jwt-token>"
```

## Comment Management

### Add a comment to a post

```bash
curl -X POST http://localhost:8000/posts/1/comments \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <your-jwt-token>" \
-d '{
  "comment": "This is a great post!"
}'
```

### Update a comment

```bash
curl -X PUT http://localhost:8000/comments/1 \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <your-jwt-token>" \
-d '{
  "comment": "Updated comment text"
}'
```

### Delete a comment
```bash
curl -X DELETE http://localhost:8000/comments/1 \
-H "Authorization: Bearer <your-jwt-token>"
```

## Reaction Management

### Add a reaction to a post
```bash
curl -X POST http://localhost:8000/posts/1/reactions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <your-jwt-token>" \
-d '{
  "reaction_type": "like"
}'
```

### Remove a reaction from a post
```bash
curl -X DELETE http://localhost:8000/posts/1/reactions \
-H "Authorization: Bearer <your-jwt-token>"
```

Note: Replace `<your-jwt-token>` with the actual JWT token received after login.
