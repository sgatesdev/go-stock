# go-stock

![License badge](https://img.shields.io/badge/license-MIT-green)

## Description

go-stock is a simple app for streaming real-time stock data and performing basic time-series analysis. The goal is to leverage the power of FinnHub.io's [free stock API](https://finnhub.io/) to provide a simple, easy to use streaming stock dashboard. This application is a personal project and is designed solely for my use. The application utilizes Go on the backend and TypeScript React on the frontend.

The backend is written in Go. During trading hours, the scheduler polls for stock price data on an interval. When data is received, the price updates are handled as follows:

1. Store price updates for historical analysis
2. Stream price updates to the frontend dashboard, if connected

The backend permits web socket connections on a designated endpoint. Once the connection is established, the backend tracks which stocks have active web socket connections. As price updates come in, the updates are send to each web socket connection by stock. This is to ensure (limited) scalability. The backend can safely handle multiple websocket connections, efficiently broadcasting updates for overlapping or different stocks. 

There are CRUD endpoints for stocks and a GET endpoint for price data.

## Deployment

The application is deployed on a home server and fetches data autonomously when markets are open. It can run in a Docker container or as a standalone Go executable. It runs as a Docker container on the server. 

Scripts are included in the repo to facilitate deployment and release. All data is stored in a Postgres database that is also running in a Docker container on the same Ubuntu server.