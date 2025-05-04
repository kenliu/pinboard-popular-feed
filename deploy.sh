#!/bin/bash

# Exit on error
set -e

# Configuration
PROJECT_ID="${GOOGLE_CLOUD_PROJECT:-pinboard-popular-feed-dev}"
JOB_NAME="pinboard-popular-feed"
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
docker build -t ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${JOB_NAME}:${IMAGE_TAG} .

# Push the image to Artifact Registry
echo "Pushing image to Artifact Registry..."
docker push ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${JOB_NAME}:${IMAGE_TAG}

# Prepare environment variables string
ENV_VARS_STRING=""
for key in "${!ENV_VARS[@]}"; do
    ENV_VARS_STRING+="${key}=${ENV_VARS[$key]},"
done
# Remove trailing comma
ENV_VARS_STRING=${ENV_VARS_STRING%,}

# Deploy to Cloud Run Job
echo "Deploying to Cloud Run Job..."

# Check if job exists
if gcloud run jobs describe ${JOB_NAME} --region ${REGION} >/dev/null 2>&1; then
    echo "Updating existing job..."
    gcloud run jobs update ${JOB_NAME} \
        --image ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${JOB_NAME}:${IMAGE_TAG} \
        --region ${REGION} \
        --set-env-vars "${ENV_VARS_STRING}" \
        --max-retries 3 \
        --task-timeout 10m
else
    echo "Creating new job..."
    gcloud run jobs create ${JOB_NAME} \
        --image ${REGION}-docker.pkg.dev/${PROJECT_ID}/${REPOSITORY}/${JOB_NAME}:${IMAGE_TAG} \
        --region ${REGION} \
        --set-env-vars "${ENV_VARS_STRING}" \
        --max-retries 3 \
        --task-timeout 10m
fi

# Set up Cloud Scheduler
echo "Setting up Cloud Scheduler..."
SCHEDULER_NAME="${JOB_NAME}-scheduler"

# Get project number for the service account
PROJECT_NUMBER=$(gcloud projects describe ${PROJECT_ID} --format='value(projectNumber)')

# Check if scheduler exists
if gcloud scheduler jobs describe ${SCHEDULER_NAME} --location ${REGION} >/dev/null 2>&1; then
    echo "Updating existing scheduler..."
    gcloud scheduler jobs update http ${SCHEDULER_NAME} \
        --location ${REGION} \
        --schedule "0 * * * *" \
        --time-zone "America/New_York" \
        --uri "https://${REGION}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${PROJECT_NUMBER}/jobs/${JOB_NAME}:run" \
        --http-method POST \
        --oauth-service-account-email "${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
        --oauth-token-scope "https://www.googleapis.com/auth/cloud-platform"
else
    echo "Creating new scheduler..."
    gcloud scheduler jobs create http ${SCHEDULER_NAME} \
        --location ${REGION} \
        --schedule "0 * * * *" \
        --time-zone "America/New_York" \
        --uri "https://${REGION}-run.googleapis.com/apis/run.googleapis.com/v1/namespaces/${PROJECT_NUMBER}/jobs/${JOB_NAME}:run" \
        --http-method POST \
        --oauth-service-account-email "${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
        --oauth-token-scope "https://www.googleapis.com/auth/cloud-platform"
fi

echo "Deployment complete!" 