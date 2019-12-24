package main

import "testing"

func TestImgHack(t *testing.T) {
	tcs := []struct {
		in        string
		html, amp string
	}{
		{
			`<h1><img src="/map.webp" alt="map of countries I&rsquo;ve visited" /></h1>`,
			`
<picture>
        <source type="image/webp" srcset="/map.webp">
        <source type="image/jpeg" srcset="/map.jpg">
        <img src="/map.png" alt="map of countries I&rsquo;ve visited">
</picture>
`,
			`
<amp-img src="/map.webp" alt="map of countries I&rsquo;ve visited" width="1.78" height="1" layout="responsive"></amp-img>
`,
		},
	}
	for i, tc := range tcs {
		html, amp := imgHack(tc.in)
		if html != tc.html {
			t.Errorf("TestImgHack %d: expected %s got %s\n", i, tc.html, html)
		}
		if amp != tc.amp {
			t.Errorf("TestImgHack %d: expected %s got %s\n", i, tc.amp, amp)
		}
	}
}
