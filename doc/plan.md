This document is a running plan for the wildcat-countdown bot.

## Features

- ~~**The countdown bot is running on PWS**~~
- **The countdown bot prints something to stdout every X days**
- **The countdown bot prints something to a Slack channel every X days**
- **I can configure the countdown bot to print stuff every X days at runtime**
- **The countdown bot prints a countdown to a certain date every X days**
- **I can configure the countdown bot to print stuff via a printf format string**

## Bugs

None!

## Implementation

The countdown-wildcat bot will sit on a channel (via a time.Ticker) that fires every `<configured duration>`. When it fires, it will calculate the number of days to reunion, format a `<message>`, and tell a Slack client to send that message to a channel.

The `<configured duration>` is able to be updated via a PUT to /api/duration. When that is updated, the program will notify the bot that there is a new timeout. The bot will then recalculate the next time that it needs to run (using the current timeout's start date).

The `<message>` is able to be updated via a PUT to /api/message.
