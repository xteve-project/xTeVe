<a href="https://hub.docker.com/repository/docker/winguru/xteve-experimental"><img alt="Docker Cloud Automated build" src="https://img.shields.io/docker/cloud/automated/winguru/xteve-experimental?style=for-the-badge"></a>&nbsp;
<a href="https://hub.docker.com/repository/docker/winguru/xteve-experimental"><img alt="Docker Cloud Build Status" src="https://img.shields.io/docker/cloud/build/winguru/xteve-experimental?style=for-the-badge"></a>

<h1 id="xTeVe"><a href="https://https://github.com/winguru/xTeVe/tree/experimental">xTeVe Docker Experimental Edition</a></h1>
<tr>

Image Maintainer:  <b>winguru</b> <github@geoffthornton.me\></a>
<br>
                          
<table style="border: 0">
<tr style="border: 0">
<td style="border: 0; padding-left: 0">For issues, questions, comments, or suggestions regarding this experminal xTeVe build, click below to join my Discord server:</td>
<td style="border: 0">For general xTeVe docker questions or issues, click below to join the official xTeVe docker Discord server.</td>
</tr>
<tr style="border: 0">
<td style="border: 0; padding-left: 0"><a href="https://discord.gg/mqHtSBSSNC"><img alt="Discord" src="https://img.shields.io/discord/431853493503131648?style=for-the-badge"></a></td>
<td style="border: 0"><a href="https://discord.gg/Up4ZsV6"><img alt="Discord" src="https://img.shields.io/discord/465222357754314767?color=%2367E3FB&style=for-the-badge"></a></td>
</tr>
</table>
Many thanks to all the contributors to xTeVe, including <a href="https://github.com/dnsforge-repo/xteve">LeeD &lt;hostmaster@dnsforge.com&gt;</a>, whom most of this docker build is based on.

<br>
<br>                                                                                                                                      
                                                                                                                                      
<h2 id="description">Description</h2>

xTeVe is a M3U proxy server for Plex, Emby and any client and provider which supports the .TS and .M3U8 (HLS) streaming formats.

<p>xTeVe emulates a SiliconDust HDHomeRun OTA tuner, which allows it to expose IPTV style channels to software, which would not normally support it.  This Docker image includes the following packages and features:

<br>

<br>

<ul>
<li>xTeVe v2.1 (Linux) x86 64 bit</li>
<li>VLC & ffmpeg Support</li>
<li>Automated XMLTV Guide Lineups</li>
<li>Runs as an unprivileged user</li>
</ul>

<br>

<h2 id="experimental-features">Experimental Features</h2>                  
<ul>
<li>Per-channel support for "Tvg-shift" field in M3U files (timezone shifting)</li>
<li>Per-filter support for setting the starting channel for automatic channel number mapping/li>
<li>Per-filter support for optionally preserving the M3U channel number mappings/li>
<li>Per-filter support for automatic xTeVe dummy EPG data for channels with no matching EPG source</li>
</ul>

<h2 id="bugfixes">Bugfixes</h2>
<ul>
<li>Ensured HDHomeRunner lineup.json file is always sorted by channel number</li>
</ul>
<br>
                    
<h2 >Docker 'run' Configuration & container mappings</h2>

The recommended <b>default</b> container settings are listed in the docker run command listed below:


<p><b> docker run -it -d --name=xteve --network=host --restart=always -v $PATH/xteve:/home/xteve/conf winguru/xteve-experimental:latest</b></p>


<br>

<br>

<h2 >Isolated (bridge) mode</h2>
<p>To isolate the container in bridge mode use 'docker run' as follows.  **Only use bridge mode if you fully understand its use**  Generally for most users, it is easier to use host mode. 

<br>

In bridge mode the docker container will assign it's own dockernet ip address (usually in the 172.17.x network).</p>

<p><b>docker run -it -d --name=xteve -p 34400:34400 --restart=always -v $PATH/xteve:/home/xteve/conf winguru/xteve-experimental:latest</b></p>

<br>

<br>

<h2>Default container paths</h2>

This container is configured with the following default environmental variables,  for reference, here are the paths of the xTeVe installation:


<table class="paleBlueRows">
<thead>
<tr>
<th>Variable</th>
<th>Path</th>
</tr>
</thead>
<tfoot>
<tr>
<td>&nbsp;</td>
<td>&nbsp;</td>
</tr>
</tfoot>
<tbody>
<tr>
<td>$XTEVE_HOME</td>
<td>/home/xteve</td>
</tr>
<tr>
<td>$XTEVE_BIN</td>
<td>/home/xteve/bin</td>
</tr>
<tr>
<td>$XTEVE_CONF</td>
<td>/home/xteve/conf</td>
</tr>
<tr>
<td>$XTEVE_CONF/data</td>
<td>/home/xteve/conf/data</td>
</tr>
<tr>
<td>$XTEVE_CONF/backup</td>
<td>/home/xteve/conf/backup</td>
</tr>
</tbody>
</table>

<br>

<h2 id="parameters">Parameters</h2>

<table class="paleBlueRows">
<thead>
<tr>
<th>Parameter</th>
<th>Description</th>
</tr>
</thead>
<tfoot>
<tr>
<td>&nbsp;</td>
<td>&nbsp;</td>
</tr>
</tfoot>
<tbody>
<tr>
<td>--name</td>
<td>Name of container image</td>
</tr>
<tr>
<td>--network</td>
<td>Set network type [ host | bridge ]</td>
</tr>
<tr>
<td>--restart</td>
<td>Enable auto restart for this container</td>
</tr>
<tr>
<td>-e TZ=Europe/London</td>
<td>Set custom Locale</td>
</tr>
<tr>
<td>-p 34400</td>
<td>Default container port mapping [ 127.0.0.1:34400:34400 ]</td>
</tr>
<tr>
<td>-e XTEVE_PORT=8080</td>
<td>Set custom xTeVe Port</td>
</tr>
<tr>
<td>-e XTEVE_BRANCH=beta</td>
<td>Set xTeVe git branch [ master|beta ] Default: master
</tr>
<tr>
<td>-e XTEVE_DEBUG=1</td>
<td>Set xTeVe debug level [ 0-3 ] Default: 0=OFF</td>
</tr>
<tr>
<td>-e XTEVE_API=0</td>
<td>Enable/Disable API [ beta ] Default: 1=ON</td>
</tr>
<tr>
<td>-v</td>
<td>Set volume mapping [ -v ~xteve:/home/xteve/conf ]</td>
</tr>
<tr>
<td>winguru/xteve-experimental:latest</td>
<td>Latest Docker image</td>
</tbody>
</table>

<br>
<br>

