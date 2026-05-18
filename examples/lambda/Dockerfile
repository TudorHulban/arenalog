# Stage 1: Build the Go binary
FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go

# Stage 2: Run in the actual Lambda environment
FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /app/bootstrap /var/runtime/bootstrap
CMD [ "bootstrap" ]