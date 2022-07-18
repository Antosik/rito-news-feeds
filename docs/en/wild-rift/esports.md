# Wild Rift Esports news

The service generates files that contains the esports news from official Wild Rift esports website (https://wildriftesports.com/en-us/news).

## File URL
The file url should be formed like
```
https://rito-news.iamantosik.me/wr/{locale}/esports.{extension}
```

### Available locales
- English (North America) (`en-us`)
- English (Europe) (`en-gb`)
- English (Oceania) (`en-au`)
- Čeština (`cs-cz`)
- Deutsch (`de-de`)
- Ελληνικά (`el-gr`)
- Español (España) (`es-es`)
- Español (Latinoamérica) (`es-mx`)
- Français (`fr-fr`)
- Italiano (`it-it`)
- Polski (`pl-pl`)
- Português (Brasil) (`pt-br`)
- Română (`ro-ro`)
- Русский (`ru-ru`)
- Türkçe (`tr-tr`)
- 日本語 (`ja-jp`)
- 한국어 (`ko-kr`)
- 繁體中文 (`zh-tw`)
- ภาษาไทย (`th-th`)
- English (Philippines) (`en-ph`)
- English (Singapore) (`en-sg`)
- Indonesian (`id-id`)
- Tiếng việt (`vi-vn`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS with italian locale - https://rito-news.iamantosik.me/wr/it-it/esports.rss
- Raw news data with korean locale - https://rito-news.iamantosik.me/wr/ko-kr/esports.json
