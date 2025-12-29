# PomoHub API Documentation

## Authentication

### Register
**POST** `/api/v1/auth/register`

Creates a new user account.

**Body:**
```json
{
  "username": "johndoe", // English letters, numbers, underscore only
  "email": "john@example.com",
  "password": "Password1!", // Min 8 chars, 1 upper, 1 lower, 1 special
  "first_name": "John",
  "last_name": "Doe"
}
```

### Login
**POST** `/api/v1/auth/login`

Logs in a user and returns a JWT token.

**Body:**
```json
{
  "email": "john@example.com",
  "password": "Password1!"
}
```

**Response:**
```json
{
  "token": "jwt_token_here",
  "user": { ... }
}
```

## Profile

### Update Profile
**PUT** `/api/v1/users/me`

Updates the current user's profile.

**Headers:** `Authorization: Bearer <token>`

**Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "bio": "Productivity enthusiast",
  "avatar_url": "https://example.com/avatar.jpg",
  "banner_url": "https://example.com/banner.jpg"
}
```

### Get User Profile
**GET** `/api/v1/users/:username`

Gets a user's public profile.

### Get User Posts
**GET** `/api/v1/users/:username/posts`

Gets a user's posts.

## Posts

### Create Post
**POST** `/api/v1/posts`

**Headers:** `Authorization: Bearer <token>`

**Body:**
```json
{
  "content": "Just finished a 4 hour deep work session!",
  "image_url": "https://example.com/image.jpg" // Optional
}
```

### Delete Post
**DELETE** `/api/v1/posts/:id`

**Headers:** `Authorization: Bearer <token>`

## Spaces & Chat

### Create Space
**POST** `/api/v1/spaces`

### Get My Spaces
**GET** `/api/v1/spaces`

### Send Message
**POST** `/api/v1/spaces/:spaceId/messages`

### Get Messages
**GET** `/api/v1/spaces/:spaceId/messages`

## Productivity

### Todos
- **GET** `/api/v1/todos`
- **POST** `/api/v1/todos`
- **PUT** `/api/v1/todos/:id/toggle`
- **DELETE** `/api/v1/todos/:id`

### Habits
- **GET** `/api/v1/habits`
- **POST** `/api/v1/habits`
- **PUT** `/api/v1/habits/:id/toggle`
- **DELETE** `/api/v1/habits/:id`

### Pomodoro
- **POST** `/api/v1/pomodoro/sessions`
- **GET** `/api/v1/pomodoro/stats`
