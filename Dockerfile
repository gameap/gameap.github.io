FROM jekyll/jekyll as jekyllbuilder

COPY --chown=1000:1000 . /site

WORKDIR /site

RUN jekyll build --trace


FROM golang:1.21-alpine as gobuilder

WORKDIR /app
COPY server/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=gobuilder /server /server
COPY --from=jekyllbuilder /site/_site /srv

ENV LANG=en
ENV REDIRECTS=""
ENV PORT=80
ENV SITE_ROOT=/srv

EXPOSE 80

CMD ["/server"]
