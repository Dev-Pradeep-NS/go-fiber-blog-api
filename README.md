# Social Media API

This is a RESTful API for a social media application, providing endpoints for user management, post creation and interaction, commenting, and reactions.

## Table of Contents
1. [Features](#features)
2. [Getting Started](#getting-started)
3. [Authentication](#authentication)
4. [API Endpoints](#api-endpoints)
5. [Error Handling](#error-handling)
6. [Rate Limiting](#rate-limiting)

## Features

- User registration and authentication
- User profile management (including avatar upload)
- Post creation, updating, and deletion
- Commenting on posts
- Reacting to posts (like, etc.)
- Following/unfollowing users
- Fetching user posts and followers

## Getting Started

1. Clone the repository
2. Install dependencies
3. Set up your environment variables
4. Run the server


git clone https://github.com/yourusername/social-media-api.git
cd social-media-api
npm install
cp .env.example .env
# Edit .env with your configuration
npm start


## Authentication

This API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints, include the JWT token in the Authorization header of your requests:


Authorization: Bearer <your_jwt_token>


## API Endpoints

### Authentication
- `POST /register`: Register a new user
- `POST /login`: Login and receive a JWT token

### User Management
- `GET /users/:id`: Get user details
- `PUT /users/:id`: Update user details
- `POST /users/:id/avatar`: Upload user avatar
- `POST /users/:id/follow`: Follow a user
- `DELETE /users/:id/unfollow/:targetId`: Unfollow a user
- `GET /users/:id/followers`: Get user followers
- `GET /users/:id/following`: Get users being followed

### Post Management
- `GET /posts`: Get all posts
- `POST /posts`: Create a new post
- `PUT /posts/:id`: Update a post
- `DELETE /posts/:id`: Delete a post
- `GET /users/:id/posts`: Get posts by user
- `GET /posts/:username/:slug`: Get post by username and slug

### Comment Management
- `POST /posts/:id/comments`: Add a comment to a post
- `PUT /comments/:id`: Update a comment
- `DELETE /comments/:id`: Delete a comment

### Reaction Management
- `POST /posts/:id/reactions`: Add a reaction to a post
- `DELETE /posts/:id/reactions`: Remove a reaction from a post

## Error Handling

The API uses standard HTTP status codes to indicate the success or failure of requests. In case of an error, the response will include a JSON object with an `error` field describing the issue.

## Rate Limiting

To prevent abuse, this API implements rate limiting. Please refer to the response headers for information on your current rate limit status.

---

For detailed information on request and response formats, please refer to our [API Documentation](API_DOCS.md).

For any issues or feature requests, please [open an issue](https://github.com/yourusername/social-media-api/issues) on our GitHub repository.
