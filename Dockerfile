# Build React UI
FROM node:18 AS ui-build
WORKDIR /app
COPY ./UI/package*.json ./
RUN npm install
COPY ./UI ./
RUN npm run build

# Build Go server
FROM golang:1.24 AS server-build
WORKDIR /go/src/packCalculator/server
COPY ./server/go.mod ./server/go.sum ./
RUN go mod download
COPY ./server ./
COPY --from=ui-build /app/build ./cmd/static
RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev
RUN go build -o main ./cmd

# Final image
FROM debian:bookworm-slim
WORKDIR /root/
COPY --from=server-build /go/src/packCalculator/server/main .
COPY --from=server-build /go/src/packCalculator/server/cmd/static ./static
RUN chmod +x ./main
EXPOSE 8080
CMD ["./main"]