FROM golang:alpine 
WORKDIR /GObotDias
COPY . .
RUN ["go", "mod", "verify"]
RUN ["go", "mod", "tidy"]
RUN ["go", "build", "-o", "main", "main.go"]
CMD [ "./main" ]