# CompactBro
**Warning**: This software is yet to reach v1. Some functions are simply unimplemented; development started in spring of 2023, after reddit removed `/.compact` functionality, however once the use of `.i` and a relevant userscript has been discovered, this project has been immediately obendened (v0.83.0) and stood there until the end of 2024. In spring of 2024 the `.i` has been shut down by reddit as well, but since the author pretty much left reddit at that time it was barely noticed and no effort has been taken to develop the software and move it towards v1. Now I decided to release it as is, allowing to browse and make comments/replies.

Miss the old reddit mobile web/compact/wap interface (`i.reddit.com`, `.compact`)? Worry not, CompactBro brings it back!

1. Download & unzip a version for your platform -> https://github.com/rdtmaster/compactbro/releases/latest
2. Rename `compactbro.sample.toml` -> `compactbro.toml` and copy to your default config directory
3. From your reddit account, Create a new app with script access, then specify your `secret`, `id`, reddit login and password in `compactbro.toml`. Link to do it -> https://old.reddit.com/prefs/apps
4. Launch the app
5. http://localhost/ <- in your browser address bar

## what it is?
A reddit client for web that focuses on fast load time, old browser compatability and delivering user experience similar or at least reasonably close to that provided by i.reddit.com mobile interface, also known as `.compact`.
 
## Why?
In late march 2023 the mobile web interface got removed by reddit staff.

### but alternatives (teddit and the like) exist
None of them are like the old ``.compact`` interface. Most lack old browser compatability, some are 'privacy-oriented' which means you cannot log into your reddit account.

## What it's not?
It's not a privacy-focused client or a reddit automation software.

## Is pasting my login/password safe
Short answer: yes. CompactBro is fully open-source and contains no backdoors. It is generally safe, but I am not responsible for 3rd party libraries and also cannot guarantee your account won't die because of API use (even though the API is public, most websites do not like when you use anything but the official web/mobile client). If your account poses great value I wouldn't advice to risk it.

## How to install
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
EcoMode = true
MarkMsgsUnreadOnView = true
CheckMsgs = true
Logging = true
DisplayFlairEmojis = true
NightMode = false
DefaultLimit = 25
LocalAddress = "127.0.0.1:80"

[HTTPS]
Use = true
LocalAddress = "127.0.0.1:443"
KeyPath = "/certs/key.pem"
CRTPath = "/certs/certificate.pem"


[Auth]
Use = false
Username = "user"
Password = "pass"

[Credentials]
ID = "<ID>"
Secret = "<Secret>"
Username = "<user>"
Password = "<pass>"

[TemplateOptions]
PrettyPrint = true
LineNumbers = true
```
replace values in `<...>` with relevant data, do not include `<` and `>`


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
If you require traffic between compactbro instance and your browser to be encrypted there's an option to use HTTPS. Usually you don't need it if you run the program locally.

To generate keys, install OpenSSL (there are binaries available for most platforms, refer to their site for more info) and generate SSL certs:
```shell
cd <path-to-compactbro>/certs
openssl ecparam -genkey -name prime256v1 -out key.pem
openssl req -new -sha1 -key key.pem -out csr.csr -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost,DNS:localhost,IP:127.0.0.1"
openssl req -x509 -sha1 -days 3650 -key key.pem -in csr.csr -out certificate.pem
```
This will generate self-signed SSL certificate, which you have to import to your browser as a trusted root CA (for firefox: `settings` => `advanced` => `certificates` => `import`). 

Then open your `compactbro.toml` config file. All HTTPS settings are stored in the respective section which should look like this (you can append it to the end of the file):
```toml
[HTTPS]
Use = true
LocalAddress = "127.0.0.1:443"
KeyPath = "/certs/key.pem"
CRTPath = "/certs/certificate.pem"
```

You will need to "Add security exception" when visiting https://localhost/ for the first time. During these operations you will see a lot of security here-be-dragons-like warnings. In general it is safe to ignore them since you are dealing with certificates you just signed yourself, but note that you do everything at your own risk.
You can use your own domain as well.

## V1 roadmap
- [x] Front page
- [x] View a subreddit
- [x] View Posts+comments
- [x] Edit posts
- [x] Edit comments
- [x] Comment under post and reply to comments
- [x] Infinite scrolling subreddit, DMs and user overview page
- [x] Check inbox and display orange icon, add config option to turn it on/off
- [x] View messages (DMs, comment replies, post replies, username mentions, sent)
- [x] Auto mark viewed as read (option to disable this in config)
- [x] User overview
- [x] Coloreful post/comment flairs
- [x] User attrs (submitter, mod, admin)
- [x] Image/video thumbnails
- [x] Mark NSFW, spoiler
- [ ] Handle more comments (partially complete)
- [ ] GIFs display
- [ ] Delete comment
- [ ] Delete post

## Roadmap for future versions
- [ ] Full support for DM messages
- [ ] Submit new post
- [ ] Create sub
- [ ] Open video links
- [ ] Moderation and modmail
- [ ] Caching
- [ ] night mode (placeholder and option created)
- [ ] Use better template engine


## Versioning policy
It'll be nice if v1 can ever see the light and even nicer if something beyond that gets released.
New versions are released on every change pushed to the repo, in general you should download the latest version.

## Compatability and intentional decisions
1. Javascript (e.g. frontend) has been completely re-written without use of `jQuery` library. Adopting reddit's codebase would be even harder then coding everything from scratch. Initial plan was to leave HTML and CSS untouched, but it was impossible
2. Compactbro web interface supports Firefox 47.0 and server supports Windows 7x32. Effort should be made to try and support older versions of Firefox. At no point should Compactbro ever drop Windows 7 (incl.32 bit) support or bump minimum supported browser version. If adding any new feature requires to raise those minimum system requirements, that feature should be forgotten. Note: the last version of Go officially supporting Win7 is `go1.20.14` so it must be used to compile CompactBro.
3. Compactbro assumes you are over 18 and are willing to view NSFW and spoiler content. By using this software you confirm you reached that age and have no issue with adult materials.
4. CompactBro uses the `amber` engine. It turned out to be pretty backward and limited in terms of functionality, most notably it has no support for recursive mixins. That said, I enjoyed using this engine, it makes creating templates very fast, probably as rapid as it could be. Making page templates still eight around 70% of the time spent to develop this software, with any other engine I doubt I can ever finish it.