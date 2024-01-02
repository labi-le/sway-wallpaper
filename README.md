# sway-wallpaper

Sway Wallpaper - a project that dynamically manages the wallpaper on the Sway desktop. It can use browser history to search for images, and it can also use a search phrase to search for images on the web. Several wallpaper managers are supported, as well as several image search APIs
## Working mode
- Browser history search mode
- Search by phrase mode (`search-phrase` option)

## Dependencies

- Wayland
- Chromium\Firefox-based browser - parsing history file, non required if use `search-phrase` option
- Utilities for setting wallpaper -`swaybg`
- `xrandr` - for resolution detection, non required if you specify resolution manually
- Internet connection - for find and downloading images
- Go - for building, non required if you use prebuilt binaries

## Compile from source

```sh
make install
```

## Prebuilt binaries
https://github.com/labi-le/sway-wallpaper/releases

## Usage

```
sway-wallpaper
```

## Examples of usage

```
Usage of sway-wallpaper:
  -browser string
        browser to use. Available: [vivaldi chrome chromium brave opera firefox] (default "vivaldi")
  -follow string
        follow a time interval and update wallpaper. e.g. 1h, 1m, 30s
  -history-file string
        browser history file to use. Auto detect if empty (only for chromium based browsers)
        e.g ~/.mozilla/icecat/gxda4hpz.default-1672760493248/formhistory.sqlite
  -resolution string
        resolution to use. e.g. 1920x1080 (default "1920x1080")
  -save-image-dir string
        directory to save image to (default "/home/labile/Pictures")
  -search-phrase string
        search phrase to use
  -wp-api string
        wallpaper api to use. Available: [unsplash] (default "unsplash")
  -wp-tool string
        wallpaper tool to use. Available: [swaybg] (default "swaybg")
```

Update your wallpaper every hour with wbg manager:

```
sway-wallpaper -wp-tool wbg -follow 1h
```

Update your wallpaper every 30 minute with swaybg manager:

```
sway-wallpaper -follow 30m
```

Use sway-wallpaper as dynamic wallpaper:

```
sway-wallpaper -follow 1h -search-phrase space
```

Direct path to browser history file (firefox-based)
```
sway-wallpaper -browser firefox -history-file ~/.mozilla/icecat/gxda4hpz.default-1672760493248/formhistory.sqlite
```

Direct path to browser history file (chromium-based)
```
sway-wallpaper -browser vivaldi -history-file ~/.config/vivaldi/Default/History
```

Add sway autostart:

```
exec sway-wallpaper [options]
```

## TODO (maybe never)

- [x] Add swaybg tool
- [x] Add logger
- [x] Add wbg tool
- [x] Unsplash image provider
- [x] Chromium based browser support
- [x] Firefox support
- [x] Resolution auto detection
- [ ] Google image provider without API key
- [ ] Pinterest?