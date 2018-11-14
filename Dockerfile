FROM scratch

ENV PORT 8000
EXPOSE $PORT

COPY ./serv4kub /serv4kub
CMD ["/serv4kub"]
