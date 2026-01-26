# chiasma

dynamically manages wallpapers based on browser history or search phrases

## features

- **source**:
  - **browser history**: extracts last search query (chromium-based/firefox).
  - **manual phrase**: static keyword search.
- **api**:
  - **unsplash**
  - **nasa**
- **backends**:
  - `swww`
  - `swaybg`
- **output**:
  - auto-detection via `xrandr`.
  - multi-monitor support.
- **modes**:
  - one-shot.
  - daemon (`--follow`).

## dependencies

- **wayland**
- **backend** (one of):
  - `swww`
  - `swaybg`
- **resolution**:
  - `xrandr` (optional, for auto-detection).
- **browser** (optional):
  - chromium-based (chrome, brave, vivaldi, opera, etc.)
  - firefox

## installation

- [Prebuilt binaries](https://github.com/labi-le/chiasma/releases)

## usage

```bash
chiasma [flags]
```

### flags

```shell
      --api string              image source api (default "nasa")
      --browser string          browser name (default "google-chrome")
      --follow                  enable periodic updates
      --history-file string     path to history file
      --interval duration       update interval (default 1h0m0s)
      --output monitor          monitor output (e.g. eDP-1)
      --phrase string           search phrase
      --resolution resolution   target resolution (e.g. 1920x1080) (default 0x0)
      --save-dir string         save directory (default "/home/$USER/Pictures/chiasma")
      --tool string             wallpaper tool (default "swaybg")
      --verbose                 enable verbose logs
```

## examples

**1. one-shot update using last google search from vivaldi:**
```bash
chiasma --browser vivaldi --api unsplash
```

**2. daemon mode (update every 30 minutes) using specific keyword:**
```bash
chiasma --phrase "cyberpunk city" --follow --interval 30m --tool swww
```

**3. specific monitor and resolution with nasa api:**
```bash
chiasma --output HDMI-A-1 --resolution 2560x1440 --api nasa
```

**4. firefox usage (requires manual history path):**
```bash
chiasma --browser firefox --history-file ~/.mozilla/firefox/PROFILE_ID/formhistory.sqlite
```

**5. chromium-based browser with custom history path:**
```bash
chiasma --browser brave --history-file ~/.config/BraveSoftware/Brave-Browser/Default/History
```

## supported providers

### browsers
*   **chromium-based**: `google-chrome`, `vivaldi`, `chromium`, `brave`, `opera`.
*   **firefox**: `firefox`.

### apis
*   **nasa**
*   **unsplash**

### tools
*   **swww**
*   **swaybg**

## todo

- [x] add swaybg tool
- [x] add swww tool
- [x] unsplash image provider
- [x] nasa image provider
- [x] chromium based browser support
- [x] firefox support
- [x] resolution auto detection
- [ ] refactor codebase
- [ ] add more providers
