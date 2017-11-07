FROM trimetric-web AS web-build

FROM trimetric-api AS api-build

FROM busybox:glibc

COPY --from=web-build /opt/trimetric/web/dist /opt/trimetric/web/dist
COPY --from=api-build /go/bin/trimetric /opt/trimetric/
WORKDIR /opt/trimetric

EXPOSE 80

VOLUME ["/opt/trimetric"]

CMD ["./trimetric"]