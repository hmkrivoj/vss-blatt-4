FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o main main.go

FROM iron/go
COPY --from=builder /app/main /app
EXPOSE 52000-53000
ENTRYPOINT [ "/app/main" ]
