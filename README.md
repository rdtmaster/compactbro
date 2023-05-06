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
- Login to reddit and head over to **preferences** => **apps**, create a new APP. The name could be any (e.g. `compactbro`), type should be `script` and other fields are optional. Note down App ID and secret.
- create config file `compactbro.toml` in your `<default-config-directory>\compactbro` (if you are unsure of the location, launch the software and it will print the path in uss, create the file there)
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

## Configuration
Work in progress
