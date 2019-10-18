# Endpoints

## Whoami

```bash
$ curl localhost:3000/whoami
[
    "127.0.0.1/8",
    "::1/128",
    "172.18.0.4/16",
    "fe80::42:acff:fe12:4/64"
]
```

## Metrics

Exposes [Prometheus](https://prometheus.io/) Metrics.

## Read todo's

```bash
$ curl http://localhost:3000/todo
[
    "Eat",
    "Sleep",
    "Code",
    "Repeat"
]

```

## Insert todo

```bash
$ curl -XPOST http://localhost:3000/todo/Hello
[
  "Eat",
  "Sleep",
  "Code",
  "Repeat",
  "Hello"
]
```

## Delete todo

```bash
$ curl -XDELETE http://localhost:3000/todo/Hello
[
  "Eat",
  "Sleep",
  "Code",
  "Repeat",
]
```

## Health endpoint

```bash
$ curl http://localhost:3000/health
{
    "redis-master-0": "ok",
    "redis-slave-0": "ok",
    "self": "ok"
}
```
