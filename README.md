# github_traffic_exporter
Export Github repo traffic data to Prometheus

WORK IN PROGRSS - currently is collecting data from GitHub and exporting it through the /metrics endpoint - NOT READY FOR USE!

## What this is

This is an exporter of data from Github to Prometheus.  It is written in go (mostly because I wanted to renew my go skills).  

This code uses Google's [go-github library](https://github.com/google/go-github/) which in turn uses the GitHub [traffic](https://docs.github.com/en/rest/metrics/traffic) and [repository](https://docs.github.com/en/rest/repos/repos) APIs.

## Wait - isn't there already a github_exporter?

Yes.  Infinityworks wrote their [Prometheus HitHub Exporter](https://github.com/infinityworks/github-exporter) but it does not expose data on views or clones, which is what I wanted the most.  Their code also seemed pretty complicated, but probably because I have been away from go for awhile.  I wanted more than just to refactor an addition into their code.  I wanted to build something myself.  

I might find that after building this that the right thing is to backport these new additions to their exporter, but time will tell.

## What problem are you trying to solve?

I want to measure the "awareness" of code repositories.  I want to see if people are looking at the repo, if they are cloning it, forking it, opening issues against it, watching it.  Increases in those will tell me that "awareness" of that code is rising.  I cannot tell if folks are actually building something useful and deploying it, but I can at least see if folks are actually kicking the tires.

## Why Prometheus

Prometheus is an excellent and proven metrics collector and time-series database store designed to support easy queries and aggregation.  It also has the ability to generate Alerts through it's Alert Manager.  For what I am measuring it's a great fit and I don't need to build anything but the shim layer to get the data into it from GitHub.

## Prometheus Configuration

TBD

## Other useful github tooling

TBD