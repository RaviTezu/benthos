---
title: beanstalkd
type: output
status: experimental
categories: ["Services"]
---

<!--
     THIS FILE IS AUTOGENERATED!

     To make changes please edit the corresponding source file under internal/impl/<provider>.
-->

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

:::caution EXPERIMENTAL
This component is experimental and therefore subject to change or removal outside of major version releases.
:::
Write messages to a Beanstalkd queue.

Introduced in version 4.7.0.

```yml
# Config fields, showing default values
output:
  label: ""
  beanstalkd:
    address: ""
    max_in_flight: 64
```

## Fields

### `address`

An address to connect to.


Type: `string`  

```yml
# Examples

address: 127.0.0.1:11300
```

### `max_in_flight`

The maximum number of messages to have in flight at a given time. Increase to improve throughput.


Type: `int`  
Default: `64`  


