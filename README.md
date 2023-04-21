# PowerSlide CLI

## Get it

### Automatically (Linux/Mac)

`/bin/bash -c "$(curl -fsSL https://github.com/power-slide/cli/releases/latest/download/install.sh)"`

### Manually (Linux/Mac/Windows)

Go to the
[latest release page](https://github.com/power-slide/cli/releases/latest),
then click on the release asset that matches your OS & CPU to download it,
if needed mark it as executable, rename it to `pwrsl`, then place it in a
folder on your `$PATH`. Run `pwrsl setup` and follow the on screen insructions.

## Update it

The CLI will check for updates automatically and prompt you to update.
Run `pwrsl update --help` for more details.

## Build it

Be sure to have golang >= 1.20, `sed`, `jq`, `curl` and `make` installed.

```shell
git clone https://github.com/power-slide/cli.git pwrsl
cd pwrsl
make # Development build
make release # Release build for all platforms
make test_release # Release build for current platform
make clean # Clean up builds
```

## Deployment process

When a PR is merged to master the following pipeline happens:

- A new release and tag (version) are automatically created and pushed
- Binaries for that release are built and uploaded as release assets

The CLI checks for new versions via the GitHub latest release API,
and downloads updates via the release assets.
