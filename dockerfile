FROM golang:alpine3.16

# Install git.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

ENV GIN_MODE=release
ENV PORT=3000

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .
#COPY db-config.sql /docker-entrypoint-initdb.d/
# Download all the dependencies 
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN go build -o /build

# Expose port 3000 to the outside world
EXPOSE ${PORT}

# Run the executable
CMD [ "go", "run", "." ]
