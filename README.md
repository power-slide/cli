# PowerSlide CLI

## Get it

### Manually (Linux/Mac/Windows)

Go to the
[latest release page](https://github.com/power-slide/cli/releases/latest),
then click on the release asset that matches your OS & CPU to download it,
if needed mark it as executable, rename it to `pwrsl`, then place it in a
folder on your `$PATH`. Run `pwrsl setup` and follow the on screen insructions.

## Build it

Be sure to have golang >= 1.20 and `make` installed.

```shell
git clone https://github.com/power-slide/cli.git pwrsl
cd pwrsl
make # Development build
make clean # Clean up builds
```

## Deployment process

When a PR is merged to master the following pipeline happens:

- A new release and tag (version) are automatically created and pushed
- Binaries for that release are built and uploaded as release assets

The CLI checks for new versions via the GitHub latest release API,
and downloads updates via the release assets.
