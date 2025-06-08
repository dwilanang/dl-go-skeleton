# Builder
FROM golang:1.23

ENV FILE_LOC=./.dev/.air.toml

RUN apt update && apt upgrade -y && \
    apt install -y git \
    make openssh-client

WORKDIR /go/src/ppms

COPY go.* ./
RUN go mod download && go mod verify

COPY . .

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh v1.42.0 && cp ./bin/air /bin/air

CMD air -c $FILE_LOC