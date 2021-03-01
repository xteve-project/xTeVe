#!/bin/bash

URL='http://www.xmltvepg.nl'

while IFS='' read -r line || [[ -n "$line" ]]; do
 echo "Running for file: ${line}"
 echo "Will fetch for: ${URL}/${line}"
 curl -s -o /tmp/$line -L $URL/$line
 xz -d -c /tmp/$line > /config/$line.xml
 rm /tmp/$line
 echo "XML written to /config/${line}.xml"
done < "$1"
exit 0
