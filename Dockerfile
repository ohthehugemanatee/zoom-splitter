FROM jrottenberg/ffmpeg:4.1
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
