FROM scratch

COPY podexec /podexec
ADD static /static 
ADD templates /templates 
WORKDIR /

ENTRYPOINT ["/podexec"]
