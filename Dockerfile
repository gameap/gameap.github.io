FROM jekyll/jekyll as jekyllbuilder

COPY --chown=1000:1000 . /site

WORKDIR /site

RUN jekyll build --trace && \
    mkdir -p _site/en/css _site/ru/css && \
    cp css/*.css _site/en/css/ && \
    cp css/*.css _site/ru/css/


FROM golang:1.25-alpine as gobuilder

WORKDIR /app
COPY server/main.go .
COPY --from=jekyllbuilder /site/_site ./site

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server main.go


FROM scratch

COPY --from=gobuilder /server /server

ENV LANG=en
ENV REDIRECTS=""
ENV PORT=80

EXPOSE 80

CMD ["/server"]
