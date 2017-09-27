datadogbot
----------

Build Status: [![Build Status](https://travis-ci.org/packateloop/datadogbot.svg?branch=master)](https://travis-ci.org/packetloop/datadogbot)

[Binary Releases](https://github.com/packetloop/datadogbot/releases)

We rely on Datadog to trigger alert events however it
is not easy to view different alerts in a big picture.
Hence, we would collect these events from a event
stream and this would provide us information such
as frequent alerts for us to address them or use
these information to make platform smarter.

TBD:

1. Improve test coverage.

Usage:
======

Set Slack API token:

```bash
$ cp env.example .env
```

Development:
============

Run tests:

```bash
$ make test
$ make coverage
```

Below command would run unit tests and produce a binary

```bash
$ make install
```

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
