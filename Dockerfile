# Build the go application into a binary
FROM golang:alpine as builder
WORKDIR /under-maintenance
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
RUN echo "Under maintenance" >> under-maintenance.html

# Run the binary on an empty container
FROM scratch
COPY --from=builder /under-maintenance/main .
COPY --from=builder /under-maintenance/under-maintenance.html .
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/main"]
