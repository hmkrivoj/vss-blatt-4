FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o reservation/main reservation/main

FROM iron/go
COPY --from=builder /app/reservation/main /app/reservation
EXPOSE 52000-53000
ENTRYPOINT [ "/app/reservation/main" ]
