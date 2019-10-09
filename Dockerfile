FROM golang:1.13.1-alpine3.10

# Need git for dep
RUN apk add --no-cache git

# Need dep to get dependencies
RUN go get -u github.com/golang/dep/cmd/dep

# Copy current dir (outside docker) to the proper directory (inside docker)
COPY . src/github.com/etzelm/consistent-graph-store-api/

# Change current directory
WORKDIR src/github.com/etzelm/consistent-graph-store-api/

# Get dependencies
RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o graphstore .

# # Install assignment/code
# RUN go install .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=0 /go/src/github.com/etzelm/consistent-graph-store-api/graphstore .
CMD ["./graphstore"]  

EXPOSE 50051
EXPOSE 8080