## Discord Call Bot
This repository contains the code for a Discord bot that can initiate a phone call using Twilio when a message is posted to a specific channel. This bot is written in Go.

### Installation

To install and run the bot, follow these steps:

1. Create a [Twilio trial account](https://www.twilio.com/docs/usage/tutorials/how-to-use-your-free-trial-account) to obtain the necessary account information for the environment variables.
2. Create a [Discord bot](https://discordpy.readthedocs.io/en/stable/discord.html) by following the instructions.
3. Install Go on your system.
4. Clone the repository using the command git clone git@github.com:bartosian/dstwilio.git.
5. Change to the cloned directory using the command cd dstwilio.
6. Build the executable using the command go build -o dstwilio main.go.
7. Move the executable to /usr/local/bin using the command sudo mv dstwilio /usr/local/bin.
8. Set the required environment variables and create a system service file. To create the service file, run the following commands:

```shell
echo "[Unit]
Description=Discord Call Bot
After=network.target

[Service]
Type=simple
User=root
Environment=\"DISCORD_BOT_TOKEN=<DISCORD_BOT_TOKEN>\"
Environment=\"TWILIO_ACCOUNT_SID=<TWILIO_ACCOUNT_SID>\"
Environment=\"TWILIO_AUTH_TOKEN=<TWILIO_AUTH_TOKEN>\"
Environment=\"TWILIO_PHONE_NUMBER=<TWILIO_PHONE_NUMBER>\"
Environment=\"YOUR_PHONE_NUMBER=<YOUR_PHONE_NUMBER>\"
Environment=\"DISCORD_CHANNEL=<DISCORD_CHANNEL>\"
ExecStart=/usr/local/bin/dstwilio
Restart=on-failure
RestartSec=always

[Install]
WantedBy=multi-user.target" > $HOME/dstwilio.service

mv $HOME/dstwilio.service /etc/systemd/system/

sudo tee <<EOF >/dev/null /etc/systemd/journald.conf
Storage=persistent
EOF
```
7. Reload the systemd daemon and enable and start the service using the following commands:

```shell
sudo systemctl daemon-reload
sudo systemctl enable dstwilio
sudo systemctl start dstwilio
sudo systemctl status dstwilio
```

8. To run the bot manually and read environment variables from a .envrc file, run the following command:

```shell
go run main.go --envrc
```

# License

Apache2.0