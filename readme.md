# ER-Network-Calculator

### Hello, World!
`ER-Network-Calculator` is a project developed and maintained by `@Pheric` and `@Naluca` on Github. The purpose of this
web application is to provide a set of tools to make working with IP addressing easier. It is built around what is taught
in Networking I & II at Dakota State University, and it is based off of
[Dan Dirks' Networking Calculator](https://singleton.dsu.edu/subnet-calculator/) which is maintained by Dakota State
University's ITS team. While this calculator is based off of that one, it shares no code (`@Pheric`: I specifically
avoided looking into it). In addition to a new (less impressive) interface, this new version also has more support
for IPv6 addressing (using a custom library). More features are in the works.

### Why make a new version?
While Dan Dirks' version is excellent, we wanted to see if it could be improved. In its current state, it cannot
subnet IPv6 addresses, and simple tasks like finding a network address must be done manually.

*This calculator is not intended to be a replacement.* Dan's calculator has a number of benefits over this one:
1. All computations are done client-side
2. The interface is easy to use, and lets you see a history of past operations
3. It has generally useful tools, like being able to convert numbers to binary and back, and performing the binary AND
function on two or more numbers
4. Best of all: it shows _how_ to do everything

Finally, we wanted the challenge. Hats off to you, Mr. Dirks, for designing such a useful tool that has impacted hundreds of
students over the years, and will continue to do so. We hope that this project will be useful to someone as well,
and even if it isn't, it has been a great learning experience.

### How do I run this?
In the first commits in this repository, the binary was generated and pushed to GitHub under the name `www.elf`. As long
as the other files are in the same directory and are readable, there should be no problems. This project was designed on
Linux, for Linux, so your mileage may vary.

`git clone https://github.com/Pheric/ER-Network-Calculator.git && cd ER-Network-Calculator && chmod +x www.elf && nohup sudo ./www.elf --port 80`

There is an online version, but it is currently restricted to those with the link until bugs have been worked out and
more permanent hosting is decided.