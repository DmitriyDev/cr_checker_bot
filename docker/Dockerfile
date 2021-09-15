FROM golang:1.15-alpine

ENV CGO_ENABLED=0

ADD ./_source/*.go /app/_build/
WORKDIR /app/

COPY ./_source/go.mod ./
COPY ./_source/go.sum ./
RUN go mod download

RUN go build -o ./cr_checker ./_build/

CMD [ "./cr_checker" ]