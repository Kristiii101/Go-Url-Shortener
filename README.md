# URL Shortener

A simple, fast URL shortener service built by Karasin Kristian with Go for ATAD Project.

## Features

- Generate short URLs with random 6-character codes
- Custom aliases for memorable links
- Click tracking and analytics
- Thread-safe in-memory storage
- No dependencies beyond Go standard library

## User Stories
- As a user, I can paste a long URL and receive a short link
- As a user, I can customize my short link
- As a user, I can view statistics: total clicks, unique visitors, geographic location
- As a user, I can generate a QR code for my short link
- As a user, I can set an expiration date for my links

## Technical Requirements
- RESTful API with endpoints for shortening, redirecting, and analytics
- Short code generation (6-8 characters, alphanumeric)
- Collision detection and handling
- Rate limiting (e.g., 10 requests/minute per IP)
- Web dashboard showing all links and statistics
- Database for URLs and click events

The server will start on `http://localhost:8080`

