OSM installation issues
-----------------------
> ...documenting the struggle for posterity.

I've captured some of the OSM install sessions that resulted in broken
installs and parked them in this dir just in case I ran into the same
issues again I've got an idea how to fix stuff quickly. One sore point
was the PGP keys---quite a few of them. Here's an example of how to fix
the K8s ones, the others are similar:

- https://stackoverflow.com/questions/49877401

Also notice the install script was broken. See the `multipass*` and
`patched*` scripts in the `osm-install` dir. It looks like the OSM
guys fixed it though:

- https://osm.etsi.org/gitlab/osm/devops/-/commit/fdbe776e9bb9e43f7d4dc0f8c023b93d258666e2

