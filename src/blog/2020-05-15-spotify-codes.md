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
Each bar has a top and bottom section with 8 states,
which corresponds to 40 bytes of data.
This is more than enough to hold the ~18 bytes of data for a
standard 22 char bas62 [spotify id][sid]

[web]: https://www.spotifycodes.com/index.html#
[sc]: https://user-images.githubusercontent.com/11343221/82091242-4b249a80-96f7-11ea-9ede-3f993708c677.png
[sid]: https://developer.spotify.com/documentation/web-api/
