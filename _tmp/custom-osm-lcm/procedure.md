Custom OSM LCM image
--------------------
> what a schlep!

Below are the steps to build and deploy a custom OSM LCM Docker image.
Mostly what OSM devs told Gabriele to do, plus some guesswork. Make
sure to keep hydrated b/c the procedure takes a few hours(1) and will
make you sweat alot :-)

(1) my hardware: MacBook Pro 13'', 2 GHz Dual-Core Intel Core i5,
8 GB RAM. Make sure to shut down every app since the below procedure
needs alot of horse power to run decently.

Notice at the moment we still can't get all the steps below to work.
Details [over here][fails].


### Build OSM 11 VM

We'll build and deploy our custom LCM image in an OSM release 11 VM.
Not explicitly mentioned by the OSM devs, but I don't see any other
easy way of doing that given I've got no clue about how they set up
their dev env.

```console
$ multipass launch --name osm11 --cpus 2 --mem 6G --disk 40G 20.04
$ multipass shell osm11
% wget https://osm-download.etsi.org/ftp/osm-11.0-eleven/install_osm.sh
% chmod +x install_osm.sh
% ./install_osm.sh 2>&1 | tee install.log
```

See:

- https://osm.etsi.org/docs/user-guide/latest/03-installing-osm.html


### Set up source workspace

OSM devs say:

> Clone these repositories in your workspace on OSM host:
> cd workspace

So that probably means create a `workspace` directory in your home on
the OSM VM you've just built...


```console
$ multipass shell osm11
% mkdir workspace && cd workspace
```

Cloning repos

```console
% git clone https://osm.etsi.org/gerrit/osm/LCM
% git clone https://osm.etsi.org/gerrit/osm/N2VC
% git clone https://osm.etsi.org/gerrit/osm/devops
% git clone https://osm.etsi.org/gerrit/osm/common
% git clone https://osm.etsi.org/gerrit/osm/IM
% git clone https://osm.etsi.org/gerrit/osm/RO
% git clone https://osm.etsi.org/gerrit/osm/NBI
```

Notice the RO and NBI repos weren't in the original instructions
they gave us, but then their build command requires them. So we
clone those two as well.

Setting up OSM's git commit hook in each repo


```console
% curl https://osm.etsi.org/gerrit/tools/hooks/commit-msg > commit-msg
% chmod +x commit-msg
% for r in IM LCM N2VC NBI RO common devops; do cp commit-msg $r/.git/hooks/; done
```


### Install additional deps

OSM devs mentioned you've got to install QHttp too.

```console
% devops/tools/local-build.sh --install-qhttpd
Attempting to open the browser failed, but the server might still work
This might happen if you're running this with sudo, a none graphical session, are lacking xdg-desktop portal support or have disabled the desktop interface
Attempting to serve files from /home/ubuntu/snap/qhttp/common, press control + c to exit
Serving HTTP on 0.0.0.0 port 8000 (http://0.0.0.0:8000/) ...
```

Hit `Ctrl+c` to exit.


### Build LCM image

First you've got to build the artifacts that make up the LCM image

```console
% devops/tools/local-build.sh --module common,IM,N2VC,LCM,NBI stage-2
```

Notice the original build command they gave us included RO too:

```console
% devops/tools/local-build.sh --module common,IM,N2VC,RO,LCM,NBI stage-2
```

but it looks like trying to build RO is a lost cause. Details
[over here][fails]. So we skip building RO for the moment.

Then build a Docker image from the above components. The image name is
`opensourcemano/lcm:devel`.
 
```console
% devops/tools/local-build.sh --module LCM stage-3
```


### Deploy LCM image

Finally, patch your OSM deployment to use the dev image you've just
built:

```console
% kubectl -n osm patch deployment lcm --patch '{"spec": {"template": {"spec": {"containers": [{"name": "lcm", "image": "opensourcemano/lcm:devel"}]}}}}'
deployment.apps/lcm patched
```

And as a sanity check:

```console
% kubectl -n osm get deployment lcm -o yaml | grep 'image: open'
        image: opensourcemano/lcm:devel

% kubectl -n osm get pod | grep lcm
lcm-7cf9644d9b-zthgf            1/1     Running   0              2m33s
```

### Grief down the line?

Notice we didn't build RO earlier. While we manage to build and deploy
LCM in the end, the LCM image might have some missing components, i.e.
those the build process supposedly fetched from RO. So we've got to
test the custom image thoroughly to make sure it works for our use
case.


### From the horse's mouth

For the record, these are the actual instructions Gabriele got from
the OSM devs. Copy-paste from the chat, original text, no edits.

After installing OSM via the script on the appropriate VM, this is
how to build the LCM image:

1. Clone these repositories in your workspace on OSM host:

cd workspace

git clone "https://osm.etsi.org/gerrit/osm/LCM" && (cd "LCM" && curl https://osm.etsi.org/gerrit/tools/hooks/commit-msg > .git/hooks/commit-msg ; chmod +x .git/hooks/commit-msg)

git clone "https://osm.etsi.org/gerrit/osm/N2VC" && (cd "N2VC" && curl https://osm.etsi.org/gerrit/tools/hooks/commit-msg > .git/hooks/commit-msg ; chmod +x .git/hooks/commit-msg)

git clone "https://osm.etsi.org/gerrit/osm/devops" && (cd "devops" && curl https://osm.etsi.org/gerrit/tools/hooks/commit-msg > .git/hooks/commit-msg ; chmod +x .git/hooks/commit-msg)

git clone "https://osm.etsi.org/gerrit/osm/common" && (cd "common" && curl https://osm.etsi.org/gerrit/tools/hooks/commit-msg > .git/hooks/commit-msg ; chmod +x .git/hooks/commit-msg)

git clone "https://osm.etsi.org/gerrit/osm/IM" && (cd "IM" && curl https://osm.etsi.org/gerrit/tools/hooks/commit-msg > .git/hooks/commit-msg ; chmod +x .git/hooks/commit-msg)

 

2. Install HTTP server:

devops/tools/local-build.sh --install-qhttpd

 

3. Build artifacts:

devops/tools/local-build.sh --module common,IM,N2VC,RO,LCM,NBI stage-2

 

4. Build image (this generates a “devel” tagged image using previous artifacts):

devops/tools/local-build.sh --module LCM stage-3

 

5. Patch deployment to use “devel” image:

kubectl -n osm patch deployment lcm --patch '{"spec": {"template": {"spec": {"containers": [{"name": "lcm", "image": "opensourcemano/lcm:devel"}]}}}}'




[fails]: ./failed-steps.md
