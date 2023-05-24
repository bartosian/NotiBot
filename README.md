## NotiBot
NotiBot is a powerful notification system designed to keep you informed about important events happening in your Discord channels and alert manager services. This bot listens to designated Discord channels for new messages and promptly sends an SMS or initiates a phone call to inform you about these updates.

Simultaneously, the bot can connect to remote or local alert manager services, continuously monitoring for any triggered alerts. When an alert is triggered, the bot instantly calls or sends an SMS to inform you about the situation.

Moreover, this bot doesn't just inform you about the problems, it also keeps you updated on their resolution. When an alert is resolved, the bot notifies you via call or SMS, ensuring you're always up-to-date about the status of your systems.

### Installation

To enable the bot to call you and send SMS on new messages received from Discord channels you are subscribed to, follow these steps. If you do not wish to be notified about new messages in Discord channels, you can skip this section:

1. **Create a [Discord bot](https://discordpy.readthedocs.io/en/stable/discord.html) by following the instructions.**
2. **Create a Discord channel:**

   - On your Discord server, right-click on the space where your channels are listed and click "Create Channel."
   - Give your channel a name and click "Create."

3. Add the Bot to the Channel:

   - Navigate to the settings of the bot you created previously.
   - In the "Bot" section, you'll find a URL to add the bot to a server. Click on it.
   - Select your server from the dropdown list and click "Continue."
   - On the next page, ensure that the bot has the necessary permissions, then click "Authorize."

4. **Follow these steps to set up notifications from the specific external Discord channels that you are interested in.**

   - Navigate to the external Discord channel you want to follow. Note that you need to have the 'Manage Channel' permission in your own channel, and the external channel must have 'Allow anyone to @mention this channel' enabled.
   - In the external Discord channel, find a message that has been marked as an announcement (these messages have a small megaphone icon next to them). Click on the "Follow" button on this message.
   - A dialog box will open asking where you want to send the announcements. Choose your own Discord channel from the dropdown menu and click "Follow."

At the moment, our notification system only supports Twilio as the external provider. However, we have plans to expand our support in the near future to include other providers such as PagerDuty. Stay tuned for updates on the availability of additional notification providers.


To create a Twilio account and acquire the phone number(s) with the necessary permissions, please follow the instructions below:

1. **[Create a Twilio account](https://www.twilio.com/try-twilio).**
2. **[Purchasing a Phone Number](https://www.twilio.com/console/phone-numbers/search).**
3. **[Finding Account SID and Auth Token](https://www.twilio.com/console).**

Once you have completed the configuration of the external dependencies, proceed with installing the required packages and setting up a systemd service to run the bot:

1. **Install Go on your system.**
2. **Clone the repository using the command:**
```shell
git clone git@github.com:bartosian/notibot.git
```
3. **Change to the cloned directory using the command:** 
```shell
cd notibot
```
4. **Build the executable using the command:**
```shell
go build -o notibot ./cmd/main.go
```
5. **Move the executable to ``/usr/local/bin`` using the command:** 
```shell 
sudo mv notibot /usr/local/bin
```
6. **Set the required environment variables and create a system service file. To create the service file, run the following commands:**

```shell
echo "[Unit]
Description=Twilio-Discord Integration Bot
After=network.target

[Service]
Type=simple
User=root
Environment=\"DISCORD_BOT_TOKEN=<DISCORD_BOT_TOKEN>\"
Environment=\"TWILIO_ACCOUNT_SID=<TWILIO_ACCOUNT_SID>\"
Environment=\"TWILIO_AUTH_TOKEN=<TWILIO_AUTH_TOKEN>\"
Environment=\"TWILIO_PHONE_NUMBER=<TWILIO_PHONE_NUMBER>\"
Environment=\"CLIENT_PHONE_NUMBER=<YOUR_PHONE_NUMBER>\"
Environment=\"DISCORD_CHANNEL=<DISCORD_CHANNEL>\"
Environment=\"ALERT_MANAGER_URL=<ALERT_MANAGER_URL>\"
ExecStart=/usr/local/bin/notibot --alerts --discord
Restart=on-failure

[Install]
WantedBy=multi-user.target" > $HOME/notibot.service

mv $HOME/notibot.service /etc/systemd/system/

sudo tee <<EOF >/dev/null /etc/systemd/journald.conf
Storage=persistent
EOF
```
7. **Reload the systemd daemon, enable and start the service using the following commands:**

```shell
sudo systemctl daemon-reload
sudo systemctl enable notibot
sudo systemctl start notibot
sudo systemctl status notibot
```

If you choose not to use a systemd service to run the bot or need to provide environment variables using a .envrc file for development purposes, utilize the following flag when executing the command:

```shell
go run main.go --envrc
```

The optional flag ``--alerts`` allows the bot to listen to Alertmanager (via the ``ALERT_MANAGER_URL`` environment variable). When active alerts are encountered, the bot will initiate a phone call and send an SMS to the provided phone number. If you do not require this type of alerting, you can omit the ``--alerts`` flag. Additionally, the bot will notify you when all alerts have been resolved.

```shell
go run main.go --alerts
```

When the optional flag ``--discord`` is enabled, the bot will actively monitor specified Discord channels for new messages. Upon receiving any messages, the bot will initiate a phone call and send an SMS to the configured Twilio number.

```shell
go run main.go --discord
```

# License

Apache2.0