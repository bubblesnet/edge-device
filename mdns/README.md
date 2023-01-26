# mdns

mdns is a Balena container that exists solely to find the bubblesnet-controller and communicate that
knowledge to the other containers on the device. bubblesnet-controller.local is the host name of the 
bubblesnet controller and is registered with avahi mdns on the controller which means that any 
mdns device on the same subnet can resolve that name to the right IP address. However, the functional 
containers on the edge-device (sense-go, sense-python and store-and-forward) all run in bridge mode (172.X.X.X)
so they are NOT on the same subnet as bubblesnet-controller.local (e.g. 192.168.x.x). This protects the
functional containers from the network.

The key parts of this are:

 - The service file [find_controller.service](find_controller.service)
 - ThE line from [../docker-compose.yml](docker-compose.yml) "network_mode: host" that defines this service as running in the same subnet as bubblesnet-controller.local
 - The file /config/bubblesnet-controller.txt that contains the output of the service.

