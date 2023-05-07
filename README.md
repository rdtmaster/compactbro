# CompactBro
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

- Download binary for your operating system
- Place `templates` and `static` in the same folder (TODO <- automate it with Github Actions)

### Install from source
```shell
git clone https://github.com/rdtmaster/compactbro
cd compactbro
go mod tidy
go build
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

(*replace values in `<...>` with relevant data, do not include `<` and `>`*)

### HTTPS
If you want compactbro to become a drop-in replacement for `i.reddit.com` you need to specify it in your hosts file, on windows it is located at `C:\Windows\System32\drivers\etc/`. Open `hosts` file with text editor of your choice and add the following:
```
127.0.0.1 i.reddit.com
```
Then save and exit. You can verify it: (press Win+r => type `cmd` => enter and type following command):
```
ping i.reddit.com
```
You should see `127.0.0.1` IP address; sometimes changes require reboot to become active.

However, the point of doing this is to revive links pointing to `i.reddit.com` so you will need to generate SSL certs.
```shell
cd <path-to-compactbro>/certs

openssl ecparam -genkey -name prime256v1 -out key.pem

openssl req -new -sha1 -key key.pem -out csr.csr -subj "/CN=i.reddit.com" -addext "subjectAltName=DNS:i.reddit.com,DNS:i.reddit.com,IP:127.0.0.1"

openssl req -x509 -sha1 -days 365 -key key.pem -in csr.csr -out certificate.pem
```
Then open your `compactbro.toml` config file. All HTTPS settings are stored in the respective section which should look like this (you can append it to the end of the file):
```toml
[HTTPS]
Use = true
LocalAddress = "127.0.0.1:443"
KeyPath = "/certs/key.pem"
CRTPath = "/certs/certificate.pem"
```

This will generate self-signed SSL certificate, which you have to import to your browser as a root CA (for firefox: `settings` => `advanced` => `certificates` => `import`). Most likely the browser still won't allow you to visit `https://i.reddit.com/*` pages due to transport security (because the certificate had changed and could not be verified). Find your firefox profile folder and look for `SiteSecurityServiceState.txt` file. Although this is considered insecure you can simply delete/rename it or find records for `i.reddit.com` and remove them.

You will need to "Add security exception" when visiting https://i.reddit.com for the first time. During these operations you will see a lot of security here-be-dragons-like warnings. In general it is safe to ignore them since you are dealing with certificates you just signed yourself, but note that you do everything at your own risk.