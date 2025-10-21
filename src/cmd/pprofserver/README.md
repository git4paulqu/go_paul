# PProf Server with Enhanced File Browser

This is an enhanced version of pprof that includes a modern file browser functionality for browsing, comparing, and analyzing .pprof files with persistent state management and improved user experience.

## Quick Start

```bash
# Build the executable
go build -o pprof-server ./pprofserver

# Run with default settings
./pprof-server

# Run with custom settings
./pprof-server -port 9090 -path /path/to/profiles -no_browser=false
```

## Command Line Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `-port` | `8088` | HTTP server port |
| `-no_browser` | `true` | Disable automatic browser opening |
| `-path` | `executable directory` | Default file browser path |
| `-help` | - | Show help information |

## Features

### 🗂️ **Advanced File Browser**
- **Hierarchical Navigation**: Browse local directories with expandable tree view
- **Smart Filtering**: Automatically filters directories containing .pprof files
- **Windows Drive Support**: Shows all available drive letters on Windows
- **Lazy Loading**: Efficient directory scanning with intelligent caching

### 📊 **Profile Management**
- **Base/Source Selection**: Select profiles for comparison analysis
- **Visual Indicators**: Color-coded selection markers (green for base, yellow for source)
- **File Validation**: Automatic file existence checking before processing
- **Flexible Work Mode**: Work with base-only or base+source comparison

### 💾 **State Persistence**
- **Session Memory**: Remembers selected directory and expanded folders
- **Profile Persistence**: Retains base and source selections across page refreshes
- **Cross-Navigation**: Maintains state when switching between views
- **File Change Detection**: Automatically detects and refreshes when files change

### 🎨 **Modern UI/UX**
- **Two-Row Layout**: Base and Source profiles displayed on separate lines
- **Aligned Paths**: Perfect alignment of profile paths for better readability
- **Gradient Buttons**: Modern button design with hover effects and animations
- **Color-Coded Icons**: Green folders, blue files for intuitive navigation
- **Responsive Design**: Optimized for different screen sizes

## Usage Examples

### Basic Usage
```bash
# Start server on default port 8088
./pprof-server

# Start on custom port
./pprof-server -port 9090

# Start with custom profile directory
./pprof-server -path /path/to/your/profiles

# Start with browser auto-opening
./pprof-server -no_browser=false
```

### Advanced Usage
```bash
# Custom configuration
./pprof-server -port 8080 -path /data/profiles -no_browser=false

# Show help
./pprof-server -help
```

## Web Interface

### File Browser View
- **Default Landing Page**: Opens directly to File Browser
- **Directory Navigation**: Click folders to expand/collapse
- **Profile Selection**: Click "Select base" or "Select source" on .pprof files
- **Manual Directory**: Use "Select Directory" button for custom paths

### Profile Analysis
- **Base Profile**: Required for analysis (enables Work button)
- **Source Profile**: Optional for comparison analysis
- **Work Button**: Processes selected profiles and redirects to graph view
- **Reset Button**: Clears all selections

### State Management
- **Auto-Save**: All selections and expanded states are automatically saved
- **Auto-Restore**: Page refresh or navigation returns to previous state
- **File Monitoring**: Detects file changes and refreshes content automatically

## Technical Details

### Architecture
- **Backend**: Go-based HTTP server with RESTful API
- **Frontend**: Modern JavaScript with CSS3 animations
- **State Storage**: In-memory server-side state management
- **File System**: Cross-platform file system access

### API Endpoints
- `GET /api/files` - Get directory contents
- `POST /api/setbase` - Set base profile
- `POST /api/setsource` - Set source profile
- `POST /api/work` - Process profiles
- `GET /api/getstate` - Get current state
- `POST /api/setdir` - Set selected directory
- `POST /api/setexpanded` - Set folder expansion state

### Performance Features
- **Lazy Loading**: Directories loaded on-demand
- **Caching**: Directory modification time caching
- **Efficient Scanning**: Bounded recursive directory scanning
- **Memory Management**: Optimized state storage

## Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd go_paul

# Build the executable
cd src/cmd
go build -o pprof-server ./pprofserver

# Run the server
./pprof-server
```

## Requirements

- **Go**: Version 1.19 or later
- **Browser**: Modern web browser with JavaScript support
- **OS**: Windows, macOS, or Linux

## Troubleshooting

### Common Issues
1. **Port Already in Use**: Change port with `-port` argument
2. **Permission Denied**: Ensure executable has read access to target directories
3. **Browser Not Opening**: Use `-no_browser=false` to enable auto-opening
4. **Files Not Found**: Verify file paths and permissions

### Debug Mode
```bash
# Run with verbose output
./pprof-server -port 8088 2>&1 | tee server.log
```

## Contributing

This is a modified version of the Google pprof tool. Contributions and improvements are welcome.

## License

Based on the original pprof tool license. See LICENSE file for details.
