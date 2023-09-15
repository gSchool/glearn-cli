# Start with a light-weight image and install BASH
FROM alpine:3.18.3
RUN apk add --no-cache bash

# Set current directory
WORKDIR /app

COPY test.sh .

# For local testing only. Learn will provide the actual submission which overrides this
COPY submission.txt .
