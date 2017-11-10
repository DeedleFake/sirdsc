FROM scratch

EXPOSE 8080

COPY ./web /web
COPY ./interface/build /pub

ENTRYPOINT ["/web", "-root", "/pub"]
