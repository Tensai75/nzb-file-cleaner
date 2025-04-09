[![Release Workflow](https://github.com/Tensai75/nzb-file-cleaner/actions/workflows/build_and_publish.yml/badge.svg?event=release)](https://github.com/Tensai75/nzb-file-cleaner/actions/workflows/build_and_publish.yml)
[![Latest Release)](https://img.shields.io/github/v/release/Tensai75/nzb-file-cleaner?logo=github)](https://github.com/Tensai75/nzb-file-cleaner/releases/latest)

# NZB File Cleaner

Command line tool to manipulate the metadata and the filename of either a single NZB file or for batch processing all NZB files in a folder.

- can add the password from the filename (in the {{password}} format) as password meta tag to the NZB file
- can add the the filename as title meta tag to the NZB file
- can remove the titel or password meta tag from the NZB file
- can remove the password (in the {{password}} format) from the filename
- can add the password from the password meta tag to the filename (in the {{password}} format)
- can use the title from the title meta tag as the filename

## Usage

`nzb-file-cleaner.exe [--apm] [--apf] [--atm] [--utf] [--rpm] [--rpf] [--rtm] [--verbose] NZBFILE [DESTPATH]`

### Positional arguments:
- **`NZBFILE`**

  Path to the NZB file to be processed or a folder containing NZB files (required).

- **`DESTPATH`**

  Destination path where the new NZB file(s) should be saved (optional).
  If DESTPATH does not exist, the user is prompted if the path should be created.
  If DESTPATH is omitted, the new NZB file(s) will be written in the same directory as NZBFILE.
  **Caution**: NZB files with same filename are overwritten without warning. The use of a DESTPATH is recommended.

### Option arguments:
- **`--apm`**

  Add the password from the filename (if available in the {{password}} format) to NZB file metadata.
- **`--apf`**

  Add the password from NZB file metadata (if available) to the filename ({{password}}).
- **`--atm`**

  Add the filename to the NZB file metadata as title.
- **`--utf`**

  Use the title in the NZB file metadata (if available) as the filename for the NZB file.
- **`--rpm`**

  Remove the password from the NZB file metadata (if available).
- **`--rpf`**

  Remove the password from the filename (if available in the {{password}} format).
- **`--rtm`**

  Remove the title from the NZB file metadata (if available).
- **`--verbose, -v`**

  Enable verbose output.
- **`--help`, `-h`**

  Display a help text and exit.
- **`--version`**

  Display the program version and exit.

## Contribution

Feel free to send pull requests.

## Change log
#### v1.0.0
- first public release