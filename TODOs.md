- [Rework Test's](#sec-1)
- [When packaging check PREFIX is honored](#sec-2)
- [Add OID to plan struct](#sec-3)
- [Branches](#sec-4)
  - [Don't ever hard code branches](#sec-4-1)
  - [Check branch when building](#sec-4-2)
- [Via initialization](#sec-5)
  - [Create an init via function](#sec-5-1)

\#+TITLE TODO's

# TODO Rework Test's<a id="sec-1"></a>

Test should only handle input's and output's. Many of the test's right now simple invoke function's and test for go error's.

# TODO When packaging check PREFIX is honored<a id="sec-2"></a>

We should check the sanity of PKGDIR.

We only record files in our manifest not directories. so an empty directory can end up in our tarball. But not be listed in the manifest this means we could try to untar a directory outside of our PREFIX.

We'll have to check the total entries in each of these directories

1.  PKGDIR
2.  /usr
3.  /usr/local

# TODO Add OID to plan struct<a id="sec-3"></a>

We should sha256sum check tarballs. So we can reference the OID when downloading and installing. This will ensure packages in publish match the package in plan git repository.

This gets tricky with the manifest in the tarball because there is no way to embed the OID into it. We'll have to update the manifest once the package is installed. The side benefit is we can then use OID to upgrade packages. And we don't have to worry about package increments or version increments.

# Branches<a id="sec-4"></a>

## TODO Don't ever hard code branches<a id="sec-4-1"></a>

We should never hard code branches they should be explicitly set in config.json

## TODO Check branch when building<a id="sec-4-2"></a>

when building we should check that plans git branch that is checkedout matches the configuration branch

# TODO Via initialization<a id="sec-5"></a>

When we first run via, it is dependent on the plans git repo for meta data. currently we git clone recursive the via repo, which contains the plans repo as well. We also do not respect the user and just blindly clone it on first run.

## TODO Create an init via function<a id="sec-5-1"></a>

do not assume and fetch the plans repository. Error gracefully and suggest user to init the plans repo.
