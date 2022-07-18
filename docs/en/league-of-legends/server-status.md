# League of Legends Server status

The service generates files that contains theserver status message from League of Legends server status page (https://status.riotgames.com/lol).

## File URL
The files url looks like
```
https://rito-news.iamantosik.me/lol/{locale}/status.{server}.{extension}
```

### Available servers and locales
- Brazil (`br`)
    - Português (`pt-br`)
- EU Nordic & East (`eune`)
    - English (`en-gb`)
    - Čeština (`cs-cz`)
    - Ελληνικά (`el-gr`)
    - Magyar (`hu-hu`)
    - Polski (`pl-pl`)
    - Română (`ro-ro`)
- EU West (`euw`)
    - English (`en-gb`)
    - Deutsch (`de-de`)
    - Español (`es-es`)
    - Français (`fr-fr`)
    - Italiano (`it-it`)
- Japan (`jp`)
    - 日本語 (`ja-jp`)
- Korea (`kr`)
    - 한국어 (`ko-kr`)
- Latin America North (`lan`)
    - Español (`es-mx`)
- Latin America South (`las`)
    - Español (`es-ar`)
- North America (`na`)
    - English (`en-us`)
- Oceania (`oce`)
    - English (`en-au`)
- Russia (`ru`)
    - Русский (`ru-ru`)
- Turkey (`tr`)
    - Türkçe (`tr-tr`)
- Public Beta Environment (`pbe`)
    - All of above

### Available extensions
- RSS (`.rss`)
- Atom (`.atom`)
- JSONFeed (`.jsonfeed`)
- Raw data (`.json`)

### Examples
- RSS for NA server status messages - https://rito-news.iamantosik.me/lol/en-us/status.na.rss
- Raw server status data for Korean server - https://rito-news.iamantosik.me/lol/ko-kr/status.kr.json
