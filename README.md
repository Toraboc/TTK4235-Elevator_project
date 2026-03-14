# Sanntid prosjekt

Dette er kode til TTK4145 prosjekt.

- [ ] Exercise 1
- [ ] Exercise 2
- [ ] Exercise 3
- [ ] Exercise 4
- [ ] Exercise 5

## Program flagg
Det finnes ulike typer flagg som kan brukes til programmet. Bruk `-h` flagget for å få hjelp. Legg også merke til at simulatoren også kan lytte på ulike porter med flagget `--port 5432`

## Start koden på andre maskiner

Bruk scriptet til å starte heiskoden på flere maskiner. Gjør dette ved å spesifisere det siste tallet i ip-adressen til maskinene man skal deployere til.

```bash
./deploy_and_run.sh <number> [number] [number]
```

For eksempel slik

```bash
./deploy_and_run.sh 35 36 37
```

Kan også brukes med en fast mode. Da vil kun koden kopieres og kjøres.

```bash
./deploy_and_run.sh -f 35 36 37
```

## TODO

- [ ] Fikse de siste bugsene
- [ ] Refactore koden litt
    - [ ] Fjerne TODO
    - [ ] Se over alle mutex
    - [ ] Endre hvordan kanalene settes opp
    - [ ] Fjerne ubrukte funksjoner
    - [ ] Se over variabel navn
    - [ ] Se over side effects, `getMyId()` og `elevio` greier
- [ ] Teste alt sammen
    - [ ] Pakketap
- [ ] Supervisor

## Docker

This can be used to run multiple simulators that will talk to eachother.

```bash
docker compose up --build
```

To access on of the elevators run
```bash
docker exec -it sanntid55-simulator-1-1 tmux attach
```
or replace the second last `1` to `2` or `3`.

Use `ctrl-b` then `d` to detach from the tmux session.
