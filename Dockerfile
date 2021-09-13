FROM golang:1.15.2-alpine

WORKDIR $GOPATH/src/mysql-user-db-go-interface
ARG SERVER_PORT

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

RUN apk add --update g++
RUN go mod tidy
RUN cd $GOPATH/src/mysql-user-db-go-interface/ && go build main.go

# This container exposes SERVER_PORT to the outside world.
# Check .env for the actual value
EXPOSE $SERVER_PORT

RUN chmod 0766 $GOPATH/src/mysql-user-db-go-interface/scripts/init.sh

# Run the executable
CMD ["./scripts/init.sh"]