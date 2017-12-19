FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY FinteqGCEMonitor /app/
ENTRYPOINT ["/app/FinteqGCEMonitor"]
EXPOSE 8001