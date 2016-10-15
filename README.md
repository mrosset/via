<div id="table-of-contents">
<h2>Table of Contents</h2>
<div id="text-table-of-contents">
<ul>
<li><a href="#orgheadline1">1. Description</a></li>
<li><a href="#orgheadline2">2. Installing the via binary</a></li>
<li><a href="#orgheadline3">3. Installing a package</a></li>
<li><a href="#orgheadline4">4. Build bash</a></li>
<li><a href="#orgheadline5">5. Installing the development group</a></li>
</ul>
</div>
</div>


# Description<a id="orgheadline1"></a>

Via is a systems package manager written in go language. It's primary purpose is
is to download and install binary packages.  It uses json based plans to
download and build the packages. The packages once installed do not effect the
host system. And can be run on most modern Linux systems.

# Installing the via binary<a id="orgheadline2"></a>

    go get bitbucket.org/strings/via

# Installing a package<a id="orgheadline3"></a>

    via install bash

# Build bash<a id="orgheadline4"></a>

    via -c build bash

# Installing the development group<a id="orgheadline5"></a>

    via install devel
