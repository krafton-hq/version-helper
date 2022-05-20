# version-helper

----------------

# Version Styles
## Client Team
0.3.2-feature-charactermove.313.38cdcab

# Server & DevOps Team
0.3.2-feature-charactermove.38cdcab

----------------

# Commands
## Base Command
```shell
# Generate Version Metadata Client Team 
$ versionmaker client raise
$ versionmaker client set <semver> <count>
# Client Specific Parameters
--project : Project Name (default "client")
--gen-header-file : Header Metadata File Name C++ header (default "GeneratedVersion.h")
--gen-version-file : Version Metadata File Name (json or yaml) (default "version.yaml")
-d, --tmpl-header-file : Template Header File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///GeneratedVersion.h")
-v, --tmpl-version-file : Template Version File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///client.yaml")

# Generate Version Metadata for Server Team 
$ versionmaker server
# Server Specific Parameters
--override-project : Override Project Name (default is repository-name)
--gen-file : Output File Name (default version.yaml)
-t, --tmpl-file : Template File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///server.yaml")

# Generate Version Metadata for DevOps Team 
$ versionmaker devops
# Server Specific Parameters
--override-project : Override Project Name (default is repository-name) 
--gen-file : Output File Name (default version.yaml)
-t, --tmpl-file : Template File Url (embded:///PATH, ./PATH or file:///PATH) (default "embed:///devops.yaml")
```
## Parameters
```yaml
--help : Help for versionhelper
--version : Print Program Version
--debug : Enable to Print Debug Messages (default false)
--json-log : When Enable this, Print Log Message as Json (default false)

--ci-hint : (default "", none means automatically detect what CI uses)
--counter : Build Revision Counter Type (local or network, default is local)
--counter-local-path : Local Counter DB File Path (default "~/.versionhelper/db.json")
--counter-fox-addr : Network Counter gRPC Server Address (default "")
--counter-fox-secure : Use Tls Flag to Connect Network Counter Server (default true)
--gen-dir : Output Files Directory (default is same as workdirectory)
```
