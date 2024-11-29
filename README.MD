# Blog Backend Application

## Overview

This is a backend application for a blog platform. It is built using Go and the Gin framework. The application provides APIs for user authentication, managing posts, and handling comments.

## API Documentation

### Authentication

- **Register a new user**
  - **Endpoint:** `POST /auth/register`
  - **Request Body:**
    ```json
    {
      "username": "string",
      "password": "string"
    }
    ```
  - **Response:**
    ```json
    {
      "message": "User registered successfully"
    }
    ```

- **Login**
  - **Endpoint:** `POST /auth/login`
  - **Request Body:**
    ```json
    {
      "username": "string",
      "password": "string"
    }
    ```
  - **Response:**
    ```json
    {
      "token": "jwt_token"
    }
    ```

- **Get all users**
  - **Endpoint:** `GET /auth/users`
  - **Response:**
    ```json
    [
      {
        "id": "string",
        "username": "string"
      }
    ]
    ```

### Posts

- **Create a new post**
  - **Endpoint:** `POST /api/posts`
  - **Request Body:**
    ```json
    {
      "title": "string",
      "content": "string"
    }
    ```
  - **Response:**
    ```json
    {
      "id": "string",
      "title": "string",
      "content": "string",
      "author": "string"
    }
    ```

- **Get all posts**
  - **Endpoint:** `GET /api/posts`
  - **Response:**
    ```json
    [
      {
        "id": "string",
        "title": "string",
        "content": "string",
        "author": "string"
      }
    ]
    ```

- **Get a post by ID**
  - **Endpoint:** `GET /api/posts/:id`
  - **Response:**
    ```json
    {
      "id": "string",
      "title": "string",
      "content": "string",
      "author": "string"
    }
    ```

- **Update a post**
  - **Endpoint:** `PATCH /api/posts/:id`
  - **Request Body:**
    ```json
    {
      "title": "string",
      "content": "string"
    }
    ```
  - **Response:**
    ```json
    {
      "id": "string",
      "title": "string",
      "content": "string",
      "author": "string"
    }
    ```

- **Delete a post**
  - **Endpoint:** `DELETE /api/posts/:id`
  - **Response:**
    ```json
    {
      "message": "Post deleted successfully"
    }
    ```

### Comments

- **Add a comment to a post**
  - **Endpoint:** `POST /api/posts/:id/comments`
  - **Request Body:**
    ```json
    {
      "content": "string"
    }
    ```
  - **Response:**
    ```json
    {
      "id": "string",
      "content": "string",
      "author": "string",
      "postId": "string"
    }
    ```

- **Get all comments for a post**
  - **Endpoint:** `GET /api/posts/:id/comments`
  - **Response:**
    ```json
    [
      {
        "id": "string",
        "content": "string",
        "author": "string",
        "postId": "string"
      }
    ]
    ```

- **Get all comments**
  - **Endpoint:** `GET /api/posts/comments`
  - **Response:**
    ```json
    [
      {
        "id": "string",
        "content": "string",
        "author": "string",
        "postId": "string"
      }
    ]
    ```

- **Update a comment**
  - **Endpoint:** `PATCH /api/posts/comments/:commentID`
  - **Request Body:**
    ```json
    {
      "content": "string"
    }
    ```
  - **Response:**
    ```json
    {
      "id": "string",
      "content": "string",
      "author": "string",
      "postId": "string"
    }
    ```

- **Delete a comment**
  - **Endpoint:** `DELETE /api/posts/comments/:commentID`
  - **Response:**
    ```json
    {
      "message": "Comment deleted successfully"
    }
    ```

## Running the Application in a Container

### Prerequisites

- Docker
- Docker Compose

### Steps

1. **Clone the repository:**
   ```sh
   git clone https://github.com/Takeso-user/blog-backend.git
   cd blog-backend
    ```
2. **Create a .env file:**
    ```sh
    cp .env.example .env
    ```
3. **Ubdare the .env file with your values**
    ```toml
     MONGO_DATABASE=blog
     JWT_SECRET=your_jwt_secret
     MONGO_USER=admin
     MONGO_PASSWORD=adminpassword
    ```
4. **Build and run the containers:**
   ```sh
   docker-compose up --build
   ```
5. **Access the application at `http://localhost:8080` or `http://localhost:8888` for Docker**

### Conclusion
This documentation provides an overview of the blog backend application, its API endpoints, and instructions on how to run the application in a Docker container using Docker Compose.