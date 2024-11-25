# TickrFetch

TickrFetch is a Go-based application that fetches cryptocurrency data from the CoinLore API and provides a Telegram bot interface to query coin prices and other details.

## Features

- Fetch cryptocurrency data from the CoinLore API.
- Query coin prices, market cap, and percentage changes via a Telegram bot.
- Format responses with aligned text and escape special characters for Markdown.

## Installation

### Prerequisites

- Docker
- Telegram Bot Token (create a bot using [BotFather](https://core.telegram.org/bots#botfather))

### Environment Variables

Create a `.env` file in the project root with the following content:

```env
BOT_TOKEN=your-telegram-bot-token
```

Replace `your-telegram-bot-token` with your actual Telegram bot token.

## Building and Running the Docker Container

### Build the Docker Image

Navigate to the project directory and build the Docker image:

```sh
docker build -t tickrfetch .
```

### Run the Docker Container in Detached Mode

Run the Docker container in detached mode:

```sh
docker run -d --name tickrfetch-container --env-file .env tickrfetch
```

### Verify the Container is Running

You can verify that the container is running by using the following command:

```sh
docker ps
```

## Usage

### Interacting with the Bot

1. Open Telegram and search for your bot.
2. Start a chat with your bot.
3. Use the following commands to interact with the bot:
   - `$<symbol>`
   - `/<symbol>`

### Example Commands

- `$ETH`
- `/XMR`

### Example Response

```
#1 Bitcoin (BTC)

Price        $50,000
M-Cap        $1,000,000,000

1h Change    0.5%
24h Change   2.3%
7d Change    5.6%
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
