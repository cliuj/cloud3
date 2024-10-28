# 6 WEEK DESIGN DOC - cloud3
* Project: `cloud3`
* Date: 10/27/24
* Authors: Charlie Liu

# Introduction
The purpose of this document is to provide an overview of what needs to be done to prepare and expand
the proof-of-concept version of `cloud3` for production (customer-facing) in a 6 week timeline.

## Architecture (proof-of-concept version)

#### About
```
   +---------- upload ------- + -------- upload --------+
   |                          |                         |
   v                          |                         v
[client 1] --- upload ---> [server] <--- upload --- [client 2]
   |       (compare SHASUM)   |    (compare SHASUM)     |
   |                          |                         |
get SHASUM                get SHASUM                 get SHASUM
(polls)                      (polls)                  (polls)
   |                          |                         |
   |                          |                         |
[shared dir]             [storage dir]             [shared dir]
```

As shown in the chart above, both the `server` and the `clients` poll for changes in their
respective directory by generating the checksum (SHA256) of the files.

The clients poll the `server` every 1 second for any changes on the `server`'s checksum,
if there's a difference, the `client` that has the difference uploads it's local contents of the
shared directory onto the `server`. When the `server` receives the uploaded files, it uploads it
to the rest of the `clients` (*1 ideally to the one it recieved the uploads from).

The flow is as follows:
1. clients, servers are up
2. user makes a change in the shared directory via adding or editing files
3. The user client detects the changes in the directory, and polls the remote's checksum
4. The user client checks if the remote and local checksums match
5. If `match`: Do nothing
6. If `not match`: The user client's state of the shared directory gets uploaded to the server
7. Server uploads its state of its shared directory to all the other clients via client list
8. Repeat


(*1, There's an issue with the code that is suppose to enforce this, but the comparison isn't
working ATM, thus all clients and servers upload their files to each other).

#### Limitations
The current proof-of-concept program, while does seem to work, isn't scalable nor ideal. It contains multiple limitations:
* Syncing between `clients` and `servers` results in redundant upload calls back to the original client
  - Can be considered OK if it wasn't uploading everything
* Client list is hardcoded.
  - Tried implementing a registration system where clients connected to the server are added to a client list
    for the server to upload back to
    - Couldn't get it working due to some issue retrieving the request Host URL
* System **does not** support any deletion of files
  - This also means it doesn't support diffs and renames either
* System **does not** support nested directories
* System **does not** have any better means of storing files (files are stored raw) thus, high disk usage
* To be added...

## 6 Weeks Improvements
Assuming a team was given 6 weeks to expand upon and prepare the program for production and user testing.
There are a couple of priority items needed before release:
1. **User Accounts/Security** - This is the most important one as users should expect that the files they upload belong to them and are private to their accounts.
    Currently, no account system is in-place and thus user uploaded files are public and unsecured.
    * Goals:
      * Setup User accounts
      * Setup User authentication - only users that own the account should have access to it
      * Paritioning of user storage - do not want users to share storage
      * Check if files are malicious

2. **Cleanup of PoC limitations** - The proof-of-concept has multiple limitations listed above that should be addressed before release.
    The top priority ones are the following:
    * Goals
      * Setup registration of clients to servers - remove the hardcoded client list and create a registration system where the
      server knows what clients to broadcast to. Highest priority as this may even block account creation.
      * Review architecture again - Current implementation needs to be expanded upon
        * Ideally, server should broadcast a message to online clients that they may pull changes rather than the current sitation where the server pushes to the clients.
        (probably not secure either)
        * Look to limit polling and uploading - if a client is malicious and somehow manages to modify the polling times, it could potentionally turn into a DDOS attack.
      * Sync only the changes of files. i.e pull or upload the diffs and not the entire directory

3. **Scaling** - Assuming the above are addressed, the next step is to scale the application. Currently, the expectation is that there's only 1 user connected at a time to a client.
    As the number of users and clients increase, the system needs to be more robust in handling multiple connections and concurrency issues
    * Some potential issues:
      * Race conditions:
        - 1 user making changes on 2 or more clients -> maybe resolved with first come first serve messaging queue
      * Storage:
        - Files are stored raw thus take a lot of space, should we aim to compress them?
      * Performance:
        - Is REST the best way?
        - Maybe use another protocol (gRPC possible?) to transport files over the network
