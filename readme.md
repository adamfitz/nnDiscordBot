# nnDiscord Bot

## Credential handling

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


## Config file

nnDiscordBot configuration file is located in the home directory of the user that is runing the application, it should 
be places in the following directory:

`~/.config/nnDiscordBot/`

The configuration file must be named: `nnDiscordCBot.config`

If you are intending to use any of the features that require configuration (sonarr integration etc) The contents of said 
file must be as follows:

```bash
$ cat ~/.config/nnDiscordBot/nnDiscordCBot.config 
{
	"sonarr_instance": "10.1.1.1",
	"sonarr_port": "8989"
}
```

Test master branch update...