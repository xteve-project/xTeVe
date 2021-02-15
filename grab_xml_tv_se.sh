#!/usr/bin/bash

FILE=/home/xteve/.xmltv/tv_grab_se_tvzon.conf
OUTPUT=/home/xteve/xmltvse/xmltvse_guide.xml
echo 'Grabber configuration file is: '$FILE
if test -f "$FILE"; then
    if [[ -z "$XMLTV_DAYS" ]]; then
        XMLTV_DAYS="7"
    fi
    echo Running: /usr/bin/tv_grab_se_tvzon --config-file ${$FILE} --output ${OUTPUT} --days ${XMLTV_DAYS}
    /usr/bin/tv_grab_se_tvzon --config-file ${$FILE} --output ${OUTPUT} --days ${XMLTV_DAYS}
    echo "grabber finished, exiting..."
else
    echo "$FILE does not exist"
fi
exit 0