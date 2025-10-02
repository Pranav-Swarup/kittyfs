# KittyFS

A sleek terminal-based file explorer. Navigate your files and folders easily with a simple TUI interface, all while enjoying a cute cat-themed aesthetic.

<img width="1477" height="743" alt="kittyfs1" src="https://github.com/user-attachments/assets/04c3b51a-f66e-46b0-974e-40742a5e43db" />
<img width="1463" height="748" alt="kittyfs2" src="https://github.com/user-attachments/assets/eb17f348-7cbb-4ca3-880c-d98bcefc9621" />


## Installation

### Windows

1. Download the installer (`KittyFS_Installer.exe`) from the [Releases](https://github.com/Pranav-Swarup/kittyfs/releases) page.
2. Run the installer by clicking on it.
3. Launch KittyFS from the Start Menu or Desktop shortcut.

### Mac-OS

1. The `releases/Mac-OS` folder contains a `KittyFS.app` bundle as well as the zipped version of it
2. Download `KittyFS-mac.zip` and unzip it
3. Drag the unzipped folder: `KittyFS.app` to /Applications
4. The app should show up in Finder or Dock

#### If you want to build it from source:

1. Make sure you have [Go](https://golang.org/dl/) installed.
2. Clone this repository:

```bash
git clone https://github.com/Pranav-Swarup/kittyfs.git
cd kittyfs
cd sourcefiles
```

3. Run the `.exe` or Build the TUI binary:

```bash
go build -o KittyFS.exe
```


## Features

- Lightweight, terminal-based UI, no BS.
- Browse drives and directories in a terminal interface
- Navigate with **↑ ↓ ← →** 
- **Enter** to open folders or files
- **'o'** to Open the file/folder location
- **Backspace** to go up to parent directory
- Search for files using `/`
- Toggle extended help with `?`
- Change color themes with `t`
- Cross-platform


## Usage

| Key        | Action                                    |
| ---------- | ----------------------------------------- |
| ↑ / ↓      | Navigate files                            |
| Enter      | Open file/folder                          |
| Backspace  | Go to parent directory or drive selection |
| /          | Filter/search files                       |
| Esc        | Clear filter                              |
| ?          | Toggle extended help                      |
| t          | Change color theme                        |
| q / Ctrl+C | Quit                                      |

## Configuration

You can customize border and highlight colors by editing config.json (created automatically after first run):

``` js
{
  "border_color": "#FF69B4",
  "highlight_color": "#FF1493"
}
```

Press t inside the app to cycle themes and automatically save changes.

## Contributions

Feel free to:

- Report bugs

- Suggest new features

- Improve the TUI styling

- Add cross-platform support

## License

MIT License © Pranav Swarup

## Acknowledgements

Bubble Tea
 — Terminal UI framework

Lip Gloss
 — Styling in the terminal
