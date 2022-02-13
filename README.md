
# Shellder Robot
===============

A bot that uses telegram as a console shell, written in go lang. 



==============
In the config, there is a variable for setting up the trigger

[general]
bot_token = 123456789:abcdefg
owner_ids = 1234587, 12345784, 1234578
handler_command = vega
drop_updates = true
debug = false

It's handler_command variable
consider it's set to vega
so by using /vega cmd it runs the command in shell and returns the result
/vegaDownload, /vegaDownload <path> and  /vegaDownload <file_name> will download the file (first one downloads the file to current bot directory and saves the file with its file id as its name)
as for file path, you can give it a relative path, or an absolute path

it have /vegaUpload <path> as well
if you give it a file name, it will search for file in the current directory of the bot
similar to download, you can give absolute or relative path as well




