# Go Server Dockerfile
#FROM golang:1.20 AS builder
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . .
#RUN go build -o main .
#
#FROM alpine:latest
#WORKDIR /root/
#COPY --from=builder /app/main .
#COPY packs.db .
#EXPOSE 8080
#CMD ["./main"]

# Stage 1: Build the React UI
FROM node:18 AS ui-build
WORKDIR /app
COPY ./UI/package*.json ./
RUN npm install
COPY ./UI ./
RUN npm run build

# Stage 2: Build the Go server
FROM golang:1.20 AS server-build
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
COPY --from=ui-build /app/build ./static
RUN go build -o main .

# Final Stage: Create the runtime container
FROM alpine:latest
WORKDIR /root/
COPY --from=server-build /app/main .
COPY --from=server-build /app/static ./static
EXPOSE 8080
CMD ["./main"]