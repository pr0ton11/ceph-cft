# Ceph Config Format Tool

This tool allows to configure your Ceph cluster simply by defining environment variables prefixed with *CEPH_*

Uses restrictions of https://docs.ceph.com/en/latest/rados/configuration/ceph-conf

## How does it work?

It iterates through the defined environment variables and:

* Detects if a variable if prefixed with Ceph
* Contains a valid configuration section (default: global)
* Supports fine grained configuration of specific individual daemons (e.g client.1) by replacing *__* with *.*
* Writes the updated ceph configuration file

## How do I use it

Just run ceph-cft before you start your ceph daemons and it will update / override the configuration values in ceph.conf based on environment variables.

## How do I define a custom path to ceph.conf?

This can be accomplished by setting the environment variable ``` CFT_CONFIG_PATH ```. Default is ```/etc/ceph/ceph.conf```.

## Examples

The following environment variables were set:

```
CEPH_GLOBAL_LOG_FILE='/var/log/ceph/$cluster-$type.$id.log'
CEPH_OSD_OP_QUEUE=wpq
CEPH_MON_LOG_TO_SYSLOG=true
CEPH_TEST_WITHOUT_SECTION=works
CEPH_CONTAINS_WHITESPACES="Hello World"
CEPH_OSD__1_OBJECTER_INFLIGHT_OPS=512
```

This leads to the generation of the following configuration file:

```
[global]
log_file             = /var/log/ceph/$cluster-$type.$id.log
test_without_section = works
contains_whitespaces = "hello world"

[osd]
op_queue = wpq

[mon]
log_to_syslog = true

[osd.1]
objecter_inflight_ops = 512
```


## Where is it used

This tool is used for the following projects:

* https://github.com/pr0ton11/radosgw
* https://github.com/pr0ton11/ceph-yocto
