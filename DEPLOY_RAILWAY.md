# Deploying PomoHub Backend to Railway

This guide will help you deploy the Go backend, PostgreSQL database, and Redis (Valkey) to Railway.app.

## Prerequisites
- A [Railway](https://railway.app/) account.
- [GitHub CLI](https://cli.github.com/) or Git installed.

## Steps

### 1. Create a New Project on Railway
1.  Log in to Railway.
2.  Click **"New Project"**.
3.  Choose **"Provision PostgreSQL"**.
4.  Add another service: **"Provision Redis"**.

### 2. Deploy the Go Backend
1.  Connect your GitHub repository to Railway or use the Railway CLI to deploy from local.
2.  If connecting GitHub:
    *   Select your repo.
    *   Go to **Settings** -> **Root Directory** and set it to `/backend`.
    *   Railway should automatically detect the `Dockerfile` in the `/backend` folder.

### 3. Configure Environment Variables
Go to the **Variables** tab of your Go service and add the following:

*   `PORT`: `8080` (Railway sets this automatically usually, but good to have)
*   `DATABASE_URL`: `${{PostgreSQL.DATABASE_URL}}` (Reference the Postgres service)
*   `REDIS_URL`: `${{Redis.REDIS_URL}}` (Reference the Redis service)
*   `JWT_SECRET`: Generate a strong random string (e.g., `openssl rand -hex 32`).

### 4. Verify Deployment
1.  Wait for the build to finish.
2.  Railway will generate a public URL (e.g., `https://pomohub-production.up.railway.app`).
3.  Test the health endpoint: `GET https://<your-url>/api/v1/health`.

## Database Migrations
The application runs `AutoMigrate` on startup, so tables will be created automatically when the service starts successfully.

## Security Notes
- Ensure your `JWT_SECRET` is strong and not committed to git.
- The `DATABASE_URL` provided by Railway includes the username/password, so it's secure within the internal network.
