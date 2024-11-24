# yuriposting
Take two! This is a program that fetches a single random post from Danbooru and uploads it to Mastodon.

## Run

1. Copy and configure `config.example.json`
2. Build with `go build cmd/yuriposting/main.go`
3. Run with `./main` (or `.\main.exe` on Windows)

## Automatically post every N minutes

It runs with a cronjob. Here's my cronjob as a reference:

`0,30 * * * * cd /home/vanilla/projects/yuriposting && ./main >> ~/yuriposting.log 2>&1`
