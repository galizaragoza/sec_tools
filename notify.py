#!/usr/bin/python3

# If this script is being executed, that means the last process ended successfully,
# since it should be runned like this:
# command1 && ./notify.sh (or add it to your path if you want it to be globally avaible)
# Now you can also run it like this: command1 && notify "Your custom notification"
# It literally just waits for the previous job to finish and then sends the string
# stored in to the chat_id of your choice. Of course, you need to
# create a bot and change the values in TG_token and chat_id (you can also change
# the message). For creating the bot, download Telegram desktop and use BotFather.

#!/usr/bin/python3
import os
import re
import requests
import sys
import urllib.parse


TG_token = <YOUR_TG_TOKEN>
bot_url = f"https://api.telegram.org/bot{TG_token}"
chat_id = <YOUR_CHAT_ID>
custom_msg = ""
default_msg = "The task you was waiting for is already finished ;)"
default_msg_encoded = urllib.parse.quote_plus(default_msg)
full_url = f"{bot_url}/sendMessage?chat_id={chat_id}&text="

if __name__ == "__main__":
    if len(sys.argv) == 1:
        msg = default_msg_encoded
        requests.post(full_url + msg)
        print(default_msg)
    elif len(sys.argv) == 2:
        custom_msg = sys.argv[1]
        custom_msg_encoded = urllib.parse.quote_plus(custom_msg)
        requests.post(full_url + custom_msg_encoded)
        print(custom_msg)
    else:
        print("Too many arguments")
