# version-helper

----------------

# Version Styles
## Client Team
0.3.2-feature-charactermove.313.38cdcab

# Server & DevOps Team
0.3.2-feature-charactermove.38cdcab

----------------

# Commands
## Client Commands
```shell
# Generate Version Metadata Client Team
$ versionhelper client raise
$ versionhelper client set <semver> <count>
# Client Specific Parameters
--project : Project Name (default "client")
--gen-header-file : Header Metadata File Name C++ header (default "GeneratedVersion.h")
--gen-version-file : Version Metadata File Name (json or yaml) (default "version.yaml")
-d, --tmpl-header-file : Template Header File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///GeneratedVersion.h")
-v, --tmpl-version-file : Template Version File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///client.yaml")
```

## Server and DevOps Commands
```shell
# Generate Version Metadata for Server Team
$ versionhelper server
# Server Specific Parameters
--override-project : Override Project Name (default is repository-name)
--gen-file : Output File Name (default version.yaml)
-t, --tmpl-file : Template File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///server.yaml")

# Generate Version Metadata for DevOps Team
$ versionhelper devops
# Server Specific Parameters
--override-project : Override Project Name (default is repository-name)
--gen-file : Output File Name (default version.yaml)
-t, --tmpl-file : Template File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///devops.yaml")
```

## latestversion Commands
```shell
# get or update latestversion data
$ versionhelper latestversion

# get 
Usage:
  versionhelper latestversion get {namePrefix} [flags]

Examples:
latestversion get mapdlc

Flags:
  -b, --branch string       branch name
  -c, --ci-hint CI 이름   CI 힌트 (default none)
  -h, --help                help for get
  -n, --namespace string    [REQUIRED] namespace name (ex: mapdlc-metadata)
  -r, --repo string         repository name

# update
Usage:
  versionhelper latestversion update {namePrefix} {version} [flags]

Examples:
latestversion update mapdlc abc123

Flags:
  -b, --branch string          branch name
  -c, --ci-hint CI 이름      CI 힌트 (default none)
  -h, --help                   help for update
  -l, --label stringToString   latest version label (ex: hello=world) (default [])
  -n, --namespace string       [REQUIRED] namespace name (ex: mapdlc-metadata)
  -r, --repo string            repository name
```

## Parameters
```yaml
--help : Help for versionhelper
--version : Print Program Version
--debug : Enable to Print Debug Messages (default false)
--json-log : When Enable this, Print Log Message as Json (default false)

--ci-hint : (default "", none means automatically detect what CI uses)
--counter : Build Revision Counter Type (local or redfox, default is local)
--counter-local-path : Local Counter DB File Path (default "~/.versionhelper/db.json")
--gen-dir : Output Files Directory (default is same as workdirectory)
```

## Edit Metadata Commands
```shell
$ versionhelper version append --platform=<> --target=<> --artifact-type=<> --uri=<> [--file=<>]
$ versionhelper version upload [--confilct-resolve-policy=<merge|overwrite> --conflict-retry=<int> --upload-fox-addr=<> --upload-fox-secure=<> --file=<>]
```

## cf
if you redfox server use, add environment set REDFOX_HOST.
