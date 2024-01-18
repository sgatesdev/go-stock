# go-stock

![License badge](https://img.shields.io/badge/license-MIT-green)

## Description

go-stock is a simple app for streaming real-time stock data and performing basic time-series analysis. The goal is to leverage the power of FinnHub.io's [free stock API](https://finnhub.io/) to provide a simple, easy to use streaming stock dashboard. This application is a personal project and is designed solely for my use. The application utilizes Go on the backend and TypeScript React on the frontend.

The backend uses Go to execute the fetching of live price data on a set interval. Go threads are used to concurrently poll the API for fresh stock price data. The backend does two things: 

1. Fetches and stores new price data for historical analysis
2. Streams price updates to the frontend dashboard, if connected

Websockets are used to stream price updates. React Echarts is used to provide charting capability - both on the live dashboard page and on the historical analysis page. The historical analysis page permits me to specify a time interval I wish to chart, so I can see price fluctuations during that time period. I also leverage Echarts to provide some basic information about the data itself, such as min and max prices during the chosen time period.

There are also various CRUD endpoints for managing which stocks I want to poll, and for fetching price information by stock. 

## Deployment

The application is deployed on a home server and fetches data autonomously when markets are open. It can run in a Docker container or as a standalone Go executable. Scripts are included in the repo to make deployment and release easier. All data is stored in a Postgres database that is also running in a Docker container on the same Ubuntu server.