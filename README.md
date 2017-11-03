<h2>Application Platform for trading of valuable medicinal herbs using blockchain</h2>

To overcome the smuggling of the valuable medicinal herbs and to save them from going extinct ,we have built an application that is using
the promising blockchain networks to enable the trust and transparency between producer,multiple parties or the consumer.The underlying 
blockchain network is built using hyperledger fabric network.It provides a web interface to register the produce or trade the production.


<h4> The trade details are secured using the private channels </h4>
Unlike the public blockchain network where transparency of the business transactions can be an issue for the enterprise,the network allows
private channels to be created to do the transaction.


<h5>To get started ...</h5>
Make sure you have the docker,node,golang installed in your setup.
<strong>Start by starting the fabric network that should spin up various containers needed to run the hyperledger network</strong>

<code>cd hyperledger-herb-app/herb-src/fabric-material/herb-app/</code>

<code>./startFabric.sh</code>

if everthing goes fine u should see the network spin up and output the initially stored herb details.

Next,
<code>npm i</code>
<code>node registerAdmin</code>
<code>node registerUser</code>
<code>node server</code>

cheers! If no error

head over to over browser to point to localhost:3001.

<strong>WIP</strong>
