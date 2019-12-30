--- title
Chaos Communication Congress Day 1
--- description
First!!!1!!l!!
--- main

### Day _1_

Like a deer blinded by headlights I was.

Anyways, setup phone wifi with dedicated app,
laptop with wpa_supplicant settings,
see [wiki](https://events.ccc.de/congress/2019/wiki/index.php/Static:Network/802.1X_client_settings).

Also setup wireguard VPN to always on,
worked _perfectly_.
Total usage over 4 days:
phone tx 2.96GiB / _rx_ 3.73GiB,
laptop tx 0.11GiB / _rx_ 1.61GiB

Stll don't know why it _didn't_ work when I tried crafting it myself

```
 network={
 	ssid="36C3"
 	key_mgmt=WPA-EAP
 	eap=TTLS
 	identity="edward"
 	password="snowden"
 	# ca path on debian 7.x, modify accordingly
 	ca_cert="/etc/ssl/certs/DST_Root_CA_X3.pem"
 	altsubject_match="DNS:radius.c3noc.net"
 	phase2="auth=PAP"
 }
```

#### Talks

things i heard,
may or ay not remember actually hearing it

##### Open Source is Insufficient to Solve Trust Problems in Hardware

_Like duh_
Who knows of a good hashing algorithm for hardware?
Is this the same as the universe's simulation theory?
**If as simulations we can tell we are being simulated
then amaybe software can tell if hardware is running properly**

##### What's left for private messaging?

_Uh-huh_
still no good messaging app.
We get more apps,
they're not obviously better,
refer to [xkcd 927 Standards](https://xkcd.com/927/), [xkcd 1810 Chat Systems](https://xkcd.com/1810/)

##### phyphox: Using smartphone sensors for physics experiments

_Oooh_
Admittedly didn't expect too much out of this one,
but entertaining nevertheless.
Maybe I've seen parts of this on YouTube?

##### How to Break PDFs

_Urgh..._
5 mins in after they explained a PDF file structure,
yeah, broken in all the expected places.
At least they followed through on testing in all 22 apps.

##### Plundervolt: Flipping Bits from Software without Rowhammer

_Yo-ho-ho_
Who doesn't like pirates?
Very fun style of talk,
probably accurately reflects software development flows.

##### (Post-Quantum) Isogeny Cryptography

_Zzzzzz_
did i really just listen to 30 mins of RSA thingies,
then handwavy generalizations of this elliptical curve thing?
Who is the target audience?
someone who doesn't understand crypto
(and won't understand the rest)
or someone who will understand the maths
(and wasted half an hour listening to crypto 101).
Still don't know how this compares to other post-quantum things

##### Practical Cache Attacks from the Network and Bad Cat Puns

_Ehhhh_
Maybe slightly overhyped?
Interesting chain of optimizations that allow you to observe caches,
but you can only see that something is there,
not what.
Back to stats and AI to guess what you wrote.

##### SELECT code_execution FROM \* USING SQLite;

_Wow_
seriously impressive thought process.
`SELECT` needs info from a master table which can be rewritten
to execute arbitrary code (with a weird self registered).
Also **:faceplam:**.
Maybe running a virtual machine to do SQL queries isn't the best of ideas?
