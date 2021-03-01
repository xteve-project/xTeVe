#!/bin/bash

echo "Running: tv_grab_eu_xmltvse"
echo "Configured days: ${XMLTV_DAYS}"
/usr/local/bin/tv_grab_eu_xmltvse --config-file /config/.xmltv/tv_grab_eu_xmltvse.conf --quiet --output /config/tv_grab_eu_xmltvse_guide.xml --days $XMLTV_DAYS
echo "Grabber finished, exiting..."
exit 0