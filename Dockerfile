FROM golang:1.24.2-alpine3.21 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
ARG SERVICE
RUN go build -o /bin/${SERVICE} cmd/${SERVICE}/main.go

FROM alpine:latest AS final
RUN apk update && apk add --no-cache ca-certificates
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
ARG SERVICE
ENV SERVICE=${SERVICE} 

COPY --from=build /bin/${SERVICE} .
# RUN chmod +x ./${SERVICE}
ENTRYPOINT ["sh", "-c", "./${SERVICE}"]