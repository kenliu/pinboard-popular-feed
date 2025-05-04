#!/bin/bash

# Exit on error
set -e

# Configuration
PROJECT_ID="${GOOGLE_CLOUD_PROJECT:-pinboard-popular-feed-dev}"
SERVICE_NAME="pinboard-popular-feed"
REGION="us-central1"
IMAGE_TAG="latest"
REPOSITORY="pinboard-popular-feed"

# Check for required environment variables
required_vars=(
    "MASTODON_ACCESS_TOKEN"
    "MASTODON_SERVER_DOMAIN"
    "DB_USERNAME"
    "DB_PASSWORD"
    "DB_HOST"
    "DB_PORT"
    "DB_NAME"
)

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: Required environment variable $var is not set"
        exit 1
    fi
done

# Environment variables
declare -A ENV_VARS=(
    ["MASTODON_ACCESS_TOKEN"]="${MASTODON_ACCESS_TOKEN}"
    ["MASTODON_SERVER_DOMAIN"]="${MASTODON_SERVER_DOMAIN}"
    ["DB_USERNAME"]="${DB_USERNAME}"
    ["DB_PASSWORD"]="${DB_PASSWORD}"
    ["DB_HOST"]="${DB_HOST}"
    ["DB_PORT"]="${DB_PORT}"
    ["DB_NAME"]="${DB_NAME}"
)

# Configure Docker to use Artifact Registry
echo "Configuring Docker authentication..."
gcloud auth configure-docker ${REGION}-docker.pkg.dev

# Build the Docker image
echo "Building Docker image..."
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${SERVICE_NAME}:${IMAGE_TAG} .

# Push the image to Artifact Registry
echo "Pushing image to Artifact Registry..."
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${SERVICE_NAME}:${IMAGE_TAG}

# Prepare environment variables string
ENV_VARS_STRING=""
for key in "${!ENV_VARS[@]}"; do
    ENV_VARS_STRING+="${key}=${ENV_VARS[$key]},"
done
# Remove trailing comma
ENV_VARS_STRING=${ENV_VARS_STRING%,}

# Deploy to Cloud Run
echo "Deploying to Cloud Run..."
gcloud run deploy ${SERVICE_NAME} \
    --image ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${SERVICE_NAME}:${IMAGE_TAG} \
    --region ${REGION} \
    --platform managed \
    --allow-unauthenticated \
    --max-instances 1 \
    --set-env-vars "${ENV_VARS_STRING}"

echo "Deployment complete!" 