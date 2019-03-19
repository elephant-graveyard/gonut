# Simple Cloud Foundry demo App

## Goal
This is a very basic swift demo app for cloud foundry, called **johnny 5**.

## Pre
Access to a Cloud Foundry with a swift buildpack is pre installed.

## Run this app

```sh
git clone https://github.com/idev4u/johnny-5.git
cd johnny-5
cf push my-johnny-5
```

## Advanced
Swift is less consuming memory language, to move this advantage to the cloud you should push the app with this command:
```sh
cf push my-johnny-5 -m 32M
```

## How to use

```sh
curl 'http://my-johnny-5.mybluemix.net/'
Hello. I'm Johnny 5.
```
