# vpp

A simple video and/or audio downloading, re-encoding, clipping program created for my own use cases. Many functions are dependent on ffmpeg to function as they are just a higher level wrapper.

> [!NOTE]
> This tool is not meant to download content you do not own or have permission to download. I take no responsibility for what you do with this piece of software. The URL's provided in the example commands are meant only as examples.

## Overview

- [vpp](#vpp)
  - [Overview](#overview)
  - [Quick Start](#quick-start)
  - [Examples](#examples)
    - [Downloading a video](#downloading-a-video)
      - [Using a URL](#using-a-url)
      - [Using an ID](#using-an-id)
      - [Specifying the quality](#specifying-the-quality)
    - [Downloading audio](#downloading-audio)
      - [mp3](#mp3)
      - [wav](#wav)
    - [Get info about a video](#get-info-about-a-video)
    - [Experimental features](#experimental-features)
      - [Re-encoding](#re-encoding)
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

To download a video, use the `-dl` flag and specify the URL or ID.

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

To specify the quality what you want to download the video as, use the `-q` flag followed by a number you got from the [`-info`](#get-info-about-a-video) flag. The default - 0, represents the highest available quality.

```shell
./vpp -dl -q 2 www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Downloading audio

#### mp3

To download only the audio part of the video as .mp3, you just need to add the `-mp3` flag to your command.

```shell
./vpp -dl -mp3 www.youtube.com/watch?v=dQw4w9WgXcQ
```

#### wav

To download only the audio part of the video as .wav, you just need to add the `-wav` flag to your command.

```shell
./vpp -dl -wav www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Get info about a video

To get info about a specific video, use the `-info` flag. The output will be displayed, ordered by a number. That number you can use to specify what quality you want to download the video as.

```shell
./vpp -info www.youtube.com/watch?v=dQw4w9WgXcQ
```

### Experimental features

> [!NOTE]
> These experimental features are subject to bugs, changes, and potential removal. Use with caution.

#### Re-encoding

Using the `-reencode` flag when downloading a video, will fully re-encode it using the x264/AAC codecs after the video and audio parts have been downloaded, instead of simply copying them. In most cases, this feature will make the video acquiring time significantly longer. This feature also fixes embed issues with Discord as it doesn't support h265 or AV1 encoding.

```shell
./vpp -dl -reencode www.youtube.com/watch?v=dQw4w9WgXcQ
```

## Contributing

Feel free to open issues or submit pull requests, but beware, they might not get accepted, as this is made for my own use cases.

## License

This project is licensed under the [MIT License](LICENSE).
