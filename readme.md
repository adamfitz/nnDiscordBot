# nnDiscord Bot

## Credential handling

Credentials are loaded from the home user directory file `~/.discordrc` this includes the api key/secret for opnsense FW
 and sonarr api token.

These creds need to be added to the above file in the below format:

```bash
{
    "bot_token": "xxx",
    "sonarr_api_token": "xxx",
	"opnsense_api_key": "xxx",
	"opnsense_api_secret":"xxx"
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

If you are intending to use any of the features that require configuration (sonarr api, remote DB Opnsense FW etc) The 
contents of said file must be as follows:

```bash
$ cat ~/.config/nnDiscordBot/nnDiscordCBot.config 
{
	"sonarr_instance": "10.23.0.3",
	"sonarr_port": "8989",
	"db_server": "db server name or IP",
	"db_port": "5432",
	"db_user": "db username",
	"db_user_pass": "db users password",
	"db_name": "your database name",
	"opnsense_wan_int":"your fw wan interface name",
	"opnsense_fw_ip":"your FW management IP",
}
```