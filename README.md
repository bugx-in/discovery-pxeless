# Foreman PXE-less Installer (FPI)

**This project stills under development**

Helps DHCP-less bare metal infrastructure to be migrated.

It offers a way to automate the creation of Foreman Discovery ISO images when DHCP is not available in the network 
and mounting them into the server's console to start the installation.

Foreman discovery offers the solution for DHCP-less installations by creating your own discovery image for the specific server (https://theforeman.org/plugins/foreman_discovery/18.0/index.html). However, this process needs to be executed per server since the network details need to be embedded in the image.

This solutions automates the creation of the ISO images and it exposes them via HTTP so they can be mounted to the server's console.
