# Reddit-Bot

## Introduction

This is a quick and dry implementation of reddit-bot in Go programming language.

## Features to be implemented

**Account Creation**
- Can create verified and unverified reddit accounts
- Has Kopeechka.store enabled so you dont have to worry about providing emails
- Can use custom Passwords, usernames or create with a catchall email address
- Can set account as NSFW, enable NSFW feed Settings, Disable Safe Browsing and Unsubscribe from receiving Reddit Emails

**Upvote/Downvote**
- Can upvote posts from post link,twitter link or can browse through “new” or “hot” sections to locate post and upvote it
- Can upvote other posts randomly before and after main upvote to obfuscate the upvote target
- Can upvote latest post on profile link, skipping pinned posts
- Can follow the profile after upvoting
- Can give award to profile after upvoting (if available)
- Can upvote all comments in a profile post link
- Can comment on profile post after upvote

**Post Comments**
- Can post comments from a list
- Can browse through “hot” or “new” or via a twitter link or direct reddit link to make the comment

**Comment and Upvote**
- Can make comment and upvote post
- Can comment and Downvote post

**Upvote/Downvote Comments**
- Can upvote/downvote comments
- Can upvote/downvote all comments

**Posting (Karma Farming)**
- Can make posts to list of subreddits on a timed delay
- Posts can be from single or multiple accounts
- Can post top comment after post has been made, with another account
- Can Autoupvote new post or comment made
- Can crosspost from a list of post links
- Can crosspost with single or multiple accounts
- Can change account passwords
- Can subscribe accounts to a subreddit
- Can mass report a post or comment

**Check Accounts**
- Can confirm if accounts are alive or dead (shadowbanned)
- Can login accounts and save cookies so it does not login again when performing other actions
- Can toggle NSFW and NSFW Feed settings when logging in accounts
- Can load accounts to manual mode
- Can load accounts and randomize avatar andclaim award if available
- Can assign new proxies to accounts
- Standby mode can immediately upvote any post/comment link it finds in specific file which it scans periodically (suitable for mini SMM)
- Warm Accounts
- Can warm up accounts by visiting random posts from reddit homepage
- Can randomly upvote any post according to preset upvote probability
- Can set average time to spend with each post
- Can join reddit website suggested subreddits
- Can comment on random posts and randomize avatar

**Settings**
- Is multithreaded
- can elect to use warmed up accounts or not for actions
- Can drip feed its actions to spread them over stated time.
- Is enabled to solve captchas with 2captcha.com
- Is enabled to use lates Anti-fingerprinting profiles
- Is enabled to use http or socks5 proxies
- Can disable images to save on bandwidth

**Scheduler**
- has an inbuilt scheduler with which you can schedule the bot to do any of its functions at specific intervals or on specific dates and at specific times

**Upcoming feature**
- Can upvote a post till it reaches a predetermined position on the “hot” section of the subreddit and then stop upvotes
- Can resume upvotes if the post position drops from the predetermined position.