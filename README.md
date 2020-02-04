# blackdesert-monitor

A Telegram bot that notifies you when your Black Desert stops running or loses connection to the server.

## Setup

After downloading the executable (or compiling it yourself):

1. Create a new Telegram Bot [(Read this)](https://core.telegram.org/bots#6-botfather)
2. Copy your Bot **token** and insert it into your **settings.yaml** file.
3. Start the program
4. Send **/myid** to your Bot, it should respond with your **User ID**
5. Copy the **User ID** into your **settings.yaml** file
6. Restart the Bot
7. Start Black Desert Online
8. ...
9. Profit!

## Frequently Asked Questions

### How does this work?

The Bot doesn't interact with the Black Desert client directly (besides killing it when you picked that option),
it checks which Process ID is attached to the executable and then checks if that process has
any network connections via [win-netstat](https://github.com/pytimer/win-netstat).

### Why will the bot not close my game?

The bot needs admin privileges in order to kill the Black Desert process.

## License

MIT