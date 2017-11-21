Take a valid OpenAPI spec file and turn it into a bunch of yaml files that you `kubectl apply` in order to have a working Fissionized API
Below is an example 
```
[centos@ip-10-0-1-162 hackathon]$ ./APIrance --env nodejs generate ./gringotts-spec.yaml
Current working directory is /home/centos/hackathon
cmd: /usr/local/bin/fission env list
2017/11/21 19:23:58 Environment 'nodejs' already present, skipping creation
Created dir: /home/centos/hackathon/fission/DELETE
Created dir: /home/centos/hackathon/fission/POST
Created dir: /home/centos/hackathon/fission/GET
Created dir: /home/centos/hackathon/fission/PUT
cmd: /usr/local/bin/fission function create --name getticketsla --env nodejs --code fission/GET/getTicketSla.js
cmd: /usr/local/bin/fission function create --name updateticket --env nodejs --code fission/PUT/updateTicket.js
cmd: /usr/local/bin/fission function create --name deleteticket --env nodejs --code fission/DELETE/deleteTicket.js
cmd: /usr/local/bin/fission function create --name saveticketsla --env nodejs --code fission/POST/saveTicketSla.js
cmd: /usr/local/bin/fission function create --name searchtickets --env nodejs --code fission/GET/searchTickets.js
cmd: /usr/local/bin/fission function create --name saveticket --env nodejs --code fission/POST/saveTicket.js
function 'deleteticket' created
function 'searchtickets' created
function 'saveticket' created
function 'updateticket' created
function 'getticketsla' created
function 'saveticketsla' created
cmd: /usr/local/bin/fission route create --method DELETE --url /deleteticket --function deleteticket
cmd: /usr/local/bin/fission route create --method GET --url /searchtickets --function searchtickets
cmd: /usr/local/bin/fission route create --method POST --url /saveticket --function saveticket
cmd: /usr/local/bin/fission route create --method GET --url /getticketsla --function getticketsla
cmd: /usr/local/bin/fission route create --method PUT --url /updateticket --function updateticket
trigger 'be5917f6-4211-4ecc-be53-b81f437d0156' created
cmd: /usr/local/bin/fission route create --method POST --url /saveticketsla --function saveticketsla
trigger 'b7b76c73-c3ee-4806-b632-4e35150ce297' created
trigger 'f9d558d1-6e39-4896-9f49-b89db4d6ad99' created
trigger '91d30b42-361d-45f1-9508-2a51b5f539fd' created
trigger 'fbd9649e-bdb0-47a8-80f5-5c2ac1ea34e0' created
trigger '5b6cfa26-dfc8-4f39-86b8-63933754c798' created
[centos@ip-10-0-1-162 hackathon]$ curl http://$FISSION_ROUTER/searchtickets
Hello, world!
```