# Legends of Runeterra News

The service generates files that contains the news from official Legends of Runeterra website (https://playruneterra.com/en-us/news/).

## File URL
The file url should be formed like
```
https://rito-news.iamantosik.me/lor/{locale}/news.{extension}
```

### Available locales
- English (NA) (`en-us`)
- 한국어 (`ko-kr`)
- Français (`fr-fr`)
- Español (EUW) (`es-es`)
- Español (LATAM) (`es-mx`)
- Deutsch (`de-de`)
- Italiano (`it-it`)
- Polski (`pl-pl`)
- Português (`pt-br`)
- Türkçe (`tr-tr`)
- Русский (`ru-ru`)
- 日本語 (`ja-jp`)
- English (SG) (`en-sg`)
- 繁體中文 (`zh-tw`)
- ภาษาไทย (`th-th`)
- Tiếng việt (`vi-vn`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS for English locale news - https://rito-news.iamantosik.me/lor/en-us/news.rss
- Raw news data for Korean locale - https://rito-news.iamantosik.me/lor/ko-kr/news.json
