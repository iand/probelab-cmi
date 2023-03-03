---
title: Composition
weight: 3
---
# IPFS Network Composition

{{< hint info >}}
*Notes on what could appear on this page*

An interior detail page that provides details on the numbers and types of peers in the network. The kinds of questions we want to answer with this section are: how has the network changed over time? how diverse is the software participating in the network? what old versions of software need to be supported still? where are participants located?

- Number of peers over time
- Distribution of software agents (by peer count, current and historic)
- Newly discovered agents
- Kubo versions in use
- Geolocation of providers and clients
- Network providers (datacenter/private infra)
{{< /hint >}}


## Total Peer IDs Discovered Classification

![Peer count by classification](../plots/peer-classifications.png)

In the specified time interval from `2023-02-20` to `2023-02-27` we visited `` unique peer IDs.
All peer IDs fall into one of the following classifications:

| Classification | Description |
| --- | --- |
| `offline` | A peer that was never seen online during the measurement period (always offline) but found in the DHT |
| `dangling` | A peer that was seen going offline and online multiple times during the measurement period |
| `oneoff` | A peer that was seen coming online and then going offline **only once** during the measurement period |
| `online` | A peer that was not seen offline at all during the measurement period (always online) |
| `left` | A peer that was online at the beginning of the measurement period, did go offline and didn't come back online |
| `entered` | A peer that was offline at the beginning of the measurement period but appeared within and didn't go offline since then |

## Agent Version Analysis

### Overall

![Overall Agent Distribution](../plots/agents-overall.png)

Includes all peers that the crawler was able to connect to at least once: `dangling`, `online`, `oneoff`, `entered`. Hence, the total number of peers is lower as the graph excludes `offline` and `left` peers (see [classification](#peer-classification)).


### Classification

![Agents by Classification](../plots/agents-classification.png)

The classifications are documented [here](#peer-classification).
`storm*` are `go-ipfs/0.8.0/48f94e2` peers that support at least one [storm specific protocol](#storm-specific-protocols).

### Agents

![Crawl Properties By Agent](../plots/crawl-properties.png)

Only the top 10 kubo versions appear in the right graph (due to lack of colors) based on the average count in the time interval. The `0.8.x` versions **do not** contain disguised storm peers.

`storm*` are `go-ipfs/0.8.0/48f94e2` peers that support at least one [storm specific protocol](#storm-specific-protocols).


## Geolocation

### Unique IP Addresses

![Unique IP addresses](../plots/geo-unique-ip.png)

This graph shows all IP addresses that we found from `2023-02-20` to `2023-02-27` in the DHT and their geolocation distribution by country.

### Classification

![Peer Geolocation By Classification](../plots/geo-peer-classification.png)

The classifications are documented [here](#peer-classification). 
The number in parentheses in the graph titles show the number of unique peer IDs that went into the specific subgraph.

### Agents

![Peer Geolocation By Agent](../plots/geo-peer-agents.png)

`storm*` are `go-ipfs/0.8.0/48f94e2` peers that support at least one [storm specific protocol](#storm-specific-protocols).

## Datacenters

### Overall

![Overall Datacenter Distribution](../plots/cloud-overall.png)

This graph shows all IP addresses that we found from `2023-02-20` to `2023-02-27` in the DHT and their datacenter association.

### Classification

![Datacenter Distribution By Classification](../plots/cloud-classification.png)

The classifications are documented [here](#peer-classification). Note that the x-axes are different.

### Agents

![Datacenter Distribution By Agent](../plots/cloud-agents.png)

The number in parentheses in the graph titles show the number of unique peer IDs that went into the specific subgraph.

`storm*` are `go-ipfs/0.8.0/48f94e2` peers that support at least one [storm specific protocol](#storm-specific-protocols).
