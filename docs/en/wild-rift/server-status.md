# Wild Rift Server status

The service generates files that contains theserver status message from Wild Rift server status page (https://status.riotgames.com/wildrift).

## File URL
The file url should be formed like
```
https://rito-news.iamantosik.me/wr/{locale}/status.{server}.{extension}
```

### Available servers
- BR (`br`)
- EU (`eu`)
- JP (`jp`)
- KR (`kr`)
- LATAM (`latam`)
- MEI (`mei`)
- NA (`na`)
- RU (`ru`)
- SEA (`sea`)

### Available locales
- English (`en-us`)
- Deutsch (`de-de`)
- English (`en-gb`)
- Español (LATAM) (`es-mx`)
- Español (EU) (`es-es`)
- Français (`fr-fr`)
- Indonesian (`id-id`)
- Italiano (`it-it`)
- 日本語 (`ja-jp`)
- 한국어 (`ko-kr`)
- Polski (`pl-pl`)
- Melayu (`ms-my`)
- Português (BR) (`pt-br`)
- Русский (`ru-ru`)
- ภาษาไทย (`th-th`)
- Türkçe (`tr-tr`)
- Tiếng Việt (`vi-vn`)
- 马来简体中文 (`zh-my`)
- 繁體中文 (`zh-tw`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS for north america server status messages with english locale - https://rito-news.iamantosik.me/wr/en-us/status.na.rss
- Raw server status data for europe server with italian locale - https://rito-news.iamantosik.me/wr/it-it/status.eu.json
