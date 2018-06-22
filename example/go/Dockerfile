FROM golang

# Set environment variables
ENV PATH /go/bin:$PATH

# Cd into the api code directory
WORKDIR /go/src/github.com/miguelmota/c3

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/miguelmota/c3

# Fetch dependencies
#RUN go get  ./...

#EXPOSE 3333

CMD go run /go/src/github.com/miguelmota/c3/example/go/main.go
