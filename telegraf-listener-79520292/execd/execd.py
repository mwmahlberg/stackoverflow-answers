#!/usr/bin/env python3
import json, logging, os, sys, time
from datetime import datetime, timezone

# Telegraf marks all output on stdout as an error
# and expects only errors to be reported on stderr.
LOGLEVEL = os.environ.get("LOGLEVEL", "ERROR").upper()

log_data = {
    "measurement": "log",
    "timestamp": datetime.now(timezone.utc).isoformat(),
    "message": "Test log",
    "status": 200,
}

if __name__ == "__main__":
    # We want to write logs to stderr, since telegraf
    # can only deal with properly formatted data sent to
    # the plugin's stdout.
    logging.basicConfig(level=LOGLEVEL,stream=sys.stderr)
    logger = logging.getLogger("my execd plugin")
    while True:
        try:
            log_data["timestamp"] = datetime.now(timezone.utc).isoformat()
        
            # Write the data to stdout, where telegraf can read it
            json.dump(log_data, sys.stdout)

            # The convention is "one record per line"
            # so we add a newline at the end
            # and flush the output to ensure telegraf sees it.
            sys.stdout.write("\n")
            sys.stdout.flush()
            time.sleep(1)
        except Exception as e:
            # Log any errors to stderr
            logger.error(f"Error: {e}")

