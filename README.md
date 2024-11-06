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
      - [Specifying the quality](#specifying-the-quality)
    - [Get info about a video](#get-info-about-a-video)
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

By default, all videos are downloaded at their highest quality.

#### Using a URL

```shell
./vpp -dl www.youtube.com/watch?v=dQw4w9WgXcQ
```

#### Using an ID

```shell
./vpp -dl dQw4w9WgXcQ
```

#### Specifying the quality

To specify the quality what you want to download the video as, use the `-q` flag followed by a number you got from the `-info` flag. The default - 0, represents the highest available quality.

```shell
./vpp -dl -q 2 www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Get info about a video

To get info about a specific video, use the `-info` flag. The output will be displayed, ordered by a number. That number you can use to specify what quality you want to download the video as.

```shell
./vpp -info www.youtube.com/watch?v=dQw4w9WgXcQ
```

## Contributing

Feel free to open issues or submit pull requests, but beware, they might not get accepted, as this is made for my own use cases.

## License

This project is licensed under the [MIT License](LICENSE).
