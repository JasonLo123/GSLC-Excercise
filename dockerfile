FROM golang:1.22

RUN mkdir test
WORKDIR /test

RUN git clone https://github.com/JasonLo123/GSLC-Excercise.git

WORKDIR /test/GSLC-Excercise/app

RUN go run Server.go