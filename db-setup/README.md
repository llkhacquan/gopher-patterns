# DB Setup Pattern

## Problem

Setting up local PostgreSQL for development is repetitive and error-prone. Developers need a one-command solution to spin up a reliable database instance.

## Solution

Docker-based PostgreSQL setup with intelligent Makefile:
- Starts PostgreSQL only if not already running  
- Uses standard defaults (postgres/password/localhost:5432)
- Persistent data across restarts
- Health checks and automatic readiness detection

## Quick Start

```bash
make db            # Start PostgreSQL and ensure postgres database exists
```

## Database Credentials

- Host: localhost:5432
- User: postgres
- Password: password  
- Database: postgres

## Customization

Edit `docker-compose.yml` to change defaults or add SQL files to `init-scripts/` directory for automatic initialization.

## When to Use

Perfect for local development and testing. Avoid for production deployments.