# google-history-wallpaper

"Google History Wallpaper" - a project that uses search queries from your Google browser history to find and set images as your wallpaper. Automatically update your desktop background with relevance to your interests.


## Dependencies

- Chromium-based browser

## Install

```sh
make install
```

## Usage

```
ghw -wp-tool wbg

```

## Examples of usage

```
Usage of ghw:
  -browser string
        browser to use. Available: [vivaldi chrome chromium opera brave] (default "vivaldi")
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

## TODO (maybe never)

- [x] Add swaybd tool
- [x] Add wbg tool
- [x] Unsplash image provider
- [x] Chromium based browser support
- [ ] Firefox support
- [ ] Resolution auto detection
- [ ] Google image provider without API key
- [ ] Pinterest?