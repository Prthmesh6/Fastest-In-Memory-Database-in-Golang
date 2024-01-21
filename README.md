### Fastest-In-Memory-Database-in-Golang
Inmemory cache/Database in Golang.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Overview

This is an in memory cache which stores values in key-value format. This is similar to redis and I have made it compliant with Redis Client. This can be used with Redis CLI.
I always wanted to create a database to understand how the things actually work in backend and hence I am creating this. 

## How to use ?

1) You should have redis and Go installed in your system, then clone the repository
2) Open 2 terminals in your IDE, One for running our database server other for redis-cli
3) On root folder run "go run ." Make sure nothins else running on port 6379
4) Use "redis-cli" command to run redis client
5) Now on Redis-cli you can execute the redis commands directly.

<div id="header">
  <img src="https://private-user-images.githubusercontent.com/74038190/238200839-9c351cb9-c9a2-4b20-8420-e96b8331a53b.gif?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3MDU4NjA5NzksIm5iZiI6MTcwNTg2MDY3OSwicGF0aCI6Ii83NDAzODE5MC8yMzgyMDA4MzktOWMzNTFjYjktYzlhMi00YjIwLTg0MjAtZTk2YjgzMzFhNTNiLmdpZj9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFWQ09EWUxTQTUzUFFLNFpBJTJGMjAyNDAxMjElMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjQwMTIxVDE4MTExOVomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPWRhOTNmMWZiMmE2NzU4ZTI5ZDE4YTAzYTUyNDdjYTk4ODNjODY1NjMwZmViZjQ4MjdkNjJhNWI0MWJhZWU3N2YmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JmFjdG9yX2lkPTAma2V5X2lkPTAmcmVwb19pZD0wIn0.DIxwRhjMdFvdqxpOAMGmIfo_flae9FCbZtVUNgZjcrc" width="100"/>
</div>

1) This database works superfast since it is InMemory Database. As it is InMemory, server crash can cause data loss. But as I love to solve such problems I have added logic which avoid data loss, I have implemented [AOF](https://redis.io/docs/management/persistence/)
2) For Parsing UnParsing the data we haven't used JSON, we have used [RESP](https://redis.io/docs/reference/protocol-spec/) that is Redis serialization protocol specification you can read more about it on Redis Website, I have implemented it in code. 
3) To be continued...