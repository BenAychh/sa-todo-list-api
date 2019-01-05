FROM golang as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get -d ./...
CMD ["go", "test", "-v"]