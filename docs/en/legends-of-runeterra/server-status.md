# Legends of Runeterra Server status

The service generates files that contains theserver status message from Legends of Runeterra server status page (https://status.riotgames.com/lor).

## File URL
The file url should be formed like
```
https://data.rito.news/lor/{locale}/status.{server}.{extension}
```

### Available servers
- Americas (`americas`)
- Asia-Pacific (`apac`)
- Europe (`europe`)

### Available locales
- English (`en-us`)
- Deutsch (`de-de`)
- Español (EU) (`es-es`)
- Español (LATAM) (`es-mx`)
- Français (`fr-fr`)
- Italiano (`it-it`)
- 日本語 (`ja-jp`)
- 한국어 (`ko-kr`)
- Polski (`pl-pl`)
- Português (BR) (`pt-br`)
- Русский (`ru-ru`)
- ภาษาไทย (`th-th`)
- Türkçe (`tr-tr`)
- Tiếng Việt (`vi-vn`)
- 繁體中文 (`zh-tw`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS for americas server status messages with english locale - https://data.rito.news/lor/en-us/status.americas.rss
- Raw server status data for europe server with italian locale - https://data.rito.news/lor/it-it/status.europe.json
