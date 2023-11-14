# Jaqen

> This is very much work in progress.

A Football Manager New GAN Manager rewrite in go. The original could be found [here](https://github.com/Maradonna90/NewGAN-Manager).

I named it Jaqen based on Jaqen H'ghar having a wall of faces.

## Motivation

I just wanted to learn go and just found the original didn't run well on Linux, as I think the original version works pretty well for Windows and Mac. However, this is just a fun side project for me. Even if I finish this off, I don't intend on supporting it long term, you're very welcomed to submit PRs and Issues but I don't guarantee to fix/review it.

Could I have just found a way to fix the Linux CI on that project? probably.

## Usage

> Replace all the CAPITALIZE_SNAKE_CASE with your own names, ex. `RTF_FILE_NAME` could be whatever you named the rtf file like: newgen.rtf

In this game, all user modded graphics are placed under a `graphics/` folder. The graphics folder path might look something like this: `$HOME/.local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager FM_VERSION/graphics`

1. Download the New GAN face pack and place it in the graphics folder (TODO: exact url for this)
2. (Optional) Copy the `config.xml` file to the graphics folder
3. Place the `views/PlayerSearch.fmf` as a views folder (TODO: exact steps for this)
4. Generate the `.rtf` file needed and place it in the graphics folder
5. Place the binary inside the graphics folder
   - By now under the root of the `graphics/` folder, you should have at least the `jaqen` binary, `IMAGE_FOLDER_ROOT_NAME` image folder root and if you have a previous config xml file, you could just add it to the
6. Type the command `./jaqen RTF_FILE_NAME IMAGE_FOLDER_ROOT_NAME --preserve --xml=XML_FILE_NAME` with `--preserve`, `--xml` as options
7. (TODO: steps about reloading the skin)
