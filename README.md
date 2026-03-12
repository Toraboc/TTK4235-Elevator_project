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

- Hvilke prosesser skal vi ha, nettverk, en heismodulen

## Docker

```bash
podman build --platform linux/x86_64 -t sanntid55 .
podman run --platform linux/x86_64 --rm -it -p 15657:15657 -v ./project:/app/project --name sanntid55 sanntid55 /bin/bash
podman exec -it sanntid55 /bin/bash
```
