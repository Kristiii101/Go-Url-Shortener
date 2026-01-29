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

**Backend:**   Go (Golang) 1.23+
**Router:**    Chi (`github.com/go-chi/chi/v5`)
**Database:**  PostgreSQL (via Supabase)
**Driver:**    pgx/v5 (`github.com/jackc/pgx`)
**Frontend:**  HTML5, Tailwind CSS (CDN), Chart.js, QRCode.js

# To start the app you must:

# 1. Clone the Repository
```bash
git clone [https://github.com/YOUR_USERNAME/Go-Url-Shortener.git](https://github.com/YOUR_USERNAME/Go-Url-Shortener.git)
cd Go-Url-Shortener

# 2. Setup your Database
* Run the following script in the Supabase SQL Editor
-------------------------------------------------------------------------------------------------------{
    -- Links Table
CREATE TABLE IF NOT EXISTS links (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(20) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    is_custom BOOLEAN DEFAULT FALSE,
    is_disabled BOOLEAN DEFAULT FALSE,
    visit_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

-- Clicks/Analytics Table
CREATE TABLE IF NOT EXISTS clicks (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    link_id BIGINT REFERENCES links(id) ON DELETE CASCADE,
    visitor_hash TEXT,
    user_agent TEXT,
    referer TEXT,
    country_code VARCHAR(2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE UNIQUE INDEX idx_links_short_code ON links(short_code);
CREATE INDEX idx_clicks_link_id ON clicks(link_id);
}-------------------------------------------------------------------------------------------------------

# 3. Configure your ENV file
# Server Configuration
PORT=8080
BASE_URL=http://localhost:8080
WEB_DIR=./web

# Database Connection (Supabase Transaction Mode Recommended)
DATABASE_URL = postgresql://postgres:[YOUR-PASSWORD]@db.[id].supabase.co:5432/postgres

# Security / Rate Limiting
RATE_LIMIT_CREATE=10   # Max requests per window
RATE_LIMIT_WINDOW=60s  # Window size

# 4. Run the application
go mod tidy
go run cmd/api/main.go

Visit http://localhost:8080 in your browser to see the dashboard!

## API Endpoints
You can also interact with the service via standard REST API calls.

1. Create a Short Link
POST /v1/links

Body:
{
  "originalUrl": "[https://github.com/Kristiii101](https://github.com/Kristiii101)",
  "customAlias": "my-git",    // Optional
  "expiresAt": "2026-12-31T23:59:59Z" // Optional
}
Response:
{
  "shortCode": "my-git",
  "shortUrl": "http://localhost:8080/my-git",
  "originalUrl": "[https://github.com/Kristiii101](https://github.com/Kristiii101)",
  "expiresAt": "2026-12-31T23:59:59Z"
}

2. Get Link Stats
GET /v1/links/{short_code}/stats

Response:

JSON

{
  "total_clicks": 124,
  "last_clicked_at": "2026-01-26T14:30:00Z",
  "daily": [
    { "day": "2026-01-25", "clicks": 10 },
    { "day": "2026-01-26", "clicks": 5 }
  ]
}

3. Redirect
GET /{short_code}

Redirects to the original URL (307 Temporary Redirect).

Returns 410 Gone if the link has expired.

Returns 404 Not Found if the code doesn't exist.

#* Project Structure

URL_Shortener/
├─ .vscode/
│  ├─ c_cpp_properties.json
│  ├─ launch.json
│  └─ settings.json
├─ cmd/
│  └─ api/
│     └─ main.go
├─ db/
│  └─ migrations/
│     ├─ 01_init_links.down.sql
│     ├─ 01_init_links.up.sql
│     ├─ 02_clicks.down.sql
│     ├─ 02_clicks.up.sql
│     ├─ 03_indexes.down.sql
│     ├─ 03_indexes.up.sql
│     ├─ 04_views.down.sql
│     └─ 04_views.up.sql
├─ internal/
│  ├─ app/
│  │  └─ app.go
│  ├─ config/
│  │  └─ config.go
│  ├─ domain/
│  │  ├─ click.go
│  │  ├─ errors.go
│  │  ├─ link.go
│  │  ├─ reserved.go
│  │  ├─ stats.go
│  │  └─ validation.go
│  ├─ http/
│  │  ├─ handlers/
│  │  │  ├─ health.go
│  │  │  ├─ links.go
│  │  │  ├─ redirect.go
│  │  │  ├─ static.go
│  │  │  └─ stats.go
│  │  ├─ middleware/
│  │  │  ├─ logging.go
│  │  │  ├─ ratelimit.go
│  │  │  ├─ recover.go
│  │  │  └─ requestid.go
│  │  ├─ router.go
│  │  └─ server.go
│  ├─ id/
│  │  ├─ base62.go
│  │  └─ generator.go
│  ├─ observability/
│  │  ├─ logger.go
│  │  └─ metrics.go
│  ├─ qr/
│  │  └─ generator.go
│  ├─ rate/
│  │  └─ limiter.go
│  ├─ storage/
│  │  ├─ postgres/
│  │  │  ├─ clicks_repo.go
│  │  │  ├─ db.go
│  │  │  ├─ links_repo.go
│  │  │  └─ stats_repo.go
│  │  └─ repository.go
│  └─ util/
│     ├─ hash.go
│     └─ http.go
├─ web/
│  └─ index.html
├─ .env
├─ .gitignore
├─ env.sample
├─ go.mod
├─ go.sum
├─ README.md
└─ url_shortener.session.sql
