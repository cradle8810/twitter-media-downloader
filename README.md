# twmd: CLI twitter media downloader (without api key)

This twitter downloader doesn't require Credentials or an api key. It's based on [twitter-scrapper](https://github.com/imperatrona/twitter-scraper).

# Usage:
You can get the images and videos address with this command. You need volume mount for twmd_cookies.json.

```
docker run -t --rm -v "$(pwd)":/work ghcr.io/cradle8810/twmd -t 2007935180863103398
```

# Thanks

- https://github.com/mmpx12/twitter-media-downloader
  - An original program
- https://github.com/jeffrey12cali/twitter-scraper
  - Modified X scraper program (late 2025)
