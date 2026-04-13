<p align="center">
  <img src="frontend/src/assets/linknest-logo.svg" alt="LinkNest logo" width="360" />
</p>

<h1 align="center">LinkNest</h1>

<p align="center">
  A clean creator dashboard for writing once and publishing across social channels.
</p>

<p align="center">
  <strong>Compose content, preview channel-specific posts, connect social profiles, and publish from one focused workspace.</strong>
</p>

## Overview

LinkNest is a lightweight social publishing dashboard built for creators who want a simple place to draft posts and send them to connected channels without jumping between multiple tabs.

The current MVP supports Google sign-in, X and LinkedIn account connections, a rich post composer, image attachment previews, and channel cards that show whether each social account is connected.

## What LinkNest Does

- Provides a single creator dashboard for drafting social posts.
- Lets users connect X and LinkedIn profiles from the channel sidebar.
- Shows live previews for X and LinkedIn so the user can see how the post will look before publishing.
- Supports text posts with one image attachment in the composer.
- Publishes to LinkedIn through the backend API when LinkedIn is selected.
- Opens X in a new browser tab with the post text prefilled to avoid using X API posting credits in the temporary flow.
- Shows publish status and platform-specific errors in the dashboard.

## Current X Publishing Behavior

LinkNest currently uses X Web Intent for X publishing:

- The post text is automatically prefilled in the X compose window.
- The user manually clicks the final Post button on X.
- Local media cannot be auto-attached in the X tab through Web Intent.

This limitation exists because browsers do not allow one website to programmatically upload a local file into another website's file input. Fully automatic X text + media publishing requires X API media upload, which consumes X API access/credits.

## Tech Stack

- Frontend: React, Vite, Tailwind CSS
- Backend: Go net/http
- Database: PostgreSQL
- Authentication: Google OAuth plus app session cookies
- Social integrations: X OAuth, LinkedIn OAuth

## Project Structure

```text
LinkNest/
  backend/
    cmd/server/              # Go HTTP server entrypoint
    internal/db/             # PostgreSQL connection, migrations, user helpers
    internal/handlers/       # Auth, connection, publish, X, LinkedIn handlers
    utils/jwt/               # JWT utilities
  frontend/
    src/
      Components/            # Dashboard, landing, sidebar, auth UI
      pages/                 # Home and channel windows
      assets/                # LinkNest logo and static assets
  README.md
```

## Local Development

Cloning the Repository

Follow the steps below to clone and set up the project locally:

```bash
# Create a new directory for your project
mkdir <your-directory-name>

# Navigate into the directory
cd <your-directory-name>

# Clone your forked repository
git clone <your-forked-repository-url>
``` 

Install frontend dependencies:

```bash
cd frontend
npm install
```

Start the frontend:

```bash
cd frontend
npm run dev
```

The frontend runs on:

```text
http://localhost:8000
```

Start the backend:

```bash
cd backend
go run ./cmd/server
```

The backend runs on:

```text
http://localhost:8080
```

## Verification

Run frontend checks:

```bash
cd frontend
npm run lint
npm run build
```

Run backend checks:

```bash
cd backend
go test ./...
```

## Product Flow

1. Sign in to LinkNest.
2. Connect X and/or LinkedIn from the channel sidebar.
3. Write a post in the composer.
4. Attach one image if needed.
5. Select the publishing channels.
6. Click Publish now.
7. For LinkedIn, LinkNest publishes through the backend.
8. For X, LinkNest opens X in a new tab with the text prefilled so the user can review and post manually.

## Roadmap Ideas

- Add scheduling.
- Add post history and analytics.
- Add team workspaces.
- Add multi-image and video support.
- Add optional full X API publishing for users who want automatic text + media posting.
- Add saved drafts and reusable templates.
