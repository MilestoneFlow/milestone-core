FROM --platform=linux/amd64 public.ecr.aws/docker/library/golang:1.21

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app

CMD ["app"]

EXPOSE 3333
