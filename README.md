# history-wallpaper

"History Wallpaper" - a project that uses search queries from your browser history to find and set images
as your wallpaper. Automatically update your desktop background with relevance to your interests.

## Dependencies

- Wayland
- Chromium-based browser
- Utils swaybg\wbg
- Internet connection
- Go (for building)

## Install

```sh
make install
```

## Usage

```
hw
```

## Examples of usage

```
Usage of hw:
  -browser string
        browser to use. Available: [vivaldi chrome chromium opera brave] (default "vivaldi")
  -follow string
        follow a time interval and update wallpaper. e.g. 1h, 1m, 30s
  -resolution string
        resolution to use. e.g. 1920x1080 (default "1920x1080")
  -save-image-dir string
        directory to save image to (default "/home/labile/Pictures")
  -search-phrase string
        search phrase to use
  -wp-api string
        wallpaper api to use. Available: [unsplash] (default "unsplash")
  -wp-tool string
        wallpaper tool to use. Available: [swaybg wbg] (default "swaybg")
```

Update your wallpaper every hour with wbg manager:

```
hw -wp-tool wbg -follow 1h
```

Update your wallpaper every 30 seconds with swaybg manager:

```
hw -follow 30m
```

Use ghw as dynamic wallpaper:

```
hw -follow 1h -search-phrase space
```

Add sway autostart:

```
exec hw [options]
```

## TODO (maybe never)

- [x] Add swaybg tool
- [x] Add wbg tool
- [x] Unsplash image provider
- [x] Chromium based browser support
- [ ] Firefox support
- [ ] Resolution auto detection
- [ ] Google image provider without API key
- [ ] Pinterest?