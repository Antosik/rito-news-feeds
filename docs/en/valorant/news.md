# VALORANT News

The service generates files that contains the news from official VALORANT website (https://playvalorant.com/ru-ru/news/).

## File URL
The file url should be formed like
```
https://data.rito.news/val/{locale}/news.{extension}
```

### Available locales
- English (NA) (`en-us`)
- English (EUW) (`en-gb`)
- Deutsch (`de-de`)
- Español (EUW) (`es-es`)
- Français (`fr-fr`)
- Italiano (`it-it`)
- Polski (`pl-pl`)
- Русский (`ru-ru`)
- Türkçe (`tr-tr`)
- Español (LATAM) (`es-mx`)
- Indonesian (`id-id`)
- 日本語 (`ja-jp`)
- 한국어 (`ko-kr`)
- Português (`pt-br`)
- ภาษาไทย (`th-th`)
- Tiếng việt (`vi-vn`)
- 繁體中文 (`zh-tw`)
- العربية (`ar-ae`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS with italian locale - https://data.rito.news/val/it-it/news.rss
- Raw news data with korean locale - https://data.rito.news/val/ko-kr/news.json
