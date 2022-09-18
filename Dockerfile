FROM golang:1.19
COPY . /auth-srvc
WORKDIR /auth-srvc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN apt update --no-install-recommends
RUN apt install --no-install-recommends -y protobuf-compiler 
RUN make proto
RUN go mod download
RUN go build -o auth-srvc cmd/main.go
CMD ./auth-srvc
