# Photo Memories Telegram Bot

**Photo Memories** is a Telegram bot designed for managing photo memories. The bot allows users to upload their favorite photos, as well as delete them. Additionally, the bot provides the ability to retrieve photos by tags, such as "vacation", "holiday", etc. Furthermore, a worker sends one photo per day to remind the user of pleasant memories.

## Key Features

- **Adding photos**: Users can upload their favorite photos to the bot, which saves them for later use.
- **Deleting photos**: Users can delete photos they no longer wish to keep.
- **Retrieving photos by tags**: The bot allows users to retrieve photos associated with specific tags, such as "travel", "friends", "family", etc.
- **Daily photos**: A worker selects one photo each day and sends it to the user to remind them of nice moments.

## Technologies

- **Go** — the main programming language for implementing the bot's logic and interacting with the API.
- **Docker** — used for packaging and deploying the application.
- **PostgreSQL** — the database for storing information about photos and tags.
- **Redis** — data caching to improve performance.
- **AWS S3** — a storage solution for uploading and storing photos.
- **Telegram Bot API** — for interacting with users through Telegram.
- **Logrus** — for logging and tracking events within the application.
- **Viper** — for managing configurations and environment variables.
- **CI/CD** — continuous integration and deployment system for automated build, testing, and deployment to a remote server on AWS.

## How to Run Locally

1. Clone the repository:

   ```bash
   git clone https://github.com/your-repo/photo-memories-bot.git
   cd photo-memories-bot
