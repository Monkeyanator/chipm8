# chipm8
> Because I wanted to emulate a Gameboy but didn't know how

[![Go Report Card](https://goreportcard.com/badge/github.com/Monkeyanator/chipm8)](https://goreportcard.com/report/github.com/Monkeyanator/chipm8)

Yet another chip8 emulator, written to learn more about emulator design before tackling a more complicated system. It comes with a stripped-down command-line debugger (think GDB with 3-4 commands).


## Installation

`chipm8` should compile down to any target, but SDL is required on the system the emulator is being built from (if statically linked, SDL not needed on system being run on). The instructions for configuring SDL [can be found here](https://github.com/veandco/go-sdl2).

OS X & Linux:

```sh
go get -u github.com/Monkeyanator/chipm8
```


## Usage example

```sh
chipm8 --prog=roms/tetris.ch8
```

## Contributing

1. Fork it (<https://github.com/yourname/yourproject/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request
