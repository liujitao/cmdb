FROM golang:1.15.5-alpine

RUN mkdir /app
ADD . /app

WORKDIR /app

RUN go env -w GOPROXY="https://goproxy.cn,https://goproxy.io,direct" && go env -w GO111MODULE="on" && go env -w GOSUMDB="off"
RUN go mod download
#RUN go build -o main . 

EXPOSE 8000
#CMD ["/app/main"]
CMD ["go", "run", "/app/main.go"]