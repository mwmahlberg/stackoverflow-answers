# Code for my answer to [How to deploy a local image on K3D without pushing to a registry and upgrade Helm deployments locally?][question]

## Prerequisites

1. [Docker][docker:install] or any other container runtime compatible with [tilt][tilt]
2. [k3d][k3d:install], a lightweight way to run Kubernetes inside docker
3. [tilt][tilt:install], a development runtime for Kubernetes... stuff. ;)

## Usage

### Linux / macOS

```shell
$ git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
Cloning into 'mwmahlberg-so-answers'...
[...]
$ cd mwmahlberg-so-answers/k3d-localregistry-tilt-79540452
$ k3d cluster create --config k3d.yaml
INFO[0000] Using config file k3d.yaml (k3d.io/v1alpha5#simple) 
INFO[0000] portmapping '8080:80' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] portmapping '8443:443' targets the loadbalancer: defaulting to [servers:*:proxy agents:*:proxy] 
INFO[0000] Prep: Network                                
INFO[0000] Re-using existing network 'k3d-demo' (1b7fb303b1804e61a955fa7c82cadd20455384ef8d06caed097c9d71879fd8bb) 
INFO[0000] Created image volume k3d-demo-images         
INFO[0000] Creating node 'localregistry'                
INFO[0000] Successfully created registry 'localregistry' 
INFO[0000] Starting new tools node...                   
INFO[0001] Starting node 'k3d-demo-tools'               
INFO[0001] Creating node 'k3d-demo-server-0'            
INFO[0001] Creating node 'k3d-demo-agent-0'             
INFO[0002] Creating LoadBalancer 'k3d-demo-serverlb'    
INFO[0002] Using the k3d-tools node to gather environment information 
INFO[0002] Starting new tools node...                   
INFO[0003] Starting node 'k3d-demo-tools'               
INFO[0004] Starting cluster 'demo'                      
INFO[0004] Starting servers...                          
INFO[0005] Starting node 'k3d-demo-server-0'            
INFO[0014] Starting agents...                           
INFO[0015] Starting node 'k3d-demo-agent-0'             
INFO[0020] Starting helpers...                          
INFO[0020] Starting node 'localregistry'                
INFO[0020] Starting node 'k3d-demo-serverlb'            
INFO[0027] Injecting records for hostAliases (incl. host.k3d.internal) and for 5 network members into CoreDNS configmap... 
INFO[0030] Cluster 'demo' created successfully!         
INFO[0030] You can now use it like this:
kubectl cluster-info
$ tilt up
Tilt started on http://localhost:10350/
v0.34.0, built 2025-03-12

(space) to open the browser
(s) to stream logs (--stream=true)
(t) to open legacy terminal mode (--legacy=true)
(ctrl-c) to exit
```

[question]: https://stackoverflow.com/questions/79540452/how-to-deploy-a-local-image-on-k3d-without-pushing-to-a-registry-and-upgrade-hel "Original question on StackOverflow"
[docker:install]: https://docs.docker.com/engine/install/ "Installation instructions for docker"
[k3d:install]: https://k3d.io/stable/#releases "Installation instructions for k3d"
[tilt]: https://tilt.dev "Tilt project page"
[tilt:install]: https://docs.tilt.dev/install.html "Installation instructions for tilt"
