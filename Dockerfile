# Build the go application into a binary
FROM golang:alpine as builder
WORKDIR /under-maintenance
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Run the binary on an empty container
FROM scratch
COPY --from=builder /under-maintenance/main .
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/main"]
