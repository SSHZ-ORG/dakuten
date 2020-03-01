# 濁゙点゙ ﾞBﾞoﾞtﾞ

A Telegram Inline bot that adds 濁点 to everything.

Deployed to GCP Cloud Functions.

## Usage

Go to [@dakutenbot](https://t.me/dakutenbot) or [@kitanaibot](https://t.me/kitanaibot) and send it some message.
It will convert for you.

Or you can also use inline mode: `@dakutenbot <Whatever You Want>` or `@kitanaibot <Whatever You Want>`.

## Caveats

We process at rune (Unicode code point) level, not grapheme cluster level.
So you will see weird results if any combining character is present.
(But we do handle U+3099 and U+309A specially.)
