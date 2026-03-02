# Development container for Go + Templ
FROM golang:1.25-alpine

RUN apk add --no-cache gcc musl-dev git

# Install templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Source code is mounted via volume
CMD ["sh"]
