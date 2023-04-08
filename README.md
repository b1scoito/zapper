# zapper

redirects stuff from groups

## usage

```bash
zapper -sends (jid's, separated by comma) -receives (jid's)
```

with docker

```bash
docker-compose up
# running in the background
docker-compose up -d
```

## building

with garble

```bash
go install mvdan.cc/garble@latest
./compile.sh or compile.cmd depending on your OS
```
