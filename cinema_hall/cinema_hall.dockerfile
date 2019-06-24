FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o cinema_hall/main cinema_hall/main.go

FROM iron/go
COPY --from=builder /app/cinema_hall/main /app/cinema_hall
EXPOSE 52000-53000
ENTRYPOINT [ "/app/cinema_hall/main" ]
