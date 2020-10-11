---
description: spotify codes
title: spotify codes
---

### _spotify_ code

Ever wondered what that wavy thing under shareable spotify images were?
They're [scannable codes][web]
the scanner is the camera icon in the search bar on mobile.
Apparently launched in 2017,
but who knows, I never discovered the feature.

Recently they appeared under _Start a group session_ in the connect menu,
(the thing you click to play music on another device).

![spotify code][sc]

23 bars, I'm guessing bars 1, 12, 23 (first, middle, last)
are for calibration, that leaves 20 bars.
Each bar has a top and bottom section with 8 states each
= 6 bits per bar. _Total 120 bit id_.

[spotify id][sid]s are 22 base62 chars,
which appears to be a 128bit uuid (128/log2(62)=21.497)
or better explained on [stack overflow][so]

so where did 8 bits go? idk.

[web]: https://www.spotifycodes.com/index.html#
[sc]: https://user-images.githubusercontent.com/11343221/82091242-4b249a80-96f7-11ea-9ede-3f993708c677.png
[sid]: https://developer.spotify.com/documentation/web-api/
[so]: https://stackoverflow.com/questions/4007280/convert-md5-string-to-base-62-string-in-c
