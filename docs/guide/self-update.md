# Self-Update System

SilentCast includes a built-in self-update mechanism that allows you to easily update to the latest version directly from GitHub releases.

## Features

- **Automatic Version Checking**: Check for newer versions from GitHub releases
- **Progress Reporting**: Real-time download progress with speed and ETA
- **Checksum Verification**: Ensures download integrity using SHA256 checksums
- **Atomic Updates**: Safe binary replacement with automatic rollback on failure
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Caching**: Intelligent caching to avoid excessive API calls
- **User Confirmation**: Interactive prompts before applying updates

## Commands

### Check for Updates

Check if a new version is available without downloading:

```bash
./silentcast --check-update
```

Force a fresh check (ignoring cache):

```bash
./silentcast --check-update --force-update-check
```

### Self-Update

Update to the latest version:

```bash
./silentcast --self-update
```

For automated environments (no prompts):

```bash
./silentcast --self-update --force-self-update
```

## Update Process

The self-update process follows these steps:

1. **Check for Updates**: Query GitHub API for latest release
2. **Version Comparison**: Compare current version with latest available
3. **Display Information**: Show version info, release notes, and download size
4. **User Confirmation**: Prompt for confirmation (unless --force flag is used)
5. **Download**: Download binary with progress reporting
6. **Checksum Verification**: Verify downloaded file integrity
7. **Backup Current**: Create backup of current executable
8. **Atomic Replace**: Replace current binary with new version
9. **Cleanup**: Remove temporary files and backup
10. **Restart**: Automatically restart application (platform-specific)

## Safety Features

### Atomic Updates

Updates are performed atomically:
- Current executable is backed up before replacement
- If update fails, the backup is automatically restored
- No risk of corrupted or partial installations

### Checksum Verification

All downloads are verified using SHA256 checksums:
- Checksums are fetched from GitHub release assets
- Downloaded files are verified before installation
- Update is aborted if checksum doesn't match

### Rollback Support

If an update fails:
- Original executable is automatically restored from backup
- Error messages provide clear information about the failure
- System remains in a working state

## Configuration

The updater can be configured through code or environment variables:

### Default Configuration

```go
cfg := updater.Config{
    CurrentVersion: version.GetVersionString(),
    RepoOwner:      "SphereStacking",
    RepoName:       "SilentCast",
    CheckInterval:  24 * time.Hour,    // Cache duration
    CacheDuration:  1 * time.Hour,     // API response cache
}
```

### Cache Behavior

- Update checks are cached for 1 hour by default
- Cache prevents excessive GitHub API calls
- Use `--force-update-check` to bypass cache
- Cache is stored in the user's config directory

## Platform-Specific Features

### Windows
- Uses `taskkill` for graceful process termination
- Supports Windows service restart
- Handles file locking appropriately

### macOS
- Uses `syscall.Exec` for seamless process replacement
- Supports LaunchAgent restart
- Handles code signing verification

### Linux
- Uses `syscall.Exec` for process replacement
- Supports systemd service restart
- Handles file permissions correctly

## API Integration

### GitHub Releases API

The updater integrates with GitHub's REST API:

```
GET https://api.github.com/repos/SphereStacking/SilentCast/releases/latest
```

Features:
- Rate limiting awareness
- Proper error handling
- Timeout configuration (30 seconds)
- User-Agent header for identification

### Asset Selection

The updater automatically selects the appropriate binary:
- Detects current platform (OS/architecture)
- Looks for exact platform matches in release assets
- Filters out archives (prefers direct binaries)
- Falls back gracefully if no suitable asset is found

## Progress Reporting

Real-time progress information includes:
- Download percentage
- Downloaded/total bytes (human-readable)
- Current download speed
- Estimated time remaining (ETA)
- Visual progress bar

Example output:
```
🔍 Checking for updates...

📦 Update Available!
  Current version: v0.1.0-alpha.7
  Latest version:  v0.1.0-alpha.8
  Published:       2024-01-20
  Download size:   12.4 MB

📋 Release Notes:
✨ Add self-update mechanism
🐛 Fix hotkey detection on Linux
📝 Update documentation

❓ Do you want to update now? (y/N): y

⬇️ Downloading update...
[████████████████████████████████████████] 100.0% 12.4 MB/12.4 MB @ 2.1 MB/s ETA: complete
🔐 Verifying checksum...
✅ Update applied successfully!
🔄 Restarting SilentCast...
```

## Error Handling

Common error scenarios and their handling:

### Network Issues
- Connection timeouts are handled gracefully
- Retry logic for transient failures
- Clear error messages for network problems

### GitHub API Errors
- 404 errors indicate repository or release not found
- Rate limiting is respected and reported
- API response validation

### File System Errors
- Permission errors during backup/replacement
- Disk space checks before download
- Temporary file cleanup on errors

### Checksum Failures
- Downloaded file is immediately deleted
- Clear error message about integrity failure
- Original executable remains untouched

## Troubleshooting

### Common Issues

**"No suitable update found for platform"**
- Check that GitHub releases include binaries for your platform
- Verify the asset naming convention matches expectations

**"Checksum verification failed"**
- Network corruption during download
- Retry the update operation
- Check internet connection stability

**"Permission denied"**
- Ensure SilentCast has write permissions to its directory
- On Linux/macOS, check file ownership and permissions
- Consider running with appropriate privileges

**"GitHub API returned status: 404"**
- Repository may not exist or be private
- Check network connectivity
- Verify repository owner/name configuration

### Debug Mode

Enable debug logging for detailed update information:

```bash
./silentcast --self-update --debug
```

This will show:
- Detailed API requests and responses
- File system operations
- Cache operations
- Platform-specific details

## Security Considerations

### Checksum Verification
- Always verifies SHA256 checksums when available
- Rejects downloads with mismatched checksums
- Uses secure hash algorithms

### HTTPS Only
- All downloads use HTTPS for encryption
- GitHub's SSL certificates are verified
- No fallback to insecure protocols

### Code Signing
- Future versions will support code signature verification
- Platform-specific signing validation
- Enhanced security for enterprise environments

## Development Notes

### Testing
- Comprehensive unit tests for all updater components
- Integration tests for the complete update workflow
- Mock GitHub API responses for reliable testing

### Build Integration
- CI/CD automatically generates checksums for releases
- Platform-specific binaries are built automatically
- Release artifacts follow consistent naming conventions

### Version Management
- Semantic versioning (semver) support
- Development version handling
- Pre-release and beta version support

## Future Enhancements

### Planned Features
- Differential updates for reduced download sizes
- Code signature verification
- Update channels (stable, beta, nightly)
- Background update checking
- Automatic updates (with user consent)

### Configuration Options
- Custom update channels
- Update frequency configuration
- Bandwidth limiting for downloads
- Proxy support for corporate environments

---

For implementation details, see the `internal/updater` package documentation.