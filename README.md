<h1 align="center">Jaqen</h1>

Create and manage your image file mapping to face profiles in Football Manager. Inspired by [NewGAN-Manager](https://github.com/Maradonna90/NewGAN-Manager), I named it Jaqen based on Jaqen H'ghar having a wall of faces.

## Motivation

I found the original didn't run well on Linux, but works pretty well for Windows and Mac. **You could just hook up a Virtual Machine and a volume and use the original inside the Virtual Machine**, like [this](https://youtu.be/MV-iB1v_UdU) (credits to the reddit `u/EvensenFM` that did it), but I decided to write this fun side project. I chose Go because I could compile to multiple platforms with a relatively easy learning curve and didn't have a complicated packaging step. I don't play this game as much anymore and you're very welcomed to submit PRs and issues.

## Usage

**If you're not interested in configuring your own setup, just read the [basic setup](./docs/basic_setup.md)**

These are the flags that you could use to specify the paths for various files if you would wish to change the defaults

- `--xml` specifies the xml path. Defaults to `./config.xml`
- `--rtf` specifies the rtf path. Defaults to `./newgan.rtf`
- `--img` specifies the image root directory. Defaults to `./`
- `--preserve` preserves the current xml mapping. Defaults to not preserve.
- `--version` could specify the football manager version. Defaults to `2024`, all other values will be ignored.
- `-config` specifies the config directory. Defaults to `./jaqen.toml`

All paths are relative to the binary.

```bash 
jaqen \
    --xml=/path/to/config.xml \
    --rtf=/path/to/newgan.rtf \ 
    --img=/path/to/images/directory \
    --preserve \ 
    --version=2024 \ 
    --config=/path/to/config
```

To format the config toml file

```bash
jaqen format /path/to/jaqen.toml
```

### Config file options 

It's basically the command line flags but in a file. You could see an example [here](./example/jaqen.toml). Flags will take precendents over config file options, which itself will take precendents over the defaults. The only difference is the `[mapping_override]` section, it will look something like this: 

```toml
[mapping_override]
AFG = 'MESA'
```

The first word "AFG" represents the country initials "Afghanistan", while the second word describes the type of faces that it should use, which in this example is "MESA" (it is a code for middle east south asian). The point of this is to be able to override the default mappings from nations to a certain facepack. For example, if I want to make Afghans look South East Asian, you could change it to 

```toml
[mapping_override]
AFG = 'Seasian'
```

Or if you loaded a database that has a country not listed in the game, you could also add a country face mapping in there

```toml
[mapping_override]
XYZ = 'EECA'
```

These are the current code for faces 

| Ethnic group                |Code for the faces|
|-----------------------------|------------------|
| African                     | African          |
| Asian                       | Asian            |
| Caucasian                   | Caucasian        |
| CentralEuropean             | Central European |
| EasternEuropeanCentralAsian | EECA             |
| ItalianMediterranean        | Italmed          |
| MiddleEastNorthAfrican      | MENA             |
| MiddleEastSouthAsian        | MESA             |
| SouthAmericanMediterranean  | SAMed            |
| Scandinavian                | Scandinavian     |
| SouthEastAsian              | Seasian          |
| SouthAmerican               | South American   |
| SpanishMediterranean        | SpanMed          |
| YugoslavGreek               | YugoGreek        |

## Future Wants

This is just some notes on what I want it to do in the future.

- Build a GUI, maybe with Go Wails? But if you already play a spreadsheet game I assume you would also know how to read some instructions :p
- There are some performance left on the table, currently the way reading and writing to file works relatively slow compared to what a buffered read and a generator could do. Probably faster than Python though :p
- Remove the need to copy the `config.xml` file into the directory
