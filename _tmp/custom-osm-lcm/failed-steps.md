LCM fails
---------
> epic!

So building & deploying a custom OSM LCM image has been lots of fun!
Here's what didn't work and possible workarounds.


### OSM 11 VM

Had to build it a couple of times. Some of the install script tasks
[failed][failed-osm-install] but the script went ahead. Eventually
I ended up with a broken OSM install in my hands---some OSM services
didn't get deployed to the K8s cluster. Not sure what the cause of
those random failure is, possibly lack of enough compute resources
and the install procedure not being robust enough to cater for slow
boxes?


### LCM build failures - part 1

The command to build the LCM image artifacts took about 50 mins and
I didn't get a clean build in the end:

```console
% devops/tools/local-build.sh --module common,IM,N2VC,RO,LCM,NBI stage-2
...
dpkg-deb: building package 'python3-n2vc' in '../python3-n2vc_11.0.0rc1.post36+g23c4455-1_all.deb'.
 dpkg-genbuildinfo
 dpkg-genchanges  >../n2vc_11.0.0rc1.post36+g23c4455-1_amd64.changes
dpkg-genchanges: info: including full source code in upload
 dpkg-source --after-build .
dpkg-source: info: using options from n2vc-11.0.0rc1.post36+g23c4455/debian/source/options: --extend-diff-ignore=\.egg-info$
dpkg-buildpackage: info: full upload (original source is included)
dist run-test: commands[3] | sh -c 'rm n2vc/requirements.txt'
____________________________________________________________ summary _____________________________________________________________
  dist: commands succeeded
  congratulations :)
renamed './deb_dist/python3-n2vc_11.0.0rc1.post36+g23c4455-1_all.deb' -> '/home/ubuntu/snap/qhttp/common/python3-n2vc_11.0.0rc1.post36+g23c4455-1_all.deb'
Directory /home/ubuntu/workspace/RO does not exist
% echo $?
1
```

Going ahead to the next step anyway, just in case the build failure
wasn't critical...

```console
% devops/tools/local-build.sh --module LCM stage-3
...
Step 14/57 : RUN curl $PYTHON3_OSM_LCM_URL -o osm_lcm.deb
 ---> Running in 92e6b11d10dc
curl: no URL specified!
curl: try 'curl --help' or 'curl --manual' for more information
The command '/bin/sh -c curl $PYTHON3_OSM_LCM_URL -o osm_lcm.deb' returned a non-zero code: 2
Failed to build lcm
```

Oh deary, deary. Maybe I shouldn't have gone ahead.


### LCM build failures - part 2

So it turns out the reason for this error message

> Directory /home/ubuntu/workspace/RO does not exist

is that the command

```console
% devops/tools/local-build.sh --module common,IM,N2VC,RO,LCM,NBI stage-2
```

tries to build an OSM component called RO. In fact there's an RO repo.
Since the command also tries building NBI, we're going to clone and set
up these two repos too:

```console
% git clone https://osm.etsi.org/gerrit/osm/RO
% git clone https://osm.etsi.org/gerrit/osm/NBI
% for r in IM LCM N2VC NBI RO common devops; do cp commit-msg $r/.git/hooks/; done
```

Now running again the build command got me past the directory error,
but the build seems to get into an infinite loop when installing RO
deps

```console
...
dist_ro_vim_vmware installdeps: -r/build/requirements.txt, -r/build/requirements-dist.txt
```

it just sits there for half an hour seemingly making no progress.
Could it be an issue with VMWare deps? Well, I killed the process
and ran the command again. And again the process got stuck on installing
deps

```console
...
dist_ro_vim_vmware installdeps: -r/build/requirements.txt, -r/build/requirements-dist.txt
...
dist_ro_sdn_odl_of installdeps: -r/build/requirements.txt, -r/build/requirements-dist.txt
```

Notice how this time the VMWare deps step succeeded while the build
got stuck on another component.




[failed-osm-install]: ./osm-install/install.failed.log


