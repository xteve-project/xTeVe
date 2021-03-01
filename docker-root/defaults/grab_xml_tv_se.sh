#!/bin/bash

echo "Running: tv_grab_eu_xmltvse"
echo "Configured days: ${XMLTV_DAYS}"
/usr/local/bin/tv_grab_eu_xmltvse --config-file /config/.xmltv/tv_grab_eu_xmltvse.conf --output /config/tv_grab_eu_xmltvse_guide.xml --offset -1 --days $XMLTV_DAYS
echo "grabber finished, exiting..."
exit 0