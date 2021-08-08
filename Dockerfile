ARG APP=shorturl
FROM alpine
ARG APP
ENV APP=${APP}
ENV WORKDIR=/data
COPY bin/shorturl ${WORKDIR}
WORKDIR ${WORKDIR}
CMD ["sh", "-c", "./${APP}"]
