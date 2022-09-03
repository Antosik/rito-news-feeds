# League of Legends News

The service generates files that contains the news from official League of Legends website (https://leagueoflegends.com/en-us/news/).

## File URL
The files url looks like
```
https://rito-news.iamantosik.me/lol/{locale}/news.{extension}
```

### Available locales
- English (NA) (`en-us`)
- English (EUW) (`en-gb`)
- Deutsch (`de-de`)
- Español (EUW) (`es-es`)
- Français (`fr-fr`)
- Italiano (`it-it`)
- English (EUNE) (`en-pl`)
- Polski (`pl-pl`)
- Ελληνικά (`el-gr`)
- Română (`ro-ro`)
- Magyar (`hu-hu`)
- Čeština (`cs-cz`)
- Español (LATAM) (`es-mx`)
- Português (`pt-br`)
- 日本語 (`ja-jp`)
- Русский (`ru-ru`)
- Türkçe (`tr-tr`)
- English (OCE) (`en-au`)
- 한국어 (`ko-kr`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS for NA news - https://rito-news.iamantosik.me/lol/en-us/news.rss
- Raw news data for Korean server - https://rito-news.iamantosik.me/lol/ko-kr/news.json