# Jenkins Slaves Cleaner

This tool is to clean the Jenkins slaves with specified labels for multiple Jenkins masters.

# Use Case

There is a Container Slave Pool(CSP) for CI/CD jobs of Jenkins masters. These slaves are automatically generated from Kubernetes cluster by Jenkins Kubernetes plugin. If the slaves need to be updated for some reasons, such as image or configure update, ALL the existence slaves in the CSP need to be deleted first. This tool is to clean the old slaves in the CSP with conditions:

* With the specified labels
* In idle state