<h1 align="center">Jaqen</h1>

Create and manage your image file mapping to face profiles in Football Manager. Inspired by [NewGAN-Manager](https://github.com/Maradonna90/NewGAN-Manager), I named it Jaqen based on Jaqen H'ghar having a wall of faces.

## Motivation

I found the original didn't run well on Linux, but works pretty well for Windows and Mac. **You could just hook up a Virtual Machine and a volume and use the original inside the Virtual Machine** but I decided to write this fun side project. I chose Go because I could compile to multiple platforms with a relatively easy learning curve and didn't have a complicated packaging step. I don't play this game as much anymore and you're very welcomed to submit PRs and issues but I don't guarantee to fix/review it.

## Usage

**If you're not interested in configuring your own setup, just read the [basic setup](https://github.com/jaqen/tree/release/docs/basic_setup.md)**

These are the flags that you could use to specify the paths for various files if you would wish to change the defaults

- `-xml` specifies the xml path
- `-rtf` specifies the rtf path
- `-img` specifies the image root directory
- `-p` preserves the current xml mapping
- `-ver` could specify the football manager version, this defaults to `2024`, and all other string will be parsed as "others"

All paths are relative to the binary.

## Future Wants

This is just some notes on what I want it to do in the future.

- Build a GUI, maybe with Go Wails?
- There are some performance left on the table, currently the way reading and writing to file works relatively slow compared to what a buffered read and a generator could do. Probably faster than Python though :p
- Remove the need to copy the `config.xml` file into the directory
