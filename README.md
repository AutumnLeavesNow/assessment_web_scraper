# Word count

This Go program generates linkks (to be changed for actual link input), fetches content from web pages, and processes according to rules

## Running tests ##
```bash
git clone https://github.com/AutumnLeavesNow/assessment_web_scraper.git
cd assessment_web_scraper
go mod download
go test -v ./...

```
## Configuration ##
You can change values in config.json to alter nmber of workers, or add additional parsing rules for web pages


## Future considerations
- dynamic worker pool
- saving data to db
- separate config workers from parsing config
- for each worker add separape dlqueue for avoiding bottlenecks

