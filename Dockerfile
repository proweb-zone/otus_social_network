FROM golang:1.24-alpine

RUN apk add --no-cache make

WORKDIR /otus

COPY . .

# RUN if [ ! -f "go.mod" ]; then \
#         make init; \
#     fi

# RUN make build

RUN chmod +x /otus/app/build/*

CMD ["sh", "-c", "/otus/app/build/otus_social_network"]
