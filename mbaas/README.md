# Environment variables

 - Need to activate virtual environment where Python code is set up,
   to get access to SKIDL.

 - Need to set `MBAAS_BASE_DIR` to point to the directory containing
   the `default_rules.json` file.

 - Need to set the `KICAD_SYMBOL_DIR` variable.


# Deployment

## File layout

```
/opt/morse-blinkies
  bin/
    mbaas
    start-mbaas
  python/
    process-mbaas.py
    requirements.txt
    .venv
  lib/
    default_rules.json
    project-template/
      <template files>
    kicad/
      <KiCad library files>
    assets/
      <web app assets>
    info/
      <static pages>
  systemd/
    mbaas.service
    mbaas.env
    clean-mbaas-output
```

## Procedure

1. Run `make-install-root` script.
2. Create `/opt/mbaas` directory on deployment host and untar
   installation bundle.
3. Create directories referenced in `systemd/mbaas.env` file.
4. Create Python virtual environment in `/opt/mbaas/python` and do
   `pip install -r requirements.txt` there.
5. Install Espresso (from
   https://github.com/classabbyamp/espresso-logic) and `zip`.
6. Set up the `mbaas` systemd service: environment file in
   `/etc/mbaas.env`, plus `mbaas.service` service definition.
7. Add nginx reverse proxy for `mbaas.skybluetrades.net`, including
   setting the right port for the MBaaS service to run on.
8. Start the `mbaas` systemd service.
9. Install the `clean-mbaas-output` script as a daily CRON script for
   the `mbaas` user.
