# CompactBro
Miss the old reddit mobile web/compact/wap interface (`i.reddit.com`, `.compact`)? Worry not, CompactBro brings it back!

1. Download & unzip a version for your platform -> https://github.com/rdtmaster/compactbro/releases/latest
2. Rename `compactbro.sample.toml` -> `compactbro.toml`
3. From your reddit account, Create a new app with script access, then specify your `secret`, `id`, reddit login and password in `compactbro.toml`. Link to do it -> https://old.reddit.com/prefs/apps
4. Launch the app
5. http://localhost/ <- in your browser address bar

More (much more!) details below

## what it is?
A reddit client for web that focuses on fast load time, old browser compatability and delivering user experience similar or at least reasonably close to that provided by i.reddit.com mobile interface, also known as `.compact`.
 
## Why?
In late march 2023 the mobile web interface got removed by reddit staff.

### but alternatives (teddit and the like) exist
None of them are like the old ``.compact`` interface. Most lack old browser compatability, some are 'privacy-oriented' which means you cannot log into your reddit account.

## What it's not?
It's not a privacy-focused client. If it provides any additional privacy/security (or so you believe) then congrats, go celebrate & rejoice, however this is unintentional and not something to rely on. And neither is it a full scale reddit client utilizing all features exposed by the official API. It is and likely will remain only a subset of what the official front-end provides, much like the `i.reddit.com` used to be (it had not the edit post feature, posting was limited, no moderation etc). If we can achieve same feature coverage this will exceed my expectations by far. See the roadmap below for detail.

### Why won't you add dynamic user account switching
Because I believe this could open the door for multi accounting, automation and other stuff reddit might get upset about. Yes reddit allows alts, you can use them as well by editing configs. This decision is final as of right now, I am not taking it because I am all that invested into reddit cleanliness or following ToS, I just don't want to contribute to any shady activity, which will ultimately lead to reddit becoming more strict. If you want an automated solution for upvoting, spamming, pharming accounts, go develop one or buy it off the market, plenty of them available. I want reddit to stay as relaxed as it is.

### Why not an in-browser client?
I thought about it but there are several obstacles. One of the most ambitious goals I have so far is making this app a drop-in replacement of `i.reddit.com` through the use of `.hosts` file. I have a freaking number of bookmarks pointing to `i.reddit.com` pages, would be great to revive them.

## How to install
Work in progress!
**WARNING the software is not functional yet do not install it!**

First, download binary for your operating system from the releases page.  Alternatively, clone this repo and build it (you need `git` and functional go >= 1.20.3 installation):
```bash
git clone https://github.com/rdtmaster/compactbro
#(wait...)
cd compactbro
go mod tidy
go build
chmod +x compactbro && ./compactbro #on *nix
compactbro.exe #on windows
```
## Configuration
- Login to reddit and head over to **preferences** => **apps**, create a new APP. The name could be any (e.g. `compactbro`), type should be `script` and other fields are optional. Note down App ID and secret.
- create config file `compactbro.toml` in your `<default-config-directory>\compactbro` (if you are unsure of the location, launch the software and it will print the path in use, create the file there)
- Edit the config, it will look like this:

```toml
EcoMode = false
[Credentials]
ID = "<ID>"
Secret = "<Secret>"
Username = "<My-reddit-username>"
Password = "<My-reddit-password>"

[TemplateOptions]
PrettyPrint = true
LineNumbers = true
```
replace values in `<...>` with relevant data, do not include `<` and `>`

Before describing the config process, let's summarize it to help you decide what you need and don't need depending on your goals.

1. **Running `compactbro` locally, without replacing `i.reddit.com`**. No HTTPS and no Authentication
2. **Running `compactbro` locally, replacing `i.reddit.com`**. No authentication, HTTPS enabled
3. **Running `compactbro` on a remote server without replacing `i.reddit.com`**. Auth and optional (but recommended) HTTPS
4. **Running `compactbro` on a remote server, replacing `i.reddit.com`**. Authentication+HTTPS.

Read further to know what do these steps mean and, when you decide what kind of setup you want, this list will help you realize which steps are required.

### Auth
Generally there are two strategies of using compactbro: run it locally or on a server. If the former, you don't need any authentication whatsoever, simply fire it up; but running it on a remote server can provide several advantages, namely you don't need to have any additional process launched on your device. A low-end Linux VPS server should suffice and, once you configure Authentication and HTTPS, you can make it a drop-in replacement for `i.reddit.com`.
You should take into account that without authentication compactbro exposes its interface on `IP:port` specified in the config file. Be careful when choosing the listen address. If you specify `127.0.0.1` it is available only within the realm of your machine; if it is a local address like `192.168.5.5` (or another range used in your LAN) it can be accessed by other users in the same network. And if you choose your WAN IP or use `0.0.0.0`  it becomes available to the whole world (provided you have a direct none-NAT internet connection). The latter option is meant for running the app server-side.
TLDR: if your app listens on an address available to other people, set up authentication.

Authentication-related settings live in the `[Auth]` section of your `compactbro.toml` file, if you want to enable authentication it should look like this:
```toml
[Auth]
Use = true
Username = "<username>"
Password = "<password>"
```
Replace values in `<...>` with credentials of your choice; be advised these are basic auth credentials you will use to access the app, they have nothing to do with username/password you use to log into your reddit account, therefore these do not have to match, even though they can if so you choose.

### HTTPS
Today majority of websites and apps are overly secure, with HTTPS enforced everywhere. Most resources in fact do not need it. There's nothing wrong about using plain HTTP in most scenarios. Any half sain person understands that if you just want to read the news or even sh1tpost on reddit it should not involve cryptography (which introduces great overhead and a big number of errors). Well, at least not by default. Think twice prior to enabling HTTPS, chances are you don't need it.
However, if you transmit any data without HTTPS **in theory** it could be intercepted. So in a nutshell, if you use compactbro server-side and you live in a free/democratic country, there's not much to worry about. Yet in case you are concerned about government surveillance or are just a little bit paranoid, there's an option to use HTTPS built into Compactbro. It could come handy if you want compactbro to become a drop-in replacement for `i.reddit.com`. Below are the steps to achieve this. The manual is (work-in-progress), it covers only Windows+Firefox setup but it should work basically anywhere with little changes.

First of all, you need to specify make changes in your hosts file, on windows it is located at `C:\Windows\System32\drivers\etc/`. Open `hosts` file with text editor of your choice and add the following:
```
127.0.0.1 i.reddit.com
```
if you use compactbro on a remote server change `127.0.0.1` to its IP, below we use `127.0.0.1` as a placeholder, don't forget to edit it if needed.
Then save and exit. You can verify it: (press Win+r => type `cmd` => enter and type following command):
```
ping i.reddit.com
```
You should see `127.0.0.1` IP address; sometimes changes require reboot to become active.

The HTTPS however isn't operational yet, since any SSL connection requires a certificate. Install OpenSSL (there are binaries available for most platforms, refer to their site for more info) and generate SSL certs:
```shell
cd <path-to-compactbro>/certs

openssl ecparam -genkey -name prime256v1 -out key.pem

openssl req -new -sha1 -key key.pem -out csr.csr -subj "/CN=i.reddit.com" -addext "subjectAltName=DNS:i.reddit.com,DNS:i.reddit.com,IP:127.0.0.1"

openssl req -x509 -sha1 -days 365 -key key.pem -in csr.csr -out certificate.pem
```
This will generate self-signed SSL certificate, which you have to import to your browser as a trusted root CA (for firefox: `settings` => `advanced` => `certificates` => `import`). Most likely the browser still won't allow you to visit `https://i.reddit.com/*` pages due to transport security (because the certificate had changed and could not be verified). Find your firefox profile folder and look for `SiteSecurityServiceState.txt` file. Although this is considered insecure you can simply delete/rename it or find records for `i.reddit.com` and remove them.

Then open your `compactbro.toml` config file. All HTTPS settings are stored in the respective section which should look like this (you can append it to the end of the file):
```toml
[HTTPS]
Use = true
LocalAddress = "127.0.0.1:443"
KeyPath = "/certs/key.pem"
CRTPath = "/certs/certificate.pem"
```

You will need to "Add security exception" when visiting https://i.reddit.com for the first time. During these operations you will see a lot of security here-be-dragons-like warnings. In general it is safe to ignore them since you are dealing with certificates you just signed yourself, but note that you do everything at your own risk.

## V1 roadmap
- [x] View subs
- [x] View Posts+comments
- [x] Edit posts
- [x] Edit comments
- [ ] Handle more comments
- [x] Comment under post and reply to comments
- [ ] Delete comment
- [ ] Delete post
- [ ] Infinite scrolling subreddit and user overview page (1/2 complete)
- [x] Check inbox and display orange icon
- [ ] View messages
- [x] User overview
- [x] Coloreful post/comment flairs
- [x] User attrs (submitter, mod, admin)
- [x] Image/video thumbnails
- [x] Mark NSFW, spoiler
- [ ] Create good readme (partially complete)

## Roadmap for future versions
- [ ] Full support for DM messages
- [ ] Submit new post
- [ ] Create sub
- [ ] Open video links
- [ ] Moderation
- [ ] Caching
- [ ] night mode
- [ ] Use better template engine


## Versioning policy
It'll be nice if v1 can ever see the light and even nicer if something beyond that gets released.
Minor versions are subject for eventual delition from the releases section, except the latest one. At the moment v1 is the only major version planned, its goal in a nutshell: "*release something that compiles and lets you browse and comment*". New major version should be released once something notable happens, minor versions are released whenever I want to see if it still works.

## Compatability and intentional decisions
1. Javascript (e.g. frontend) has been completely re-written without use of `jQuery` library. Adopting reddit's codebase would be even harder then coding everything from scratch. Initial plan was to leave HTML and CSS untouched, but it was impossible
2. Compactbro web interface supports Firefox 47.0 and server supports Windows 7x32. Effort should be made to try and support older versions of Firefox. At no point should Compactbro ever drop Windows 7 (incl.32 bit) support or bump minimum supported browser version. If adding any new feature requires to raise those minimum system requirements, that feature should be forgotten.
3. Compactbro assumes you are over 18 and are willing to view NSFW and spoiler content. By using this software you confirm you reached that age and have no issue with adult materials.
4. Compactbro uses the `amber` engine. It turned out to be pretty backward and limited in terms of functionality, most notably it has no support for recursive mixins. That said, I enjoyed using this engine, it makes creating templates very fast, probably as rapid as it could be. Making page templates still eight around 70% of the time spent to develop this software, with any other engine I doubt I can ever finish it.