# Runtime Layout (Isolated Services)

This directory is the isolated runtime root used by LocalValet.

Structure:

- `runtime/linux/<service>/<version>/` : service binaries and config per version
- `runtime/pids/` : pid files written by the app
- `runtime/logs/` : stdout/stderr logs from managed services

Example for PHP:

- `runtime/linux/php/8.2/sbin/php-fpm`
- `runtime/linux/php/8.2/etc/php-fpm.conf`
- `runtime/linux/php/8.3/sbin/php-fpm`
- `runtime/linux/php/8.3/etc/php-fpm.conf`
- `runtime/linux/php/8.4/sbin/php-fpm`
- `runtime/linux/php/8.4/etc/php-fpm.conf`

Active versions are controlled in `config/runtime.json`.

If you add a new version template, make sure the actual `php-fpm` binary exists and is executable.
