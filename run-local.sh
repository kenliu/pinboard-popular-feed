#!/bin/bash

# Exit on error
set -e

# Configuration
DRY_RUN=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --dryrun)
      DRY_RUN=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Build the docker run command
CMD="docker run --rm"

# Add environment variables
CMD+=" -e MASTODON_ACCESS_TOKEN=$MASTODON_ACCESS_TOKEN"
CMD+=" -e MASTODON_SERVER_DOMAIN=$MASTODON_SERVER_DOMAIN"
CMD+=" -e DB_USERNAME=$DB_USERNAME"
CMD+=" -e DB_PASSWORD=$DB_PASSWORD"
CMD+=" -e DB_HOST=$DB_HOST"
CMD+=" -e DB_PORT=$DB_PORT"
CMD+=" -e DB_NAME=$DB_NAME"

# Add the container name
CMD+=" pinboard-popular-feed"

# Add dry run flag if specified
if [ "$DRY_RUN" = true ]; then
  CMD+=" --dryrun"
fi

# Run the command
echo "Running: $CMD"
eval $CMD 