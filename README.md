# Blog Application API

This is a RESTful API for a Blog application built with Go and Fiber, providing endpoints for user management, post creation and interaction, commenting, reactions, and more.

## Features

- User registration, authentication, and profile management
- Post creation, updating, and deletion
- Commenting on posts
- Liking and disliking posts
- Bookmarking posts
- Following/unfollowing users
- User avatar upload
- Email verification
- Password reset functionality
- Firebase integration for authentication
- CORS support
- Health check endpoint
- Graceful shutdown

## Getting Started

1. Clone the repository
2. Install dependencies
3. Set up your environment variables
4. Run the server


clone this repo
go mod download
cp .env.example .env
# Edit .env with your configuration
go run main.go


## Configuration

The application uses environment variables for configuration. Copy the `.env.example` file to `.env` and adjust the values as needed:

- `HOST`: The host to run the server on
- `PORT`: The port to run the server on
- `DATABASE_URL`: The URL for your database connection
- `FIREBASE_CONFIG`: Path to your Firebase configuration file

## Authentication

This API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints, include the JWT token in the Authorization header of your requests:


Authorization: Bearer <your_jwt_token>


## API Endpoints

### Authentication
- `POST /login`: Login and receive a JWT token
- `POST /register`: Register a new user
- `POST /refresh`: Refresh the JWT token
- `POST /logout`: Logout the user
- `GET /verifyemail/:email`: Check if an email is verified
- `PUT /reset-password`: Reset user password

### User Management
- `GET /users`: Get user profile
- `GET /users/:username`: Get user details
- `PUT /users/:id`: Update user profile
- `POST /users/:id/avatar`: Upload user avatar
- `GET /users/uploads/avatars/:filename`: Get user avatar image
- `POST /users/follow/:followingID`: Follow a user
- `DELETE /users/unfollow/:followingID`: Unfollow a user
- `GET /users/:id/followers`: Get user followers
- `GET /users/:id/following`: Get users being followed
- `GET /users-emails`: Get all usernames and emails

### Post Management
- `GET /posts`: Get all posts
- `POST /posts`: Create a new post
- `GET /posts/:username/:slug`: Get post by username and slug
- `PUT /posts/:id`: Update a post
- `DELETE /posts/:id`: Delete a post
- `GET /users/:id/posts`: Get posts by user
- `GET /uploads/:filename`: Get post image

### Comment Management
- `POST /posts/:id/comments`: Add a comment to a post
- `GET /posts/:id/comments`: Get comments and count for a post
- `PUT /comments/:id`: Update a comment
- `DELETE /comments/:id`: Delete a comment

### Reaction Management
- `POST /posts/:id/like`: Like a post
- `POST /posts/:id/dislike`: Dislike a post
- `GET /posts/:post_id/reactions`: Get reactions for a post

### Bookmark Management
- `POST /users/:post_id/bookmark`: Bookmark a post
- `GET /users/post/bookmarks`: Get user's bookmarks
- `GET /:post_id/bookmarkscount`: Get bookmark count for a post

### Contact Management
- `POST /contact-us`: Submit a contact form

## Error Handling

The API uses standard HTTP status codes to indicate the success or failure of requests. In case of an error, the response will include a JSON object with an `error` field describing the issue.

## Middleware

- CORS middleware for handling Cross-Origin Resource Sharing
- Logger middleware for request logging
- Health check middleware for monitoring application health
- Authentication middleware for protecting routes

## Database

The application uses a database for data persistence. Make sure to set up your database and provide the correct `DATABASE_URL` in the `.env` file.

## Firebase Integration

This API integrates with Firebase for authentication. Ensure you have set up a Firebase project and provided the correct configuration file path in the `FIREBASE_CONFIG` environment variable.

---

For any issues or feature requests, please open an issue on our GitHub repository.
