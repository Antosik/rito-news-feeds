# VALORANT Server status

The service generates files that contains theserver status message from VALORANT server status page (https://status.riotgames.com/valorant).

## File URL
The file url should be formed like
```
https://rito-news.iamantosik.me/val/{locale}/status.{server}.{extension}
```

### Available servers
- Asia Pacific (`ap`)
- Brazil (`br`)
- Europe (`eu`)
- Korea (`kr`)
- Latin America (`latam`)
- North America (`na`)
- PBE (`pbe`)


### Available locales
- English (`en-us`)
- العربية (`ar-ae`)
- Deutsch (`de-de`)
- Español (EU) (`es-es`)
- Español (LATAM) (`es-mx`)
- Français (`fr-fr`)
- Indonesian (`id-id`)
- Italiano (`it-it`)
- 日本語 (`ja-jp`)
- 한국어 (`ko-kr`)
- Polski (`pl-pl`)
- Português (BR) (`pt-br`)
- Русский (`ru-ru`)
- Türkçe (`tr-tr`)
- ภาษาไทย (`th-th`)
- Tiếng Việt (`vi-vn`)
- 繁體中文 (`zh-tw`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS for north america server status messages with english locale - https://rito-news.iamantosik.me/val/en-us/status.na.rss
- Raw server status data for europe server with italian locale - https://rito-news.iamantosik.me/val/it-it/status.eu.json
