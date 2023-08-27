FROM docker.io/golang:1.21 AS build
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go vet -v ./src/
RUN CGO_ENABLED=0 go build -v -o /go/bin/app ./src/

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=build /go/bin/app /app
CMD [ "/app/app" ]