# confusion

Simple configuration management for Go projects. Inspired by [cristalhq/aconfig](https://github.com/cristalhq/aconfig).

## Install

`$ go get github.com/njkleiner/confusion`

## Usage

```go
package example

import (
    "fmt"

    "github.com/njkleiner/confusion"
    "github.com/njkleiner/confusion/toml"
)

type ExampleConfig struct {
    Foo, Bar string
}

func Example() {
    opts := confusion.Options{
        Prefix: "example",
        UserPaths: []string{
            "$XDG_CONFIG_HOME",
            "$HOME/.config",
        },
        SystemPaths: []string{
            "/etc",
        },
        Loaders: map[string]confusion.Loader{
            ".toml": toml.Loader,
        },
    }

    config := ExampleConfig{}

    // Loads the config file located at "$HOME/.config/example/config.toml"
    err := confusion.LoadConfig("config.toml", opts, &config)

    if err != nil {
        panic(err)
    }

    fmt.Printf("loaded: %#v", config)
}
```

## Contributing

You can contribute to this project by [sending patches](https://git-send-email.io) to `noah@njkleiner.com`. Pull Requests are also welcome.

## Authors

* [Noah Kleiner](https://github.com/njkleiner)

See also the list of [contributors](https://github.com/njkleiner/confusion/contributors) who participated in this project.

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.
