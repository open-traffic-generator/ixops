## new_gc

Setup prerequisites and create a new K8S cluster on Google Compute Engine

### Synopsis

* Does all the pre-checks to setup gcp cluster based on user's GLCOUD_EMAIL. Once gcp cluster is deployed, installs meshnet cni, metrics_server, ixia_c_operator and KNE.

### Note
* Make sure you set your email ID like so: echo "GCLOUD_EMAIL=your.email@example.com" >> ~/.profile
