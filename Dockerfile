FROM golang:1.25.5 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /out/core-backend ./cmd/core-backend

FROM gcr.io/distroless/base-debian12:nonroot

ENV CORE_BACKEND_HOST=0.0.0.0

WORKDIR /app

COPY --from=build /out/core-backend /app/core-backend

EXPOSE 8080

ENTRYPOINT ["/app/core-backend"]
