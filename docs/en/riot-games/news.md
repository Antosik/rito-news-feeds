# Riot Games News

The service generates files that contains the news from official Riot Games website (https://www.riotgames.com/en/news).

## File URL
The files url looks like
```
https://data.rito.news/riotgames/{locale}/news.{extension}
```

### Available locales
- English (NA) - (`en`)
- Indonesian - (`id`)
- Malaysian - (`ms`)
- Português - (`pt-br`)
- Čeština - (`cs`)
- Français - (`fr`)
- Deutsch - (`de`)
- Ελληνικά - (`el`)
- Magyar - (`hu`)
- Italiano - (`it`)
- 日本語 - (`ja`)
- 한국어 - (`ko`)
- Español (LATAM) - (`es-419`)
- Polski - (`pl`)
- Română - (`ro`)
- Русский - (`ru`)
- 简体中文 - (`zh-cn`)
- Español (EUW) - (`es`)
- ภาษาไทย - (`th`)
- 繁體中文 - (`zh-hant`)
- Türkçe - (`tr`)
- Tiếng việt - (`vi`)

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS for news in english locale - https://data.rito.news/riotgames/en/news.rss
- Raw news data in korean - https://data.rito.news/riotgames/ko/news.json
