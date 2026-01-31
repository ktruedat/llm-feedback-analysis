# LLM Feedback Analysis System

A web-based feedback collection system with AI-powered analysis capabilities. Users can submit feedback with ratings and
comments, while administrators gain intelligent insights through LLM-based analysis powered by OpenAI's GPT models.

![Admin Dashboard](images/admin_dashboard_1.png)

## 📋 Table of Contents

1. [Overview](#overview)
2. [Features](#features)
3. [Architecture](#architecture)
4. [Quick Start](#quick-start)
5. [Configuration](#configuration)
6. [Usage Guide](#usage-guide)
7. [Admin Dashboard](#admin-dashboard)
8. [Database Insights](#database-insights)
9. [Development](#development)
10. [Documentation](#documentation)

---

## Overview

This system fulfills a complete feedback management workflow:

- **Users** submit feedback through a clean, intuitive web interface
- **AI Analysis** automatically processes feedback when thresholds are met
- **Administrators** view insights, statistics, and AI-generated summaries
- **Distributed Tracing** provides observability into system performance

Built with modern best practices:

- **Backend**: Go with Clean Architecture + Domain-Driven Design
- **Frontend**: Next.js 14 with TypeScript and shadcn/ui components
- **Database**: PostgreSQL with automated migrations
- **AI**: OpenAI GPT-5 Mini for intelligent feedback analysis
- **Observability**: OpenTelemetry + Grafana Tempo for distributed tracing

## IMPORTANT:
To read about the design decisions, architecture, LLM integration and the approach taken, see the
[backend documentation](backend/README.md).

---

## Features

### 🎯 Core Features

- ✅ **User Feedback Submission** - Simple form with rating (1-5 stars) and comment
- ✅ **User Authentication** - Secure JWT-based authentication
- ✅ **AI-Powered Analysis** - Automatic topic clustering and sentiment analysis
- ✅ **Admin Dashboard** - Comprehensive statistics and insights
- ✅ **Analysis History** - View all past AI analyses
- ✅ **Topic Breakdown** - Drill down into specific topics and associated feedbacks

### 🤖 AI Analysis Features

- **Automatic Triggering** - Analysis runs when configured feedback threshold is reached
- **Topic Clustering** - Groups feedback by common themes (e.g., Performance, UI/UX, Bugs)
- **Sentiment Analysis** - Identifies positive, neutral, or negative sentiment per topic
- **Overall Summary** - Generates comprehensive summary of feedback trends
- **Rate Limiting** - Configurable debounce to control OpenAI API usage

### 📊 Admin Features

- **Statistics Dashboard** - Total feedbacks, average rating, rating distribution
- **Recent Submissions** - Latest feedback with timestamps
- **Analysis Results** - View AI-generated insights in clean, organized format
- **Topic Navigation** - Click on topics to see associated feedbacks
- **Feedback Management** - View all feedbacks with filtering and pagination

### 🔍 Observability

- **Distributed Tracing** - Full request tracing through all layers
- **Grafana Tempo** - Visualize performance across services
- **OpenTelemetry** - Automatic instrumentation of HTTP, database, and external API calls

---

## Architecture

The system is divided into two main components:

### Backend (Go)

- **Clean Architecture** with clear separation of concerns
- **Domain-Driven Design** with business logic in domain layer
- **SQLC** for type-safe database queries
- **Goose** for database migrations (auto-applied at startup)
- **Chi Router** for HTTP routing
- **OpenTelemetry** for distributed tracing

See [backend/README.md](backend/README.md) for detailed architecture documentation.

### Frontend (Next.js)

- **Next.js 14** with App Router
- **TypeScript** for type safety
- **shadcn/ui** for beautiful, accessible components
- **Zustand** for state management
- **TailwindCSS** for styling

See [frontend/README.md](frontend/README.md) for frontend documentation.

### Infrastructure

- **PostgreSQL** for data persistence
- **Docker Compose** for local development
- **Grafana Tempo** for trace visualization
- **OpenTelemetry Collector** for trace aggregation

---

## Quick Start

### Prerequisites

- **Docker** and **Docker Compose** (recommended)
- **Go 1.21+** (if running backend manually)
- **Node.js 18+** (if running frontend manually)
- **OpenAI API Key** - Get one at https://platform.openai.com/api-keys

### 1. Clone the Repository

```bash
git clone https://github.com/ktruedat/llm-feedback-analysis.git
cd llm-feedback-analysis
```

### 2. Configure Backend Secrets

Create a `.env` file in the `backend/` directory:

```bash
cd backend
cp .env.example .env
```

Edit `backend/.env` and add your secrets:

```dotenv
# Generate with: openssl rand -base64 64
JWT_SECRET=your_jwt_secret_key_here

# Get from: https://platform.openai.com/api-keys
LLM_ANALYSIS_OPENAI_API_KEY=your_openai_api_key_here
```

**How to generate JWT_SECRET:**

```bash
openssl rand -base64 64
```

**How to get OpenAI API Key:**

1. Go to https://platform.openai.com/api-keys
2. Click "Create new secret key"
3. Copy the key
4. Paste it into your `.env` file

### 3. Start All Services

From the root directory:

```bash
docker-compose up -d --build
```

This starts:

- PostgreSQL database (port 5432)
- Backend API (port 8080)
- Frontend app (port 3000)
- Grafana Tempo (port 3200)
- OpenTelemetry Collector (ports 4317, 4318)

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Grafana Tempo**: http://localhost:3200

### 5. Wait for Services to Start

The backend automatically:

- Connects to PostgreSQL
- Applies database migrations
- Creates default admin user
- Starts the analyzer service

Check logs:

```bash
docker-compose logs -f feedback-backend
```

You should see:

```
starting server on port: 8080
```

---

## Configuration

### Backend Configuration (`backend/config.yaml`)

The backend configuration is **ready to use** as-is, but you can customize:

#### LLM Analysis Settings

```yaml
llm_analysis:
  # Minimum feedbacks before triggering analysis (default: 7)
  min_new_feedbacks_for_analysis: 7

  # Maximum feedbacks to include in single analysis (default: 50)
  max_feedbacks_in_context: 50

  # Enable rate limiting via debounce (default: false)
  enable_debounce: false

  # Minutes to wait after last feedback before analyzing (if debounce enabled)
  debounce_minutes: 1

  # OpenAI model to use
  openai_model: "gpt-5-mini-2025-08-07"

  # Maximum tokens per request (context window management)
  max_tokens_per_request: 5000
```

#### Server Settings

```yaml
server:
  host: 0.0.0.0
  port: 8080
  graceful_shutdown_seconds: 30
```

#### Database Settings

```yaml
database:
  dsn: postgres://postgres:postgres@postgres:5432/feedback_db?sslmode=disable
```

#### JWT Settings

```yaml
jwt:
  algorithm: "HS256"
  expiration_hours: 24
```

### Environment Variables (`.env`)

**Required secrets** that must be in `backend/.env`:

```dotenv
JWT_SECRET=<your-generated-secret>
LLM_ANALYSIS_OPENAI_API_KEY=<your-openai-api-key>
```

**⚠️ Never commit the `.env` file to version control!**

---

## Usage Guide

### Step 1: Register a User

![Sign In Page](images/signin.png)

1. Open http://localhost:3000
2. Click "Don't have an account? Sign up"
3. Enter email and password
4. Click "Sign Up"
5. You'll be automatically logged in

### Step 2: Submit Feedback

![Feedback Submission Page](images/feedback_page.png)

1. Select a rating (1-5 stars)
2. Enter your feedback comment
3. Click "Submit Feedback"
4. You'll see a success message

**Submit multiple feedbacks** to trigger AI analysis. By default, **7 feedbacks** are needed.

### Step 3: Submit Enough Feedbacks to Trigger Analysis

You need to submit at least **7 feedbacks** (or the configured `min_new_feedbacks_for_analysis`) to trigger automatic AI
analysis.

**Quick way to test:** Use example feedbacks from [feedback_examples.md](feedback_examples.md):

```markdown
Example feedbacks you can use:

- "The performance improvements in this release are massive. Everything feels snappy again." (Rating: 5)
- "Memory usage slowly increases over time—might be worth checking for leaks." (Rating: 3)
- "RAM consumption is getting out of control. I have to restart it every few hours." (Rating: 2)
- "Overall a nice update. Things feel smoother, even if some changes are subtle." (Rating: 4)
- "The UI looks cleaner, but the app definitely uses more RAM now." (Rating: 3)
- "Syntax highlighting sometimes breaks in larger files, especially with Rust." (Rating: 3)
- "This update completely broke my workflow. Constant freezes and crashes." (Rating: 1)
```

### Step 4: Monitor Analysis Trigger

Once you submit the 7th feedback, check the backend logs:

```bash
docker-compose logs -f feedback-backend
```

You should see:

![Analysis Triggered Logs](images/logs_analysis_triggered.png)

```
Analyzer: threshold reached, processing 7 feedbacks
Analyzer: analysis completed successfully
```

### Step 5: Access Admin Dashboard

![Admin Dashboard](images/admin_dashboard_2.png)

**Default admin credentials:**

- **Email**: `admin@mail.com`
- **Password**: `adminpass`

1. Logout from your user account
2. Login with admin credentials
3. You'll be redirected to the admin dashboard

**⚠️ Security Note:** These are demo credentials. In production, use proper admin account generation.
See [backend/README.md](backend/README.md#security-considerations) for details.

---

## Admin Dashboard

The admin dashboard provides comprehensive insights into user feedback.

### Overview Statistics

![Admin Dashboard Overview](images/admin_dashboard_1.png)

The dashboard displays:

- **Total Feedbacks** - Number of submissions received
- **Average Rating** - Overall rating across all feedbacks
- **Rating Distribution** - Count of feedbacks per rating level (1-5 stars)
- **Recent Submissions** - Latest feedbacks with timestamps

### Latest AI Analysis

![Analysis Summary](images/admin_dashboard_2.png)

The "Latest Analysis" section shows:

- **Overall Sentiment** - Positive, neutral, or negative
- **Analysis Summary** - AI-generated summary of feedback trends
- **Total Feedbacks Analyzed** - Number of feedbacks included in this analysis
- **Analysis Date** - When the analysis was performed
- **View Full Analysis** - Link to detailed analysis page

### Analysis History

![Analysis History Page](images/analysis_history_page.png)

View all past analyses:

- Navigate to "Analyses" from the dashboard
- See chronological list of all AI analyses
- Click on any analysis to view details

### Detailed Analysis View

![Analysis Detail Page](images/analysis_page.png)

Each analysis includes:

1. **Overall Summary** - High-level insights and trends
2. **Overall Sentiment** - Aggregate sentiment classification
3. **Topics Identified** - Common themes found in feedback

#### Topic Breakdown

Each topic shows:

- **Topic Name** - e.g., "Performance", "UI/UX", "Bugs"
- **Sentiment** - Positive, neutral, or negative for this topic
- **Feedback Count** - Number of feedbacks in this topic
- **Summary** - AI-generated summary of this topic
- **View Details** - Link to see associated feedbacks

### Topic Detail View

![Topic Detail Page](images/topic_page.png)

Click on a topic to see:

- All feedbacks categorized under this topic
- Individual ratings
- User comments
- Timestamps

This allows administrators to drill down from high-level insights to specific user feedback.

---

## Database Insights

The system stores all data in PostgreSQL. You can inspect the database directly:

### Connect to Database

```bash
docker exec -it feedback-postgres psql -U postgres -d feedback_db
```

### View Tables

```sql
-- List all tables
\dt feedback.*

-- View users
SELECT id, email, roles, created_at
FROM feedback.users;

-- View feedbacks
SELECT id, rating, comment, analyzed, created_at
FROM feedback.feedbacks
LIMIT 10;

-- View analyses
SELECT id, overall_sentiment, total_feedbacks_analyzed, analysis_timestamp
FROM feedback.analyses
ORDER BY analysis_timestamp DESC
LIMIT 5;

-- View topic analyses
SELECT topic, sentiment, feedback_count, summary
FROM feedback.topic_analyses
WHERE analysis_id = '<analysis-id>';
```

### Example Screenshots

**Users Table:**

![Users Table](images/users_table.png)

**Feedbacks Table:**

![Feedbacks Table](images/feedbacks_table.png)

**Analyses Table:**

![Analyses Table](images/analyses_table.png)

---

## Development

### Running Services Individually

#### Backend Only

```bash
cd backend

# Install dependencies
go mod download

# Create .env file
cp .env.example .env
# Edit .env with your secrets

# Start PostgreSQL (required)
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=feedback_db \
  -p 5432:5432 \
  postgres:16

# Run backend
go run cmd/main.go
```

Backend runs on http://localhost:8080

#### Frontend Only

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

Frontend runs on http://localhost:3000

**Note:** Frontend requires backend API to be running.

### Building for Production

#### Backend

```bash
cd backend
go build -o bin/backend cmd/main.go
./bin/backend
```

#### Frontend

```bash
cd frontend
npm run build
npm start
```

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests (if implemented)
cd frontend
npm test
```

---

## Documentation

### Detailed Documentation

- **[Backend Documentation](backend/README.md)** - Complete backend architecture, design decisions, and technologies
- **[Frontend Documentation](frontend/README.md)** - Frontend structure and components
- **[Task Specification](task.md)** - Original requirements and evaluation criteria

### Backend Architecture

The backend uses:

- **Clean Architecture** - Separation of concerns across layers
- **Domain-Driven Design** - Business logic in pure domain layer
- **SQLC** - Type-safe SQL query generation
- **Goose** - Automated database migrations
- **OpenTelemetry** - Distributed tracing

Read more: [backend/README.md](backend/README.md)

### API Documentation

Interactive API documentation (Swagger UI):

**http://localhost:8080/swagger/index.html**

Includes:

- All endpoints and HTTP methods
- Request/response schemas
- Authentication requirements
- Try API calls directly from browser

### Database Schema

Schema is managed via migrations in `backend/migrations/`:

1. `20260129000000_create_feedback_table.sql` - Feedback storage
2. `20260130000000_create_user_table.sql` - User authentication
3. `20260130000001_add_user_id_to_feedbacks.sql` - Link feedbacks to users
4. `20260130000002_create_default_admin_user.sql` - Demo admin
5. `20260130000003_create_analysis_tables.sql` - AI analysis storage
6. `20260130000004_add_topic_enum.sql` - Topic categorization

Migrations are **automatically applied** at application startup.

---

## How It Works

### AI Analysis Workflow

1. **User submits feedback** → Saved to database
2. **Analyzer monitors** → Counts unanalyzed feedbacks
3. **Threshold reached** → Analyzer triggers (default: 7 feedbacks)
4. **Fetch feedbacks** → Retrieves unanalyzed feedbacks (up to 50)
5. **Build prompt** → Constructs analysis request for OpenAI
6. **Call OpenAI API** → Sends feedbacks for analysis
7. **Parse response** → Extracts topics, sentiments, summary
8. **Save analysis** → Stores results in database
9. **Mark feedbacks** → Marks feedbacks as analyzed
10. **Admin views** → Results appear in dashboard

### Rate Limiting

Configure in `backend/config.yaml`:

```yaml
llm_analysis:
  # Option 1: Count-based triggering (default)
  min_new_feedbacks_for_analysis: 7

  # Option 2: Time-based debounce (optional)
  enable_debounce: true
  debounce_minutes: 5  # Wait 5 minutes after last feedback
```

**Count-based:** Analysis triggers immediately when 7 feedbacks are submitted.

**Debounce:** After 7 feedbacks, wait 5 minutes before analyzing (prevents rapid API calls).

### Observability

The system includes full distributed tracing:

1. **Access Grafana Tempo**: http://localhost:3200
2. **Search for traces**: Service name = `llm-feedback-analysis`
3. **View request flow**: See time spent in each layer

See [backend/README.md#observability](backend/README.md#observability--distributed-tracing) for details.

---

## Troubleshooting

### Backend won't start

**Check logs:**

```bash
docker-compose logs backend
```

**Common issues:**

- Missing `.env` file → Create `backend/.env` with secrets
- Invalid OpenAI API key → Verify key at https://platform.openai.com/api-keys
- PostgreSQL not ready → Wait a few seconds for database to initialize

### Frontend can't connect to backend

**Check backend is running:**

```bash
curl http://localhost:8080/health
```

**Check backend logs:**

```bash
docker-compose logs backend
```

### Analysis not triggering

**Check configuration:**

- Verify `min_new_feedbacks_for_analysis` in `backend/config.yaml`
- Default is 7 feedbacks

**Check analyzer logs:**

```bash
docker-compose logs -f backend | grep Analyzer
```

**Verify feedbacks are submitted:**

```bash
docker exec -it feedback-postgres psql -U postgres -d feedback_db \
  -c "SELECT COUNT(*) FROM feedback.feedbacks WHERE analyzed = false;"
```

### Database connection failed

**Check PostgreSQL:**

```bash
docker-compose ps postgres
```

**Should show:** `Up (healthy)`

**Restart if needed:**

```bash
docker-compose restart postgres
```

---

## Project Structure

```
llm-feedback-analysis/
├── backend/                    # Go backend with Clean Architecture
│   ├── cmd/                    # Application entry point
│   ├── internal/               # Internal packages
│   │   ├── domain/             # Domain entities (pure business logic)
│   │   └── app/                # Application layer
│   │       ├── handlers/       # HTTP handlers
│   │       ├── services/       # Business workflows
│   │       ├── repository/     # Data persistence
│   │       └── external/       # External API integrations
│   ├── migrations/             # Database migrations (Goose)
│   ├── config.yaml             # Application configuration
│   ├── .env                    # Secrets (create from .env.example)
│   └── README.md               # Backend documentation
│
├── frontend/                   # Next.js frontend
│   ├── app/                    # Next.js App Router
│   ├── components/             # React components
│   ├── lib/                    # Utilities and API client
│   └── README.md               # Frontend documentation
│
├── images/                     # Screenshots and documentation images
├── docker-compose.yml          # Docker Compose configuration
├── feedback_examples.md        # Example feedbacks for testing
└── README.md                   # This file
```

---

## Technology Stack

### Backend

- **Go 1.21+** - Programming language
- **Chi** - HTTP router
- **PostgreSQL** - Database
- **SQLC** - Type-safe SQL query generation
- **Goose** - Database migration management
- **JWT** - Authentication
- **OpenAI API** - LLM analysis (GPT-5 Mini)
- **OpenTelemetry** - Distributed tracing
- **Swagger** - API documentation

### Frontend

- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **shadcn/ui** - Component library
- **TailwindCSS** - Styling
- **Zustand** - State management

### Infrastructure

- **Docker & Docker Compose** - Containerization
- **Grafana Tempo** - Trace visualization
- **OpenTelemetry Collector** - Trace aggregation

---

## Module documentation

For detailed documentation:

- **Backend**: See [backend/README.md](backend/README.md)
- **Frontend**: See [frontend/README.md](frontend/README.md)
- **API**: Visit http://localhost:8080/swagger/index.html

