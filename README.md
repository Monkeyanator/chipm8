# chipm8
> Because I wanted to emulate a Gameboy but didn't know how

[![Go Report Card](https://goreportcard.com/badge/github.com/Monkeyanator/chipm8)](https://goreportcard.com/report/github.com/Monkeyanator/chipm8)
[![Go Report Card](https://travis-ci.com/Monkeyanator/chipm8.svg?branch=master)](https://travis-ci.com/Monkeyanator/chipm8)

![chipm8 Tetris](https://raw.githubusercontent.com/Monkeyanator/chipm8/master/images/chipm8.png)

Yet another chip8 emulator, written to learn more about emulator design before tackling a more complicated system. It comes with a stripped-down command-line debugger (think GDB with 3-4 commands).


## Installation

`chipm8` should compile down to any target, but SDL is required on the system the emulator is being built from (if statically linked, SDL not needed on system being run on). The instructions for configuring SDL [can be found here](https://github.com/veandco/go-sdl2).

OS X & Linux:

```sh
go get -u github.com/Monkeyanator/chipm8
```

In addition, this project uses [dep](https://github.com/golang/dep) for dependency management, so to fetch the needed dependencies, run:

```sh
dep ensure
```

## Usage example

To run a ROM, pass the path into `chipm8` through the `--prog` flag:
```sh
chipm8 --prog=roms/tetris.ch8
```
or, to run the Tetris ROM in debug mode:

```sh
chipm8 --prog=roms/tetris.ch8 --debug
```

Full list of command-line flags can be found with:

```sh
chipm8 --help
```

## Contributing

1. Fork it (<https://github.com/yourname/yourproject/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request
