# Basic Setup

> replace the CAPITALIZED_SNAKE_CHARACTERS with your own environment

**Locate the directory (directory means folders in Linux).**

In the Linux version of the game, football manager is installed in a directory deep inside the `.local/share/Steam` directory. Yours should look something like this `$HOME/.local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager VERSION/`, with the `ID_NUMBER` different for every version of the game. For example, Football Manager 2024 should be located in `$HOME/.local/share/Steam/steamapps/compatdata/2252570/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager 2024/`

| Football Manager Version | ID_NUMBER |
| ------------------------ | --------- |
| 2024                     | 2252570   |
| 2023                     | 1904540   |
| 2022                     | 1569040   |
| 2021                     | 1263850   |
| 2020                     | 1100600   |

A way to find the directory/folder of your football manager save is to use [`fzf`](https://github.com/junegunn/fzf) or `grep`. You would see a directory like this. **You may not have the `graphics` directory.**

![steam-directory](/docs/img/steam-directory.png)

**Create the graphics directory if you don't have it.** In the game, all user modded graphics are placed under a `graphics/` directory. The graphics directory path might look something like this: `$HOME/.local/share/Steam/steamapps/compatdata/ID_NUMBER/pfx/drive_c/users/steamuser/Documents/Sports Interactive/Football Manager FM_VERSION/graphics`. The `graphics/` directory might not be existant and you may have to create your own.

![graphics-directory](/docs/img/graphics-directory.png)

**Download the New GAN face pack and place it in the graphics directory.** The latest facepack could be found at the original [NewGAN-Manager release page](https://github.com/Maradonna90/NewGAN-Manager/releases) and would be under `Facepack Download`. I recommend renaming it as newgen.

![facepack-download-img](/docs/img/facepack-download.png)

**Extract the file**, it should have subdirectories titled: `EECA`, `MESA` etc... Like this

![newgen-graphics-subdirectory](/docs/img/newgens-subdirectories.png)

**Copy the [`config.xml` file](/example/config.xml) (under `example/`) to the Facepack directory**

![config-file](/docs/img/config-xml.png)

**Generate the `.rtf` file** needed, name it `newgen.rtf` and place it in the Facepack directory. If you don't know how to do it, watch [Zealand's video on the original NEWGan Manager](https://www.youtube.com/watch?v=pmdIkhfmY6w), the steps are between `9:28` to `12:28`

![rtf-file](/docs/img/newgen-rtf.png)

**Download the binary and place it into the football manager `graphics` directory**

![download the binary](/docs/img/download-binary.png)

![place the binary in the right directory](/docs/img/jaqen-binary.png)

**Execute the binary by double clicking the jaqen binary** and check if your `config.xml` has been changed or not

![changed-xml](/docs/img/changed-xml.png)

**Reload the skin** and enjoy. Reloading is different on each version of the game but I'm sure there's a lot of guides on the internet and you could see Zealand doing it in his video/stream. Also it might be best to not enable caching while reloading the images.
