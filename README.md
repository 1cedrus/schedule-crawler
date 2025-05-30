## Inspiration

- This project was created based on the inspiration from
  [wttr.in](https://github.com/chubin/wttr.in) and the inconvenient way to
  access my schedule at my college's site.
- A cool way to view my classes' schedule in the terminal

## Demo

![Demo](./assets/demo.gif)

## Feature

- Using credentials that would be passed through using `Basic Auth Header`,
  crawl and show the schedule using ASCII.
- Auto-caching requests' responses.

## How to run

```bash
$ git clone https://github.com/1cedrus/no.name
$ cd no.name
$ go run .
```

## How to use

```bash
# Current week
$ curl -u "your_password:your_password" localhost
# Last week
$ curl -u "your_password:your_password" localhost/1
# Next week
$ curl -u "your_password:your_password" localhost/-1
# Exam schedule
$ curl -u "your_password:your_password" localhost/exam
```
