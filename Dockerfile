FROM jekyll/jekyll as jekyllbuilder

COPY --chown=1000:1000 . /site

WORKDIR /site

RUN jekyll build --trace


FROM cupcakearmy/static

COPY --from=jekyllbuilder /site/_site /srv