FROM golang:1.24.1

WORKDIR /app

COPY elevatorserver-simulation .
COPY simulator.con .

COPY project ./project


