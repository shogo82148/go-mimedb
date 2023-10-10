# go-mimedb

go-mimedb is a drop-in replacement for [the mime package](https://pkg.go.dev/mime) that is independent of external files.

## Motivation

[The mime package](https://pkg.go.dev/mime) uses the local system's MIME-info database or mime.types file(s).

```
/usr/local/share/mime/globs2
/usr/share/mime/globs2
/etc/mime.types
/etc/apache2/mime.types
/etc/apache/mime.types
```

Therefore, it returns different results depending on the execution environment.

go-mimedb has a large build-in table and does not depend on these files.

## Data Source of MIME-info database

MIME-info database comes from [jshttp/mime-db](https://github.com/jshttp/mime-db).
