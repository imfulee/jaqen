# Jaqen

Create and manage your image file mapping to face profiles in Football Manager. Inspired by [Football Manager New GAN Manager written by Marco et al.](https://github.com/Maradonna90/NewGAN-Manager) I named it Jaqen based on Jaqen H'ghar having a wall of faces.

## Motivation

I found the original didn't run well on Linux, but works pretty well for Windows and Mac. **You could just hook up a Virtual Machine and a volume and use the original inside the Virtual Machine** but I decided to write this fun side project. I chose Go because I could compile to multiple platforms with a relatively easy learning curve and didn't have a complicated packaging step. I don't intend on supporting it at all, you're very welcomed to submit PRs and Issues but I don't guarantee to fix/review it.

## Usage

> If you're not interested in configuring your own setup, just read the basic setup

Download the repository and `go install`

There are flags that you could use to specify the paths for various files if you would wish to change the defaults

- `-xml` specifies the xml path
- `-rtf` specifies the rtf path
- `-img` specifies the image root directory
- `-p` preserves the current xml mapping
- `-ver` could specify the football manager version, this defaults to `2024`, and all other string will be parsed as "others"

## Basic Setup

> replace the CAPITALIZED_SNAKE_CHARACTERS with your own environment

**Locate the directory.** In the linux version of the game, football manager is installed in a directory deep inside the `.local/share/Steam` directory. If you are running football manager 2024, yours should look like this `$HOME/.local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager 2024/`, with the `ID_NUMBER` different for everyone. You would see a directory like this. **You may not have the `graphics` directory.**

![steam-directory](/docs/img/steam-directory.png)

**Create the graphics directory if you don't have it.** In the game, all user modded graphics are placed under a `graphics/` directory. The graphics directory path might look something like this: `$HOME/.local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager FM_VERSION/graphics`. The `graphics/` directory might not be existant and you may have to create your own.

![graphics-directory](/docs/img/graphics-directory.png)

**Download the New GAN face pack and place it in the graphics directory.** The latest facepack could be found at the original [NewGAN-Manager release page](https://github.com/Maradonna90/NewGAN-Manager/releases) and would be under `Facepack Download`. I recommend renaming it as newgen.

![facepack-download-img](docs/img/facepack-download.png)

**Extract the file**, it should have subdirectories titled: `EECA`, `MESA` etc... Like this

![newgen-graphics-subdirectory](/docs/img/newgens-subdirectories.png)

**Copy the [`config.xml` file](/example/config.xml) (under `example/`) to the Facepack directory**

![config-file](/docs/img/config-xml.png)

**Generate the `.rtf` file** needed, name it `newgen.rtf` and place it in the Facepack directory. If you don't know how to do it, watch [Zealand's video on the original NEWGan Manager](https://www.youtube.com/watch?v=pmdIkhfmY6w), the steps are between `9:28` to `12:28`

![rtf-file](/docs/img/newgen-rtf.png)

**Download the repository, build the binary and move the binary to your graphics directory**

```bash
cd REPOSITORY_PATH
go build .

fm_graphics_dir="$HOME/.local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager FM_VERSION/graphics/newgen"
install jaqen $fm_graphics_dir
```

**Execute the binary** and check if your `config.xml` has been changed or not

```bash
fm_graphics_dir="$HOME/.local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager FM_VERSION/graphics/newgen"
cd $fm_graphics_dir
./jaqen
```

![changed-xml](/docs/img/changed-xml.png)

**Reload the skin** and enjoy. Reloading is different on each version of the game but I'm sure there's a lot of guides on the internet and you could see Zealand doing it in his video/streams

If you don't want to run the long command in the future, you could also setup an alias or function for the command in your `.bashrc` or `.zshrc` or whatever you use.

```bash
newgancli() {
    fm_graphics_dir=".local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager FM_VERSION/graphics/"
    $fm_graphics_dir/jaqen
}
```

## Future Wants

This is just some notes on what I want it to do in the future.

- Build a GUI, maybe with Go Wails?
- Create a build pipeline to create executables
- There are some performance left on the table, currently the way reading and writing to file works relatively slow compared to what a buffered read and a generator could do. Probably faster than Python though :p
- Remove the need to copy the `config.xml` file into the directory
