# Frontend - Feedback Analysis System

A Next.js frontend application for submitting feedback and viewing analytics with automatic role-based access control.

## Features

- **User Authentication**: Register and login with email/password
- **Role-Based Access Control**: Automatic routing based on user roles (user/admin)
- **Feedback Submission**: Submit feedback with 1-5 star rating and comments (for all authenticated users)
- **Admin Dashboard**: View statistics, analyze feedback, and manage submissions (admin-only)
- **Protected Routes**: Authentication-required pages with role-based protection
- **Modern UI**: Built with shadcn/ui components for a polished, accessible interface
- **Responsive Design**: Works on desktop and mobile devices

## Tech Stack

- **Next.js 16** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI component library
- **React Hook Form** - Form handling
- **Zod** - Schema validation
- **Axios** - HTTP client
- **Zustand** - State management

## Getting Started

### Prerequisites

- Node.js 20+
- npm or yarn
- Backend API running on `http://localhost:8080`

### Installation

```bash
# Install dependencies
npm install

# Run development server
npm run dev
```

The frontend will be available at `http://localhost:3000`

### Environment Variables

Create a `.env.local` file in the frontend directory:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

## Role-Based Access Control

The application automatically handles different views based on user roles, as specified in the task requirements:

### Regular Users (`user` role)

**Automatic Flow:**
1. **Register/Login**: Create an account or login with existing credentials
2. **Automatic Redirect**: After login, users are automatically redirected to the **Feedback Submission** page (`/feedback`)
3. **Submit Feedback**: Users can submit feedback with 1-5 star rating and comments
4. **Access**: Users can only access the feedback submission page

**Features Available:**
- Submit feedback with rating and comments
- View their own feedback submissions
- Logout functionality

**Restrictions:**
- Cannot access admin dashboard (`/admin`)
- Attempts to access admin routes are automatically redirected to `/feedback`

### Admin Users (`admin` role)

**Automatic Flow:**
1. **Login**: Login with admin credentials
2. **Automatic Redirect**: After login, admins are automatically redirected to the **Admin Dashboard** (`/admin`)
3. **Admin Dashboard**: Admins can view statistics, analyze feedback, and manage submissions
4. **Navigation**: Admins can navigate between admin dashboard and feedback submission page

**Features Available:**
- **Statistics Dashboard**:
  - View total feedback count
  - See average rating
  - View rating distribution charts
- **AI Analysis**:
  - Run topic clustering analysis
  - Generate feedback summaries
  - View analysis results
- **Feedback Management**:
  - View all recent submissions
  - Delete feedback entries
- **Navigation**: Access both admin dashboard and feedback submission pages

### Access Control Implementation

- **Authentication**: All routes except `/login` and `/register` require authentication
- **Role-Based Routing**: 
  - Home page (`/`) automatically redirects based on role
  - Regular users → `/feedback`
  - Admin users → `/admin`
- **Route Protection**:
  - `ProtectedRoute`: Ensures user is authenticated
  - `AdminRoute`: Ensures user is authenticated AND has admin role
- **Automatic Redirects**: Non-admins attempting to access `/admin` are automatically redirected to `/feedback`

## Project Structure

```
frontend/
├── app/                    # Next.js App Router pages
│   ├── admin/             # Admin dashboard page (admin-only)
│   ├── feedback/          # Feedback submission page (all users)
│   ├── login/             # Login page (public)
│   ├── register/          # Registration page (public)
│   ├── layout.tsx         # Root layout
│   └── page.tsx           # Home page (redirects based on role)
├── components/            # React components
│   ├── ui/                # shadcn/ui components
│   ├── ProtectedRoute.tsx # Route protection for authenticated users
│   └── AdminRoute.tsx     # Route protection for admin users
├── lib/                   # Utilities
│   ├── api-client.ts      # API client for backend communication
│   └── utils.ts           # Utility functions
├── store/                 # State management
│   └── auth-store.ts      # Authentication state with role management
└── types/                 # TypeScript types
    └── index.ts           # Shared types
```

## Usage

### User Flow

1. **Register**: Create a new account at `/register`
2. **Login**: Sign in at `/login`
3. **Automatic Redirect**: 
   - Regular users → Feedback submission page
   - Admin users → Admin dashboard
4. **Submit Feedback**: Fill out the form with rating and comment
5. **View Results**: See success/error messages after submission

### Admin Flow

1. **Login**: Sign in with admin credentials
2. **Automatic Redirect**: Admin dashboard opens automatically
3. **View Statistics**: See total feedbacks, average rating, and distribution
4. **Run Analysis**: Click "Run Analysis" to analyze feedback patterns
5. **Manage Feedback**: View and delete feedback entries as needed

## API Integration

The frontend communicates with the backend API at `/api`:

- `POST /auth/register` - User registration (public)
- `POST /auth/login` - User authentication (public, returns user info with roles)
- `POST /feedbacks` - Create feedback (requires auth)
- `GET /feedbacks` - List feedbacks (requires auth)
- `GET /feedbacks/{id}` - Get feedback by ID (requires auth)
- `DELETE /feedbacks/{id}` - Delete feedback (requires auth + admin role)

### Login Response Format

The login endpoint returns user information including roles:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "roles": ["user"]  // or ["admin"] for admin users
  }
}
```

The frontend uses this information to:
- Store user roles in the auth store
- Determine which page to redirect to after login
- Show/hide admin navigation links
- Protect admin routes

## Development

```bash
# Development server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint code
npm run lint
```

## Accessing Admin Mode

To access admin mode:

1. **Default Admin User**: A default admin user is created during database migrations
2. **Login**: Use the admin credentials to login
3. **Automatic Redirect**: You will be automatically redirected to the admin dashboard
4. **Admin Features**: All admin features will be available

**Note**: Admin credentials are configured in the backend migrations. Check the backend documentation for default admin credentials.

## Notes

- Authentication tokens are stored in localStorage
- User roles are stored in the auth store and localStorage
- The admin dashboard includes basic client-side analysis. For production, integrate with the backend LLM service endpoint
- CORS is configured on the backend to allow requests from the frontend
- Role-based access is enforced both on the frontend (routing) and backend (API endpoints)
