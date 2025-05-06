FROM golang:1.24-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN if [ ! -f "go.mod" ]; then \
        make init; \
    fi

RUN make build

# # #

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/otus_social_network /app

RUN chmod +x /app/otus_social_network

CMD ["sh", "-c", "./otus_social_network"]
