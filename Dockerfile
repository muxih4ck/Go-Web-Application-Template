FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build --mod vendor -o main . 
CMD ["/app/main"]