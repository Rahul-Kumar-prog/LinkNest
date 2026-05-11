<p align="center">
  <img src="frontend/src/assets/linknest-logo.svg" alt="LinkNest logo" width="360" />
</p>

<h1 align="center">LinkNest</h1>

<p align="center">
  A modern creator workspace for writing, optimizing, scheduling, and publishing content across social platforms from one unified dashboard.
</p>

<p align="center">
  <strong>Write once. Adapt everywhere. Publish intelligently.</strong>
</p>

---

## Overview

LinkNest is a lightweight multi-platform publishing platform designed for creators, developers, founders, marketers, and personal brands who want a cleaner workflow for managing social content.

Instead of switching between multiple social apps and dashboards, LinkNest provides a centralized workspace where users can:

- Draft content
- Preview platform-specific posts
- Connect social accounts
- Publish or schedule content
- Optimize posting times
- Generate AI-assisted posts tailored for each platform

The platform focuses on simplicity, speed, creator productivity, and future AI-powered automation workflows.

---

# Core Features

## Multi-Platform Publishing

Publish content from one dashboard to connected social platforms.

### Current Supported Integrations

- X (Twitter)
- LinkedIn

---

## AI-Powered Content Assistant

LinkNest includes an upcoming AI-assisted writing workflow designed to help creators generate better content faster.

### Planned AI Features

- Platform-specific post rewriting
- Tone adaptation
- Hook generation
- Thread creation
- CTA suggestions
- Short-form and long-form variants
- Technical-to-marketing content conversion
- Audience-aware post optimization

### Example Workflow

A technical developer post can automatically be rewritten into:
- A professional LinkedIn post
- An engaging X thread
- A concise founder update

---

## Smart Scheduling System

Future scheduling support will allow creators to:

- Schedule posts across multiple platforms
- Queue content in advance
- Manage publishing calendars
- Automate publishing workflows

---

## Intelligent Timing Recommendations

LinkNest plans to introduce timing optimization powered by account engagement analytics.

The system will analyze:

- Engagement patterns
- Reach history
- Audience activity windows
- Platform-specific performance
- Historical interaction data

This helps users publish content during higher-performing engagement windows instead of relying on generic posting advice.

---

## Live Platform Previews

Preview how content will appear before publishing.

### Current Preview Support

- X post preview
- LinkedIn post preview

---

## Media Attachments

### Current Support

- Text posts
- Single image attachments

### Future Support

- Multi-image posts
- Video uploads
- Media libraries

---

# Current X Publishing Flow

LinkNest currently uses X Web Intent for X publishing.

### Current Behavior

- The post text is automatically prefilled
- X opens in a new browser tab
- The user manually confirms publishing

### Why This Approach Exists

Browsers do not allow websites to automatically upload local files into another website’s file input.

Fully automated media publishing requires X API media upload access and additional API usage costs.

### Planned Future Improvements

- Full X API publishing
- Automatic media uploads
- Background publishing workflows

---

# Tech Stack

## Frontend

- React
- Vite
- Tailwind CSS

## Backend

- Go (`net/http`)

## Database

- PostgreSQL

## Authentication

- Google OAuth
- Session-based authentication

## Social Integrations

- X OAuth
- LinkedIn OAuth

---

# Project Structure

```text
LinkNest/
  backend/
    cmd/server/              # Go server entrypoint
    internal/db/             # PostgreSQL helpers and migrations
    internal/handlers/       # Auth, publish, integrations
    utils/jwt/               # JWT utilities

  frontend/
    src/
      Components/            # UI components and dashboard modules
      pages/                 # Main pages and channel views
      assets/                # Static assets and branding

  README.md
````

---

# Local Development

## Clone Repository

```bash
mkdir <project-directory>
cd <project-directory>

git clone <your-forked-repository-url>
```

---

## Frontend Setup

Install frontend dependencies:

```bash
cd frontend
npm install
```

Run frontend development server:

```bash
npm run dev
```

Frontend runs on:

```text
http://localhost:8000
```

---

## Backend Setup

Run backend server:

```bash
cd backend
go run ./cmd/server
```

Backend runs on:

```text
http://localhost:8080
```

---

# Verification

## Frontend Checks

```bash
cd frontend

npm run lint
npm run build
```

## Backend Checks

```bash
cd backend

go test ./...
```

---

# Product Workflow

1. Sign in to LinkNest
2. Connect X and/or LinkedIn accounts
3. Write or generate content using AI assistance
4. Preview platform-specific formatting
5. Attach media
6. Choose publishing channels
7. Publish immediately or schedule for later
8. Monitor engagement and optimize future publishing timing

---

# Future Roadmap

## Publishing & Automation

* Scheduled publishing
* Content queues
* Recurring posts
* Cross-platform automation

## AI Features

* AI post generation
* AI thread creation
* AI engagement optimization
* Audience-aware content adaptation
* AI hashtag and CTA suggestions

## Analytics

* Engagement dashboards
* Reach analytics
* Platform engagement comparisons
* Best-time posting recommendations

## Collaboration

* Team workspaces
* Shared content calendars
* Approval workflows

## Media

* Multi-image support
* Video publishing
* Asset library management

## Platform Expansion

* Instagram
* Threads
* Facebook
* YouTube Community
* Bluesky
* Mastodon

---

# Vision

LinkNest is being built as a creator-focused publishing workspace that combines:

* Content creation
* AI assistance
* Scheduling
* Analytics
* Publishing automation

into a single streamlined platform for modern creators, developers, founders, and technical professionals.
