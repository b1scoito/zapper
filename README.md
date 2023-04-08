# zapper

redirects stuff from groups

## usage

```bash
zapper -sends (jid's, separated by comma) -receives (jid's)
```

with docker

```bash
docker run -it zapper -sends (jid's, separated by comma) -receives (jid's)
```

## building

with docker

```bash
docker build -t zapper .
```

with garble

```bash
go install mvdan.cc/garble@latest
./compile.sh or compile.cmd depending on your OS
```
