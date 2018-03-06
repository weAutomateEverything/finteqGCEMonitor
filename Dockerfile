FROM alpine:3.6
WORKDIR /app
# Now just add the binary
COPY finteqGCEMonitor /app/
ENTRYPOINT ["/app/finteqGCEMonitor"]
EXPOSE 8001