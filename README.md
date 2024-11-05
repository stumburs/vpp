# vpp

A simple video and/or audio downloading, re-encoding, clipping program created for my own use cases. Many functions are dependent on ffmpeg to function as they are just a higher level wrapper.

> [!NOTE]
> At the moment, it lack most aforementioned features.

## Overview

- [vpp](#vpp)
  - [Overview](#overview)
  - [Quick Start](#quick-start)
  - [Examples](#examples)
    - [Downloading a video](#downloading-a-video)
      - [Using a URL](#using-a-url)
      - [Using an ID](#using-an-id)
  - [Contributing](#contributing)
  - [License](#license)

## Quick Start

1.  ```shell
    git clone https://github.com/stumburs/vpp.git
    ```

2.  ```shell
    cd vpp
    ```

3.  ```shell
    go build cmd/vpp/vpp.go
    ```

## Examples

### Downloading a video

#### Using a URL

```shell
./vpp -dl www.youtube.com/watch?v=dQw4w9WgXcQ
```

#### Using an ID

```shell
./vpp -dl dQw4w9WgXcQ
```

## Contributing

Feel free to open issues or submit pull requests, but beware, they might not get accepted, as this is made for my own use cases.

## License

This project is licensed under the [MIT License](LICENSE).
