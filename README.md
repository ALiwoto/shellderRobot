
# Shellder Robot
A bot that uses telegram as a console shell, written in golang. 

<hr/>

## How to set up
Here is the sample config:

```ini
[general]
bot_token = 123456789:abcdefg
owner_ids = 1234587, 12345784, 1234578
download_directory = downloads
handler_command = vega
drop_updates = true
debug = false
```


- `handler_command` variable specifies the command for running scripts.
> For example if it's set to `vega`, the command will be `/vega shell command`.

- additional commands such as download and upload are implemented as well, using the `handler_command` as their prefix.
> For example `/vegaDl <specified_path>` (dl short for download and ul short for upload) will download the content of replied messages to the specified path. You can give an either relative or absolute path to download a file/media content. (uploading works similar as well with `/vegaUl <file_path>` command)

- `download_directory` is the directory in which all downloaded contents go, unless you specify a directory while using download command. Do notice that this directory should be put in `.gitignore` file as well.




