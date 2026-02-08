# CodeUniverse

CodeUniverse is an online programming judge and learning platform I created as part of my senior project for my bachelor's degree.

Inspired by LeetCode, this project helped me understand how online coding platforms work internally and deepen my Go programming skills.

## 📺 Video Showcase

> https://github.com/user-attachments/assets/419ae6ba-7780-4ccc-8ad8-818e720b2181
>
> Watch a complete walkthrough of CodeUniverse's features, from problem solving to submission judging and user management.

## ✨ Features

### Core Functionality

- **🎯 Problem Management**
  - Create, edit, and organize coding problems.
  - Support for problem descriptions with Markdown.
  - Hints system to guide users.
  - Test cases (public and hidden) for automatic judging.

- **💻 Code Editor & Submission**
  - Built-in Monaco editor (VS Code's editor) .
  - Multi-language support: Go, Python, C++, TypeScript, JavaScript, Java, and Ruby.
  - Real-time syntax highlighting.
  - Code execution with detailed feedback.

- **⚖️ Judging System**
  - Sandboxed code execution for security.
  - Memory usage and execution time tracking.
  - Docker-based language runners (isolated containers).
  - Test case validation with detailed results.
  - Support for partial scoring.

### User Features
- **🔐 Authentication & Security**
  - Secure user registration and login.
  - Email verification system.
  - Password reset functionality.
  - Multi-factor authentication (MFA) support.

- **👤 User Profiles**
  - Personal profile pages with avatars.
  - Problem completion statistics.
  - User socials.

- **📝 Problem Notes**
  - Private notes for each problem.
  - Track your approach and solutions.
  - Markdown support for formatting.

### Admin & Developer Tools
- **🛠️ Database Management**
  - Complete migration system using Goose.
  - Schema version control.

- **📧 Email System**
  - Mailpit integration for development.
  - Template-based emails.
  - Support for verification and password reset emails.

- **🐳 Docker Support**
  - Complete Docker Compose setup.
  - Isolated service containers.
  - Ready-to-use development environment.

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.25+
- **Framework:** Go-chi
- **Database:** PostgreSQL
- **Migrations:** Goose
- **Authentication:** JWT-based with MFA support
- **Email:** SMTP/Resend with template support

### Frontend
- **Framework:** React 19.2
- **Language:** TypeScript
- **Build Tool:** Vite
- **UI Library:** React Bootstrap
- **Code Editor:** Monaco Editor (VS Code)
- **Markdown Editor:** MDEditor
- **Charts:** Chart.js
- **Routing:** React Router v7

### Infrastructure
- **Containerization:** Docker & Docker Compose
- **Database:** PostgreSQL
- **Development Mail:** Mailpit
- **Build System:** Make

### Language Runners
- Docker containers for isolated code execution:
  - `codeuniverse-cpp` - Custom C++ execution environment
  - `codeuniverse-go` - Custom Go execution environment
  - `codeuniverse-node` - Custom JavaScript/TypeScript environment
  - Additional language support via base images

## 📁 Project Structure

```
codeuniverse/
├── cmd/                   # Application entry points
│   └── server/            # Main server application
├── internal/              # Private application code
│   ├── database/          # Database connection & queries
│   ├── handlers/          # HTTP request handlers
│   ├── judger/            # Code execution & judging logic
│   ├── logger/            # Logging utilities
│   ├── mailer/            # Email sending service
│   ├── middleware/        # HTTP middleware (auth, CORS, etc.)
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   ├── router/            # HTTP routing
│   ├── services/          # Business logic
│   └── utils/             # Helper functions
├── frontend/              # React TypeScript frontend
│   ├── src/               # Source files
│   └── public/            # Static assets
├── migrations/            # Database migrations
├── docker-builds/         # Docker images for language runners
├── docker-compose.yml     # Docker service definitions
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
