# Overview

Build a web-based feedback collection system with AI-powered analysis capabilities. The system should allow users to
submit feedback and provide administrators with insights through LLM-based analysis.

# Requirements

## 1. User Feedback Submission

Create a web page with a feedback form that includes:

- [x] Rating field: 1-5 star rating (required)
- [x] Comment field: Free-text feedback (required, reasonable length limit)
- [x] Submit button: Save feedback to database

The form should:

- [x] Validate input before submission
- [x] Provide user feedback on successful/failed submission
- [x] Have a clean, functional UI (styling is secondary to functionality)

## 2. Data Persistence

- [x] Store all feedback submissions in a PostgreSQL database
- [x] Include at minimum:
    - [x] Rating (1-5)
    - [x] Comment text
    - [x] Submission timestamp
    - [x] Unique identifier
    - [x] Any other fields you deem necessary
- [x] Ensure data integrity and proper schema design

## 3. Admin Dashboard

Create an admin-only view/page that displays:

**Statistics:**

- [x] Total number of feedback submissions
- [x] Average rating
- [x] Rating distribution (count per rating level)
- [x] Recent submissions (with timestamps)

**AI Analysis:**

- [x] Topic clustering: Group feedback by common themes/topics
- [x] Summary: Generate overall summary of feedback trends
- [x] Display analysis results in a readable format

## 4. AI Integration

Implement backend integration with OpenAI API:

- [x] Use an appropriate GPT model
- [x] Fetch all feedback entries and send them for analysis
- [x] Request the LLM to:
    - [x] Identify common topics/themes and cluster feedback accordingly
    - [x] Generate a summary of overall feedback sentiment and key points
- [x] Handle API errors gracefully
- [x] Consider rate limiting and token limits

## 5. Access Control

- [x] Regular users should only access the feedback submission form
- [x] Admin dashboard must be protected
- [x] Implement authentication/authorization mechanism of your choice:
    - [x] HTTP Basic Auth
    - [x] Simple login with session/JWT
    - [x] Environment-based secret token
    - [x] Or any other approach you prefer
- [x] Document how to access admin mode

# Technical Requirements

## Language and Framework

- [x] Backend: Go (required)
- [x] Database: PostgreSQL (required)
- [x] Frontend: Any approach (plain HTML/CSS/JS, htmx, Templ, or any framework you prefer)
- [x] LLM: OpenAI API (appropriate GPT model)

## Code Quality

- [x] Clean, readable, and maintainable code
- [x] Proper error handling
- [x] Appropriate use of Go idioms and best practices
- [x] Reasonable project structure

## Configuration

- [x] Externalize configuration (database connection, OpenAI API key, etc.)
- [x] Support configuration via environment variables or config file
- [x] Document all required configuration

# Deliverables

## Source Code

- [x] Complete, working application
- [x] Clear project structure
- [x] Include .gitignore

## README.md

- [x] Setup instructions
- [x] How to run the application
- [x] How to access admin mode
- [x] Required environment variables/configuration
- [x] Any assumptions or design decisions

## Database Schema

- [x] Schema definition (SQL file, migration, or documented in README)
- [x] Clear explanation of your data model

## Dependencies

- [x] go.mod and go.sum files
- [x] Instructions for installing any additional dependencies

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

- [x] Docker/docker-compose setup
- [ ] Tests (unit or integration)
- [x] Caching of AI analysis results
- [x] Pagination for admin dashboard
- [x] Advanced UI/UX
- [x] API documentation