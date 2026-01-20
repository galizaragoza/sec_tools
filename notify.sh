#!/bin/bash

# If this script is being executed, that means the last process ended successfully,
# since it should be runned like this:
# command1 && ./notify.sh
# It literally just waits for the previous job to finish and then sends the string
# stored in msg_URLencoded to the chat_id of your choice. Of course, you need to
# create a bot and change the values in TG_token and chat_id (you can also change
# the message). For creating the bot, download Telegram desktop and use BotFather.

TG_token="<your_telegram_token>"
bot_url="https://api.telegram.org/bot$TG_token"
msg_URLencoded="La%20tarea%20por%20la%20que%20estabas%20esperando%20se%20ha%20completado%20con%20%C3%A9xito%F0%9F%98%8E"
chat_id="<your_chat_or_group_id>"

echo "La tarea anterior se ha completado con éxito"
curl -X POST "$bot_url/sendMessage?chat_id=$chat_id&text=$msg_URLencoded"
