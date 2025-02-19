# jsoncdetox

## NAME

jsoncdetox - Convert JSONC into JSON

## SYNOPSIS

```
jsoncdetox [ -i input_file ] [ -o output_file ] [ -m ]
```

## DESCRIPTION

'jsoncdetox' converts [JSON with Comments](https://code.visualstudio.com/docs/languages/json#_json-with-comments) into valid JSON by removing comments, trailing commas.

## OPTIONS

### -i, --input

Read from the given file path, If omitted, read from stdin.

### -o, --output

Write output to the given file path, If omitted, write to stdout.

### -m, --minify

Output minified JSON instead of pretty-print.

## EXAMPLE

```
$ cat sample.jsonc
{
  "foo": "bar", // foobar
  /*         _                           __     __
   *        (_)________  ____  _________/ /__  / /_____  _  __
   *       / / ___/ __ \/ __ \/ ___/ __  / _ \/ __/ __ \| |/_/
   *      / (__  ) /_/ / / / / /__/ /_/ /  __/ /_/ /_/ />  <
   *   __/ /____/\____/_/ /_/\___/\__,_/\___/\__/\____/_/|_|
   *  /___/
   */
  "baz": {
   "qux": 1,
   "quux": 2,
   "quuux": 3,
   "quuuux": 4,
   "quuuuux": 5,
   "quuuuuux": 6,
  },
  "corge": [1, 2, 3, 4, 5, 6],
}

$ jsoncdetox -i sample.jsonc
{
  "foo": "bar",
  "baz": {
    "qux": 1,
    "quux": 2,
    "quuux": 3,
    "quuuux": 4,
    "quuuuux": 5,
    "quuuuuux": 6
  },
  "corge": [
    1,
    2,
    3,
    4,
    5,
    6
  ]
}
```

## RUN DIRECTLY

```
go run github.com/drillbits/jsoncdetox@latest
```

## INSTALLATION

```
go install github.com/drillbits/jsoncdetox@latest
```
