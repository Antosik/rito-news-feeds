# Teamfight Tactics News

The service generates files that contains the news from official Teamfight Tactics website (https://teamfighttactics.leagueoflegends.com/en-us/news/).

## File URL
The file url should be formed like
```
https://data.rito.news/tft/{locale}/news.{extension}
```

### Available locales
- English (NA) (`en-us`)
- English (EUW) (`en-gb`)
- Deutsch (`de-de`)
- Español (EUW) (`es-es`)
- Français (`fr-fr`)
- Italiano (`it-it`)
- English (OCE) (`en-au`)
- Polski (`pl-pl`)
- Русский (`ru-ru`)
- Ελληνικά (`el-gr`)
- Română (`ro-ro`)
- Magyar (`hu-hu`)
- Čeština (`cs-cz`)
- Español (LATAM) (`es-mx`)
- Português (`pt-br`)
- Türkçe (`tr-tr`)
- 한국어 (`ko-kr`)
- 日本語 (`ja-jp`)
- English (SG) (`en-sg`)
- English (PH) (`en-ph`)
- Tiếng Việt (`vi-vn`)
- ภาษาไทย (`th-th`)
- 繁體中文 (`zh-tw`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS with italian locale - https://data.rito.news/tft/it-it/news.rss
- Raw news data with korean locale - https://data.rito.news/tft/ko-kr/news.json
