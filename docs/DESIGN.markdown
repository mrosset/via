vbuild (Build System)
==============

1. Prototype bash. Possibly go later
1. Light weight meta/build data, with minimal logic
2. Meta data must be easily parsed by non shell
3. Support for sub packages, i.e. gcc gcclibs
4. Simple interfaces to different build stages
7. Minimal dependencies

Build interfaces
--------------
1. build
2. stage
3. download
4. package

Logging
--------------
8. Output must be simple but informative
9. Verbose output is buffered and back traced on errors

Paths
--------------
1. only 2 paths plans and cache
2. all paths are references as full paths

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

