vbuild (Build System)
==============
1. Prototype bash. Possibly go later
1. Light weight meta/build data, with minimal logic
2. Meta data must be easily parsed by non shell
3. Support for sub packages, i.e. gcc gcclibs
4. Simple interfaces to different build stages
7. Minimal dependencies

Build interfaces
---------------
1. download
2. verify
3. stage
4. build
5. package

build interfaces are simple function hooks. In most cases these do not
need to be defined. Usually there is a easy way to check if the majority 
hooked should be called. function hooks make it easy to deal with corner cases,
while still keeping plan logic and meta data to a minimum.

There are some cases where there are more corner cases then other. Ie build 
and pakckage. In those cases we can use predefined custom functions.

Variables
---------
When sourcing plan we have one point of entry function called source\_plan.

Logging
--------------
8. Output must be simple but informative
9. Verbose output is buffered and back traced on errors

Cache dir structure
--------------
cache
	\build
		\foo-0.0.1
	\package
		\foo-0.0.1
	\sources
		\foo-0.0.1.tar.gz
	\stage
		\foo-0.0.1

Via (Package Manager)
==============
1. go
2. All metadata is kept in json format
2. Package signing, file hashing. Are first rate features.

