## Twilio-Discord Integration Bot
Twilio-Discord Alert Bot is a powerful notification system designed to keep you informed about important events happening in your Discord channels and alert manager services. This bot listens to designated Discord channels for new messages and promptly sends an SMS or initiates a phone call to inform you about these updates.

Simultaneously, the bot can connect to remote or local alert manager services, continuously monitoring for any triggered alerts. When an alert is triggered, the bot instantly calls or sends an SMS to inform you about the situation.

Moreover, this bot doesn't just inform you about the problems, it also keeps you updated on their resolution. When an alert is resolved, the bot notifies you via call or SMS, ensuring you're always up-to-date about the status of your systems.

### Installation

To install and run the bot, follow these steps:

1. **Create a [Twilio trial account](https://www.twilio.com/docs/usage/tutorials/how-to-use-your-free-trial-account) to obtain the necessary account information for the environment variables.**
2. **Create a [Discord bot](https://discordpy.readthedocs.io/en/stable/discord.html) by following the instructions.**
3. **Create a Discord channel:**

   On your Discord server, right-click on the space where your channels are listed and click "Create Channel." 
   Give your channel a name and click "Create."
   Add the Bot to the Channel:

   Navigate to the settings of the bot you created previously.
   In the "Bot" section, you'll find a URL to add the bot to a server. Click on it.
   Select your server from the dropdown list and click "Continue."
   On the next page, ensure that the bot has the necessary permissions, then click "Authorize."

4. **Follow these steps to set up notifications from the specific external Discord channels that you are interested in.**
   
   Go to the External Discord Channel:

   Navigate to the external Discord channel you want to follow. Note that you need to have the 'Manage Channel' permission in your own channel, and the external channel must have 'Allow anyone to @mention this channel' enabled.
   In the external Discord channel, find a message that has been marked as an announcement (these messages have a small megaphone icon next to them). Click on the "Follow" button on this message.
   A dialog box will open asking where you want to send the announcements. Choose your own Discord channel from the dropdown menu and click "Follow."
5. **Install Go on your system.**
6. **Clone the repository using the command git clone git@github.com:bartosian/dstwilio.git.**
7. **Change to the cloned directory using the command ``cd dstwilio``.**
8. **Build the executable using the command ``go build -o dstwilio ./cmd/main.go``.**
9. **Move the executable to /usr/local/bin using the command ``sudo mv dstwilio /usr/local/bin``.**
10. Set the required environment variables and create a system service file. To create the service file, run the following commands:

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
ExecStart=/usr/local/bin/dstwilio --alerts --discord
Restart=on-failure

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

9. The --alerts flag is an optional command-line argument that allows the Discord call bot to connect to an Alertmanager running on the same machine. With this flag, the bot will listen for any alerts triggered by the Alertmanager and send a message to the specified Discord channel with information about the alert. When the alert recovers, the bot will send another message indicating that the alert has been resolved. This feature can be useful for receiving real-time alerts about system or application failures.

```shell
go run main.go --alerts
```

10. For a more customized experience, the bot supports an optional command-line flag --discord. When enabled, this flag allows the bot to listen to designated Discord channels for new messages.

```shell
go run main.go --discord
```

# License

Apache2.0