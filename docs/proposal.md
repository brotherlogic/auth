<!-----



Conversion time: 0.291 seconds.


Using this Markdown file:

1. Paste this output into your source file.
2. See the notes and action items below regarding this conversion run.
3. Check the rendered output (headings, lists, code blocks, tables) for proper
   formatting and use a linkchecker before you publish this page.

Conversion notes:

* Docs to Markdown version 1.0β37
* Wed Jul 17 2024 17:18:02 GMT-0700 (PDT)
* Source doc: Kubernetes Auth
----->



# Kubernetes Auth

<p style="text-align: right">
Brotherlogic</p>


<p style="text-align: right">
</p>


<p style="text-align: right">
Draft</p>



### Abstract

This document describes the process of doing Auth for Kubernetes services. This allows us to call kube services remotely but still securely. For internal dev machines or local cluster machines the auth tokens are passed automatically and it offers and end point to retrieve the token for remote machines (or local machines that are remote).


### Background

We need to have an easy ethod to allow us to make requests into services running on kubernetes without having to manage per-service passwords or tokens. The auth service will support this by supporting (1) token cycling, (2) token passing to known receivers and (3) retrieval of token using a password and (4) logging of events to support identification of security breaches


### Process

We add an auth server which has an internal address and an external one. Both support a GetToken API request but the internal one will ignore the necessary password to get the token, whereas the external one requires a password. Password is stored as a kubenetes secret and is rotated when necessary.

Every 24 hours auth generates a single token as a UUID with an additional suffix of a large random number. This is the token that is returned on the GetToken request. The within Kube service also supports a “Auth” method that validates that the passed token matches what is stored in Auth.

Dev machines verify their token every 5 minutes, on a failure to auth they reach out to get the new token and store it on the local filesystem. The token is attached to the rpc context and extracted as necessary by backend services and authenticated.

Key thing here is that the token is only ever stored in memory - we never persist tokens.


### Milestones



1. This doc is added to the repo
2. Add proto to support auth process
3. Builds out the server to answer API calls.
4. Writes the token cycling system
5. Adds a module to support automatic authentication
6. Adds module to seperate binary to handle authentication through a go interceptor
7. Interceptor caches auth result and updates cache on the basis of the received result
8. Add CLI to support retrieval of auth token from dev machines
9. Automatically install CLI onto all dev machines
10. Crontab method on dev machines to regularly update token










<!-- watermark --><div style="background-color:#FFFFFF"><p style="color:#FFFFFF">gd2md-html: xyzzy Wed Jul 17 2024</p></div>

