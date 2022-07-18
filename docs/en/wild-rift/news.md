# Wild Rift News

The service generates files that contains the news from official Wild Rift website (https://wildrift.leagueoflegends.com/en-us/news/).

## File URL
The file url should be formed like
```
https://rito-news.iamantosik.me/wr/{locale}/news.{extension}
```

### Available locales
- English (NA) (`en-us`)
- English (EUW) (`en-gb`)
- Français (`fr-fr`)
- Deutsch (`de-de`)
- Español (EUW) (`es-es`)
- Italiano (`it-it`)
- Polski (`pl-pl`)
- Русский (`ru-ru`)
- Türkçe (`tr-tr`)
- Indonesian (`id-id`)
- Malaysian (`ms-my`)
- Português (`pt-br`)
- 日本語 (`ja-jp`)
- 한국어 (`ko-kr`)
- 繁體中文 (`zh-tw`)
- ภาษาไทย (`th-th`)
- Tiếng việt (`vi-vn`)
- Español (latam) (`es-mx`)
- English (SG) (`en-sg`)
- العربية (`ar-ae`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS with italian locale - https://rito-news.iamantosik.me/wr/it-it/news.rss
- Raw news data with korean locale - https://rito-news.iamantosik.me/wr/ko-kr/news.json
