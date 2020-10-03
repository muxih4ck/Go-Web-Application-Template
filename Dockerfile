FROM golang:1.12.13 
ENV GO111MODULE "on"
ENV GOPROXY "https://goproxy.cn"
WORKDIR /src/Go-Web-Application-Template
COPY . /src/Go-Web-Application-Template
RUN make
EXPOSE 8080
CMD ["./main", "-c", "conf/config.yaml"]
