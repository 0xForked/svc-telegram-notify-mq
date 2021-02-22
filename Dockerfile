# Start from golang:1.12-alpine base image
# builder
FROM golang:1.15-alpine as builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk --update --no-cache add bash

# Set the Current Working Directory inside the container
WORKDIR /service

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependancies to vendor folder.
RUN go mod vendor

# Build the Go app
RUN go build -o telegram_notify_svc .

## Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --no-cache --update add tzdata && \
    mkdir /service

WORKDIR /service

COPY --from=builder /service/telegram_notify_svc /service
COPY --from=builder /service/.env /service

# Run the executable
CMD /service/telegram_notify_svc