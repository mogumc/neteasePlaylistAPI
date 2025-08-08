FROM gcr.io/distroless/static-debian12
WORKDIR /app

COPY devapi /app/devapi
COPY static /app/static

EXPOSE 15967

CMD ["/app/devapi"]
