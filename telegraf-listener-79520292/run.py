import logging
import os
import requests
import time
from datetime import datetime, timezone

TELEGRAF_URL = os.environ.get("TELEGRAF_URL", "http://localhost:8080/telegraf")
LOGLEVEL = os.environ.get("LOGLEVEL", "INFO").upper()

log_data = {
    "measurement": "log",
    "timestamp": datetime.now(timezone.utc).isoformat(),
    "message": "Test log",
    "status": 200,
}

if __name__ == "__main__":
    logging.basicConfig(level=LOGLEVEL)
    logger = logging.getLogger("agent")
    logger.info(f"Telegraf URL: {TELEGRAF_URL}")
    while True:
        log_data["timestamp"] = datetime.now(timezone.utc).isoformat()
        try:
            # Ensure the session is closed after use.
            # See https://requests.readthedocs.io/en/latest/user/advanced/#session-objects
            # for more details.
            with requests.Session() as s:
                s.headers.update({"Content-Type": "application/json"})

                # Dito for the response.
                with s.post(TELEGRAF_URL, json=log_data) as response:
                    if response.status_code != 204:
                        raise Exception(
                            f"Received unexpected status {response.status_code}: {response.text}"
                        )
        except Exception as e:
            logger.error(f"Error sending data: {e}")
        time.sleep(1)
