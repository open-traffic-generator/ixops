## ixops

ixOps is `Ixia-C` Operations

### Synopsis

ixOps is the easiest way to manage and test emulated network topologies involving `Ixia-C` and containerized DUTs.

ixOps helps you create and destroy, emulated network topologies from the command line. GCP (Google Cloud Platform) is currently supported.

### Options

```
Usage ./ixops.sh [subcommand]:
new_gc            -   Setup prerequisites and create a new K8S cluster on GCE
new_kc            -   Setup prerequisites and create a new K8S cluster on local host using KIND
ct [topopath]     -   Create KNE topology (uses otg-dut-otg topology when topopath is missing)
dt [topopath]     -   Delete KNE topology (uses otg-dut-otg topology when topopath is missing)
new_fp            -   Setup featureprofiles on local host
ls_fp             -   List relevant tests in featureprofiles
run_fp [testpath] -   Execute given test from featureprofiles
new_tc            -   Setup K8S pod for internal tests (requires Keysight network)
ls_tc             -   List relevant tests in internal tests
run_tc [testpath] -   Execute given test from internal tests (execute all if testpapth is not provided)
rm_gc             -   Teardown K8S cluster on GCE
rm_kc             -   Teardown K8S cluster on KIND
setup_pre_gc      -   Setup prerequisites ONLY for creating K8S cluster on GCE
setup_pre_kc      -   Setup prerequisites ONLY for creating K8S cluster on KIND
```

* [new_gc](new_gc.md) - Setup prerequisites and create a new K8S cluster on Google Compute Engine
* [new_kc](new_kc.md) - Setup prerequisites and create a new K8S cluster on local host using KIND
* [ct](ct.md) - Create KNE topology 
* [dt](dt.md) - Delete KNE topology
* [new_fp](new_fp.md) - Setup featureprofiles on local host
* [ls_fp](ls_fp.md) - List relevant tests in featureprofiles
* [run_fp](run_fp.md) - Execute given test from featureprofiles
* [new_tc](new_tc.md) - Setup K8S pod for internal tests
* [ls_tc](ls_tc.md) - List relevant tests in internal tests
* [run_tc](run_tc.md) - Execute given test from internal tests
* [rm_gc](rm_gc.md) - Teardown K8S cluster on GCE
* [setup_pre_gc](setup_pre_gc.md) - Setup prerequisites ONLY for creating K8S cluster on GCE
* [setup_pre_kc](setup_pre_kc.md) - Setup prerequisites ONLY for creating K8S cluster on KIND


