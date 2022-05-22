
# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).




## [ 0.5.2 ] - 2022-5-22 13:55:28

> BENS-0027 optimize release workflow

### Changed

- `88fc417` - BENS-0027 split release creation and asset release into separate jobs



## [ 0.5.1 ] - 2022-5-22 13:25:17

> BENS-0026 fix update bug introduced in 0.5.0

### Changed

- `768f7ea` - BENS-0026 changed actions workflow to point to correct project version



## [ 0.5.0 ] - 2022-5-22 3:15:33

> BENS-0025 Simple-Push Feature

### Added

- `d450559` - BENS-0025 logging overhauled to be far more useful
- `cceb092` - BENS-0025 changelog debugging added properly
- `00e681c` - BENS-0025 refactored git interface and fixups to simple-push
- `27ce510` - BENS-0025 implemented base command for simple-push feature



## [ 0.4.2 ] - 2022-5-15 21:47:26

> BENS-0024 Fixup Related to BENS-0023

### Changed

- `25c03ea` - BENS-0024 fixed project root output formatting
- `1cfedb5` - BENS-0024 fixups to documentation and install script
- `4b425c3` - BENS-0024 changed order for project root check and refactored common functions



## [ 0.4.1 ] - 2022-5-15 19:50:48

> BENS-0023 Fix Feature Push in Nested Folder Structure

### Changed

- `9d756f1` - BENS-0023 ensured lead and trailing spaces removed from GOGDir
- `7200469` - BENS-0023 implemented fix for GOG feature discovery



## [ 0.4.0 ] - 2022-4-23 16:7:4

> BENS-0022 Version Prefix Warning and Prompt

### Added

- `f047226` - BENS-0022 added demo and detailed usage to README
- `3a6d62b` - BENS-0022 warning for prefix mismatch from existing



## [ 0.3.8 ] - 2022-4-18 16:1:2

> BENS-0021 Version Prefix Local Override

### Changed

- `06f6825` - BENS-0021 implemented warning for prefix change
- `42c4a9a` - BENS-0021 implemented warning when creating feature with diff. prefix than existing
- `a2f11ec` - BENS-0021 Test Build (0)
- `a77fa23` - BENS-0021 added version prefix override to feature create



## [ 0.3.7 ] - 2022-4-14 17:55:17

> BENS-0020 Added Changelog Entry Applied to Minor Releases

### Changed

- `1abd002` - BENS-0020 minor now counts as added functionality in changelogs



## [ 0.3.6 ] - 2022-4-14 17:40:45

> BENS-0019 Update Release Action Version

### Changed

- `b0f62c4` - BENS-0019 updated version for go-release action



## [ 0.3.5 ] - 2022-4-12 19:44:56

> BENS-0018 Fix Bad Release

### Changed

- `71f7425` - BENS-0018 removed bad changelog entry
- `1b4c27b` - BENS-0018 removed breaking prefix from config



## [ 0.3.4 ] - 2022-4-7 0:23:27

> BENS-0015 Fix CD Pkg Version Reference

### Changed

- `463a99c` - BENS-0015 implemented cd fix



## [ 0.3.3 ] - 2022-4-7 0:18:43

> BENS-0014 Install dependencies + App config

### Changed

- `4a1e872` - BENS-0014 cleaned up commands and removed some unused helpers
- `0cae945` - BENS-0014 fixed default config file
- `2559cbd` - BENS-0014 refactored packages and implemented config
- `85eeabb` - BENS-0014 added git to required install dependencies



## [ 0.3.2 ] - 2022-4-4 19:35:21

> BENS-0013 Update README + Install Scripts

### Changed

- `16f0263` - BENS-0013 additional helpful instructions post-install
- `4d03c18` - BENS-0013 fixed issue with latest version capture and added tested platforms
- `b21dcf6` - BENS-0013 Test Build (1)
- `b903be1` - BENS-0013 latest version one-liner for install script
- `69abdd4` - BENS-0013 Test Build (0)
- `5f2b9a4` - BENS-0013 created install script and instructions



## [ 0.3.1 ] - 2022-4-4 15:21:36

> BENS-0012 Fix Broken Update

### Changed

- `efbd56a` - BENS-0012 implemented asset_name override on CD



## [ 0.3.0 ] - 2022-4-4 12:59:49

> BENS-0011 Cmd Shortnames and Release Binary name

### Changed

- `76e9d6d` - BENS-0011 more accurate descriptions for major,minor,patch flags
- `9c449d9` - BENS-0011 changed release binary name and added subcmd shorthands



## [ 0.2.3 ] - 2022-4-3 23:37:6

> BENS-0010 Bug Fixes + Version Flag

### Changed

- `4c4d283` - BENS-0010 cleaned tag version parsing for update cmd
- `9e7958e` - BENS-0010 update README
- `89e7c8e` - BENS-0010 rem version from main.go
- `51dc64c` - BENS-0010 updating is now a little more clear
- `454b562` - BENS-0010 caught error in update release getter



## [ 0.2.2 ] - 2022-4-3 16:4:49

> BENS-0009 =-= No Release for '*.x' Versions

### Changed

- `7f8feb8` - BENS-0009 added exclusion to release triggers
- `c8f59aa` - BENS-0009 Start Feature



## [ 0.2.1 ] - 2022-4-3 16:0:13

> BENS-0008 =-= Fixup Usage Docs

### Changed

- `e4bad8d` - BENS-0008 made usage example more useful
- `ab9cc24` - BENS-0008 Start Feature



## [ 0.2.0 ] - 2022-4-3 15:45:40

> BENS-0007 GOG Updates

### Changed

- `260182b` - BENS-0007 fixed reponame on update pkg
- `819c5bc` - BENS-0007 fixed changelog spacing issues
- `62c3160` - BENS-0007 updated release workflow
- `6f341ab` - BENS-0007 more helpful changelog error message
- `f86c959` - BENS-0007 refactor and error handling for update
- `480f5c8` - BENS-0007 implemented first iteration of update command
- `8c11a11` - BENS-0007 deployment workflow for gog binaries
- `1a40fdb` - BENS-0007 added the updated command (no impl)
- `0d26617` - BENS-0007 removed unnecessary push from feature create
- `d28f73e` - BENS-0007 Start Feature



## [ 0.1.4 ] - 2022-3-26 21:10:18

> BENS-0006 Quick Fixup Changelog Format

### Changed

- `8695df0` - BENS-0006 Start Feature



## [ 0.1.3 ] - 2022-3-26 21:6:26

> BENS-0005 Input Help and Bug Fixes

### Changed

- `790ec35` - BENS-0005 better help and usage builtin
- `8c20133` - BENS-0005 Start Feature



## [ 0.1.2 ] - 2022-3-26 3:13:48

> BENS-0004 Use Flag Package Properly

### Changed

- `d56d957` - BENS-0004 implemented better argument handling
- `ec566b3` - BENS-0004 Start Feature



## [ 0.1.1 ] - 2022-3-25 20:57:32

> BENS-0003 Error handling on Semver

### Changed

- `cd50284` - BENS-0003 fixed version match regexp
- `a6cbef9` - BENS-0003 Start Feature



## [ 0.1.0 ] - 2022-3-22 23:30:15

> BENS-0002 Refactoring + Cleanup

### Changed

- `05338a2` - BENS-0002 reorg feature functions
- `f2eb399` - BENS-0002 reorg feature functions
- `5273ff2` - BENS-0002 Continued work to simply codebase
- `45ba6bd` - BENS-0002 Initial refactor
- `2c64182` - BENS-0002 Initial refactor