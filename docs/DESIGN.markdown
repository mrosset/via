Objectives
==============

1. Light weight meta/build data 
2. Meta data must be easily parsed by non shell
3. Support for sub packages, i.e. gcc gcclibs
4. Download interfaces for different download types, git,bz2,tar,hg,bzr
5. Decompress interfaces for different decompression types, git,bz2,tar,hg,bzr
7. Minimal dependencies


Logging
==============
1. Output must be simple but informative
2. verbose output is buffered and backtraced on errors


Paths
==============
1. only 2 paths, chroot and cache
2. all paths are references as full paths under /chroot


Cache dir structure
--------------
cache
├── build
│   └── foo-0.0.1
├── package
│   └── foo-0.0.1
├── sources
└── stage
    └── foo-0.0.1


Bootstrap
==============
1. bintuils			/tools
2. gcc				/tools
3. linux-headers	/tools
4. glibc			/tools
5. adjust toolchain
