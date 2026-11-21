# Third-Party Notices

fmpickrr bundles and links against the third-party components listed below.
Each is the property of its respective owner and is used under the license
described. Full license texts are included next to the code they cover.

This project itself is released under the MIT License (see `LICENSE`), copyright
Alexey S. Kompaniets. It is a fork of
[`kompadre/fmguessr`](https://github.com/kompadre/fmguessr), same author and
same license, so no upstream attribution transfer is required.

`LICENSE` additionally carries two credits carried over from fmguessr:

- MSFA synth engine (Apache-2.0) by Raph Levien — see `external/msfa` below.
- webdx7 build tooling (MIT) by Jari Kleimola — **not present in this fork's
  tree**; the credit is retained for continuity only.

## Go modules (vendored in `vendor/`)

| Component | Version | License | Copyright |
|---|---|---|---|
| [`github.com/go-audio/audio`](https://github.com/go-audio/audio) | v1.0.0 | Apache-2.0 | 2016 Matt Aimonetti |
| [`github.com/go-audio/wav`](https://github.com/go-audio/wav) | v1.1.0 | Apache-2.0 | 2016 Matt Aimonetti |
| [`github.com/go-audio/riff`](https://github.com/go-audio/riff) | v1.0.0 | Apache-2.0 | 2019 Matt Aimonetti |
| [`github.com/madelynnblue/go-dsp`](https://github.com/madelynnblue/go-dsp) | v1.0.0 | ISC | 2011 Madelynn Blue |

Full license texts are at `vendor/<module>/LICENSE`. `go mod vendor` keeps them
in sync with `go.mod`; re-run it after changing dependencies.

## C/C++ sources (vendored in `external/`)

### `external/msfa` — Music Synthesizer for Android (DX7 engine)

- **Upstream:** <https://github.com/google/music-synthesizer-for-android>
  (archived; originally `code.google.com/p/music-synthesizer-for-android`)
- **License:** Apache-2.0
- **Copyright:** 2012–2013 Google Inc. (authored at Google by Raph Levien, as
  credited in `LICENSE`)

Almost every file carries the Apache-2.0 header inline. The two exceptions are
local C glue written for this project and are covered by this project's MIT
license:

- `external/msfa/wrapper.h`
- `external/msfa/wrapper.cc`

Full Apache-2.0 text: `external/msfa/LICENSE`.

`external/msfa/sin.cc` contains an alternative `#if 0` (never compiled)
polynomial sine approximation attributed in a comment to Basile Graf,
<http://www.rossbencina.com/static/code/sinusoids/even_polynomial_sin_approximation.txt>.

### `external/miniaudio` — audio playback/capture library

- **Upstream:** <https://github.com/mackron/miniaudio>
- **Version:** 0.11.25
- **License:** dual — Unlicense (public domain) **or** MIT-0; your choice
- **Copyright:** David Reid

Full text of both alternatives: `external/miniaudio/LICENSE` (extracted verbatim
from the statement at the end of `miniaudio.h`).

`miniaudio.h` is a single-header library that bundles third-party decoders. All
of them are permissive or public-domain, and all are the choice of Unlicense /
MIT-0 or public domain / MIT:

| Bundled component | Author | License |
|---|---|---|
| dr_wav, dr_flac, dr_mp3 ([dr_libs](https://github.com/mackron/dr_libs)) | David Reid | Unlicense or MIT-0 |
| stb_vorbis ([stb](https://github.com/nothings/stb)) | Sean Barrett | Public domain or MIT |

Other in-file acknowledgements (all public domain, credited in comments):
MurmurHash3 (Peter Scott) and pink-noise generation based on Tonic.

Note: `miniaudio_backend.c` defines `MA_NO_DECODER` and `MA_NO_ENCODER`, but
miniaudio 0.11.25 keys those sections off `MA_NO_DECODING`/`MA_NO_ENCODING`, so
the bundled decoders are still compiled into the binary.

## Provenance notes

- Upstream `fmguessr` lists `/external` in `.gitignore` yet still tracks
  `external/` (a `.gitignore` entry never untracks already-added files). The
  fork's initial commit did not include those files, so a clone could not
  build. Removing the `/external` line here restores tracking so a fresh clone
  builds with no extra downloads.
- `go mod vendor` and `external/` are both committed so the repository builds
  offline with no network access.
