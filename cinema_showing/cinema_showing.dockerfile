FROM obraun/vss-protoactor-jenkins as builder
COPY . /app
WORKDIR /app
RUN go build -o cinema_showing/main cinema_showing

FROM iron/go
COPY --from=builder /app/cinema_showing/main /app/cinema_showing
EXPOSE 52000-53000
ENTRYPOINT [ "/app/cinema_showing/main" ]
