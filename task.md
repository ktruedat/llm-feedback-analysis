# Overview

Build a web-based feedback collection system with AI-powered analysis capabilities. The system should allow users to
submit feedback and provide administrators with insights through LLM-based analysis.

# Requirements

## 1. User Feedback Submission

Create a web page with a feedback form that includes:

- Rating field: 1-5 star rating (required)
- Comment field: Free-text feedback (required, reasonable length limit)
- Submit button: Save feedback to database

The form should:

- Validate input before submission
- Provide user feedback on successful/failed submission
- Have a clean, functional UI (styling is secondary to functionality)

## 2. Data Persistence

- Store all feedback submissions in a PostgreSQL database
- Include at minimum:
    - Rating (1-5)
    - Comment text
    - Submission timestamp
    - Unique identifier
    - Any other fields you deem necessary
- Ensure data integrity and proper schema design

## 3. Admin Dashboard

Create an admin-only view/page that displays:

**Statistics:**

- Total number of feedback submissions
- Average rating
- Rating distribution (count per rating level)
- Recent submissions (with timestamps)

**AI Analysis:**

- Topic clustering: Group feedback by common themes/topics
- Summary: Generate overall summary of feedback trends
- Display analysis results in a readable format

## 4. AI Integration

Implement backend integration with OpenAI API:

- Use an appropriate GPT model
- Fetch all feedback entries and send them for analysis
- Request the LLM to:
    - Identify common topics/themes and cluster feedback accordingly
    - Generate a summary of overall feedback sentiment and key points
- Handle API errors gracefully
- Consider rate limiting and token limits

## 5. Access Control

- Regular users should only access the feedback submission form
- Admin dashboard must be protected
- Implement authentication/authorization mechanism of your choice:
    - HTTP Basic Auth
    - Simple login with session/JWT
    - Environment-based secret token
    - Or any other approach you prefer
- Document how to access admin mode

# Technical Requirements

## Language and Framework

- Backend: Go (required)
- Database: PostgreSQL (required)
- Frontend: Any approach (plain HTML/CSS/JS, htmx, Templ, or any framework you prefer)
- LLM: OpenAI API (appropriate GPT model)

## Code Quality

- Clean, readable, and maintainable code
- Proper error handling
- Appropriate use of Go idioms and best practices
- Reasonable project structure

## Configuration

- Externalize configuration (database connection, OpenAI API key, etc.)
- Support configuration via environment variables or config file
- Document all required configuration

# Deliverables

## Source Code

- Complete, working application
- Clear project structure
- Include .gitignore

## README.md

- Setup instructions
- How to run the application
- How to access admin mode
- Required environment variables/configuration
- Any assumptions or design decisions

## Database Schema

- Schema definition (SQL file, migration, or documented in README)
- Clear explanation of your data model

## Dependencies

- go.mod and go.sum files
- Instructions for installing any additional dependencies

# Evaluation Focus

Your solution will be evaluated on:

- Functionality: Does it work as specified?
- Code Quality: Is the code clean, organized, and idiomatic?
- Error Handling: Are errors handled appropriately?
- Architecture: Is the solution well-structured and maintainable?
- Documentation: Can someone else set up and run your application?
- AI Integration: Is the OpenAI API integration implemented correctly?

# Time Expectation

This assignment should take approximately 2-3 days to complete. We value quality over speed, but also respect your time.
If you find yourself spending significantly more time, consider simplifying your approach or documenting what you would
do differently with more time.

# Submission

- Provide a link to a Git repository (GitHub, GitLab, etc.)
- Ensure the repository is accessible to reviewers
- Do not include API keys or secrets in the repository

# Optional Enhancements (Not Required)

If you have extra time and want to showcase additional skills:

- Docker/docker-compose setup
- Tests (unit or integration)
- Caching of AI analysis results
- Pagination for admin dashboard
- Advanced UI/UX
- API documentation