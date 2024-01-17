FROM golang:1.20-alpine

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /go-stock

ENTRYPOINT ["/go-stock"]

# docker build . --tag go-stock:0.08
# docker save go-stock:0.08 | gzip > go-stock_latest.tar.gz
# changed env to IP of postgres container 
# docker load --input go-stock_latest.tar.gz
# docker run --env-file ./docker.env --net=host go-stock:0.08
# since linux, can use host mode!
# copy public.prices (id, created_at, updated_at, type, price, stock_id,received) FROM '/data/prices.csv' DELIMITER ',' CSV QUOTE '"' ESCAPE '''';
# copy public.stocks (id, name, symbol, created_at, updated_at, poll) FROM '/data/stocks.csv' DELIMITER ',' CSV QUOTE '"' ESCAPE '''';