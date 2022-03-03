# gmatrix
cmatrix but go

## Installation

### Pre-requisites

 - A working istallation of go tools
 - A terminal emulator and font that supports plenty of unicode glyphs. Download and view charset.txt in your terminal to check if your terminal/font supports the glyphs used.

### Installation Instructions

 - Navigate to your `$(go env GOPATH)/src` folder and clone the repo.
 - Make sure you intstall the required dependencies of this project.
 - Compile gmatrix with
 ```bash
 cd gmatrix
 go install tjweldon/gmatrix@latest 
 ```
 
 ## Usage
 
 ```
 gmatrix [--dump-charset PATH]
 ```
 
 When called with no arguments gmatrix should fill your current shell with all that good matrix nonsense. Once you're done quit with Ctrl-C or Q.
 
 The `--dump-charset PATH` option dumps the full characterset used to the file specified by `PATH`.
