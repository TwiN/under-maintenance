# Build the go application into a binary
FROM golang:alpine as builder
WORKDIR /app
ADD . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o under-maintenance .
RUN echo "Under maintenance" >> under-maintenance.html

# Run the binary on an empty container
FROM scratch
COPY --from=builder /app/under-maintenance .
COPY --from=builder /app/under-maintenance.html .
ENV PORT 80
EXPOSE 80
ENTRYPOINT ["/under-maintenance"]
