# nnDiscord Bot

Credentials are loaded from the home user directory file `~/.discordrc`

These creds need to be added to the above file in the below format:

```bash
{
    "bot_token": "xxx"
    "sonarr_api_token": "xxx
}
```

The bot token is the most important and sensitivate as this is how the bot interacts with the discord servers.  The 
credentials (client id/secret) are only for adding a bot to a channel etc.  Do **not** let someone else get hold of 
your bot token, as then they will have access to your bot.

