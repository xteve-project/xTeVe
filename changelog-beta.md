#### 2.0.3.0042-beta
**Version 2.0.3.0042 changes the settings.json.**  
Settings from the current beta can not be used for the current master version 2.0.3  
- New default options for VLC and FFmpeg  
- VLC and FFmpeg log entries in the xTeVe log  
- Less CPU load with VLC and FFmpeg  

#### 2.0.3.0035-beta
```diff
+ FFmpeg support
+ VLC support
```
**Version 2.0.3.0035 changes the settings.json.**  
Settings from the current beta can not be used for the current master version 2.0.3

#### 2.0.2.0024-beta
```diff
+ Improved monitoring of the buffer process
+ Update the XEPG database a bit faster
```

##### Fixes
- Error message if filter rule is missing
- Channels are lost when saving again (Mapping)
- Plex log, invalid source: IPTV

#### 2.0.1.0012-beta
```diff
+ Add support for "video/m2ts" video streams (Pull request #14)
```
#### 2.0.1.0011-beta
```diff
+ Original group title is shown in the Mapping Editor
```
##### Fixes
- incorrect original-air-date

#### 2.0.1.0010-beta
```diff
+ Set timestamp to <episode-num system="original-air-date">
```

#### 2.0.0.0008-beta
##### Fixes
- Pull request #6 [Error in http/https detection] window.location.protocol return "https:", not "https://"

#### 2.0.0.0007-beta
```diff
+ Buffer HLS: Add VOD tag from M3U8
+ CLI: Add new arguments [-restore]
+ CLI: Add new arguments [-info]
```
##### Fixes
- Missing images with caching for localhost URL


#### 2.0.0.0001-beta
```diff
+ Wizard: Add HTML input placeholder (M3U, XMLTV)
+ Wizard: Alert by empty value (M3U, XMLTV)
+ Image caching: Ignore invalid image URLs
```