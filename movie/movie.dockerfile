FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o movie/main movie/main.go

FROM iron/go
COPY --from=builder /app/movie/main /app/movie
EXPOSE 52000-53000
ENTRYPOINT [ "/app/movie/main" ]
