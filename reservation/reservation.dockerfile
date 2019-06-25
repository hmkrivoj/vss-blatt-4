FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN cd reservation && go build -o main

FROM iron/go
COPY --from=builder /app/reservation/main /app/reservation
EXPOSE 52000-53000
ENTRYPOINT [ "/app/reservation/main" ]
