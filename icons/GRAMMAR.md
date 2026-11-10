# Icon Motion Grammar

Use this grammar when adding or changing a hover animation.
The set follows 19 hand-tuned reference icons; each new icon should express its meaning in the same motion language.
Animations use Motion 14.7.1 through the shared runtime.

## 1. Mechanics

**One module per icon.**
`anim/<name>.js` animates `svg/<name>.svg`; the file's existence registers the animation.
Paths in this page are relative to `icons/` unless a command says otherwise.

```js
import { animate, τmax } from "../motion.js";

export default function (el) {
  const body = el.querySelector("path");
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
```

- `el` is the root `<svg>`; select parts inside it only.
- Import only the names you use from `../motion.js`; `node --check` must pass.
- Call every `animate()` synchronously inside the function.
  - Sequence with `delay` or `stagger`, never `setTimeout`, `await` or `.then`.
  - Deferred calls count as a second run and escape the runtime's completion tracking.
- Keep no state between hovers, except a `WeakMap` keyed by `el` for `orbit`'s carry.
- Never read or write `data-*` attributes in an animation; morph strips them.

**Selecting parts.**
Select by what the markup already says, in this order of preference:

- Use the tag when unique: `el.querySelector("circle")`.
- Use a geometry attribute: `circle[r="3"]`, `circle[cx="9"]`.
- Use a path prefix or substring: `path[d^="M21 15"]`, `path[d*="21v-2"]`, `path[d$="H8"]`.
- Use an index into `querySelectorAll` when parts are peers (`scroll`, `user-contact`).

**Grouping parts.**
When parts must move together, wrap them in `<g>` in `svg/<name>.svg`.
Examples are the lens and handle of `search`, the lid of `trash` and the whole `bell`.
The importer reports changed SVGs as `edited` and does not overwrite them.
Use classes only as group selectors, such as `g.gear` in `user-cog` and `g.search` in `user-search`; they carry no CSS.

**Pivots.**
`glyph/icons.css` centres rotation and scaling for elements without `transform` or `transform-origin` attributes.
For a different pivot, set `transform-origin="Xpx Ypx"` in view-box units; the CSS switches that element to the view-box.

- `palette` uses `12px 12px`, the icon centre.
- `github` uses `12.01px 11.74px`, the cat's optical centre.
- `bell` uses `12px 2px`, its hanger.
- `trash` uses `12px 6px`, the rim where the lid hinges.
- `map-pin` uses `12px 22px`, so the pin squashes on its tip.

## 2. Time

The shared `glyph/motion.js` exports the timing tokens:

| Token | Seconds | Role                                   |
| ----- | ------- | -------------------------------------- |
| τmax  | 0.42    | Whole gesture; also the hover cooldown |
| τmid  | 0.27    | `τmax − τmin`; main shared beat        |
| τmin  | 0.15    | Flick, settle, nod or final pop        |
| 0.3   | 0.30    | One item of a cascade                  |

**Beats.**
A gesture normally has at most three beats on one clock: `0`, `τmin` and `τmid`, ending by `τmax`.

- **Opening at `0`:** use `τmax`, or `τmid` when an accent follows.
- **Accent at `τmin`:** `{ duration: τmid, delay: τmin }` ends at τmax (`user-plus`, `user-x`, `user-star`).
- **Closer at `τmid`:** `{ duration: τmin, delay: τmid }` ends at τmax (head nods, `check`'s pop).

Use only `τmax`, `τmid`, `τmin` and `0.3` as durations.

**Burst cascades** use `duration: 0.3` with `stagger(s)`, where `s` is 0.03 to 0.21.

- Use `0.06` for small parts (`palette` dots, `clipboard-list` bullets).
- Use `0.09` for signal cascades (`rss`, `wifi`, `volume-2`).
- Use `0.12` for 2 peers: text lines (`scroll`), chevrons (`chevrons-down`) or knobs (`sliders-horizontal`).
- Use `0.21` for 2 far-apart pins (`user-contact`).
- Use `0.03` for 8 rays (`sun`), keeping the cascade near τmax.
- To start later, use `stagger(s, { startDelay: τmin })` or `τmid`.
- The last item should finish by 0.6 s; a cascade on the τmin beat lands at 0.57 s (`coffee` steam, `cake` flames).
  - `github`'s 0.9 s brand flourish is the only longer gesture.

`trash` uses staggered draws: its two slats draw for τmin, staggered by 0.06 from τmid.

**Delayed parts hold their first keyframe.**
Motion applies the first keyframe as soon as hover starts, even when the animation is delayed.
Choose the first keyframe as the pose to hold:

- **Hidden until its beat:** start at `0`; `user-pen` hides its body until τmid, and `house` hides its door.
- **Visible until it moves:** start near rest; `user-pen`'s head holds at `[1, 1, -2, 0]`.
- `pin` drops its head on beat 0 and draws the needle at τmid, hiding the needle while the head falls.

**Easing.**
Do not pass `ease`, `type: "spring"`, `repeat` or `repeatType`.
Motion's default tween is the family's easing; write overshoot into the keyframes (`1.25`, `-2`, `15`).
Keep `times` for cuts and contacts: `send` cuts a jump, and `clapperboard` lands contact near τmin.

## 3. Roles

An icon has up to three roles:

- **Subject:** the noun, such as a body, page, frame or envelope; it draws or breathes.
- **Accent:** the modifier, such as a plus, check, star, gear, pen or lens; it gets the meaning verb.
- **Closer:** a detail that settles last, such as a head, door or dot; it nods or pops at τmid.

Parts with no role stay still.

## 4. Verbs

Each verb has keyframes, timing and reference icons.

**Draw:** `{ pathLength: [0, 1], opacity: [0, 1] }`.

- Use it for a subject, at τmax alone or τmid before an accent (`user`, `user-check`).
- Use `opacity: [0.5, 1]` for frames that must stay readable (`palette`, `user-contact`, `house`).
- **Ring draw:** add `rotate: [270, 90]` to start a full circle at the top (`question`, `info`, `circle-check`).
  - New rings end at 90; `circle-user` ends at 115, and any angle is at rest for a centred circle.
- **Direction:** drawing follows Lucide's path direction; do not use `pathOffset`.
  - If that reads backwards, grow the part from flat or slide it.
  - `chart-line` grows its line; `trending-*` arrows slide; `ruler` ticks grow from the edge.
  - Concentric arcs burst like `wifi`; `rainbow` bursts because its bands draw in opposing directions.
- **Rects:** drawing a `rect` is fine; the runtime removes the dash at hover end, leaving no seam.
  - Prefer a rise (`scaleY: [0, 1.1, 1]`) for a wall going up (`building`).
- **Residue:** after the last animation ends, the runtime restores authored `stroke-dasharray` and `stroke-dashoffset`.
  - Keep Motion's `pathLength="1"`: Motion sets it only once per element, so removing it breaks later draws.
  - It also removes final inline transforms that are pure multiples of 90°, or a pure turn of a centred `circle`.
  - Section 6 limits the turns for which this cleanup is safe.

**Pop:** `{ scale: [0, 1.25, 1] }`, often with a turn.

- Use `rotate: [0, 360]` for asymmetric accents (`user-star`, `user-minus`, `user-cog` gear).
- Use `rotate: [0, 180]` for two-fold symmetric accents (`user-plus`, `user-x`, `plus`, `x`).
- A 180° turn also suits the `snowflake` and each `atom` orbit.
- After a draw, use a **late pop**, `scale: [1, 1.25, 1]`, as a closer (`user-check`, `check`).

**Burst:** `{ scale: [0, 1.75, 1] }` or `[0, 1.5, 1]` with `opacity: [0, 1]`, in a cascade.

- Use it for dots, arcs, rays and glints (`rss`, `palette`, `github`, `wifi`, `sun`).
- A large part overshoots only to 1.25 (`target` rings, `rainbow` bands).

**Grow:** rise from flat on one axis, `scaleY: [0, 1.25, 1]` or `scaleX: [0, 1.25, 1]`.

- Use it for bars, doors and lines that build up (`chart-line`, `house`, `text-align-center`).
- A line ending flush with a frame grows without overshoot, `[0, 1]` (`grid-2x2`, `table`).

**Nod:** `{ translateY: [1, -2, 0] }` at `{ duration: τmin, delay: τmid }`.

- Use it for a person's head; `[1, 1, -2, 0]` holds still first (`user-pen`, `user-search`).
  - The reference `user-search` starts its nod at τmin with duration τmid.
- New people icons follow `contact`, `book-user`, `users-round` and `user-round`.

**Spin:** `{ rotate: [0, 360] }` with a scale dip or rise.

- `settings` uses `rotate: [0, 360], scale: [0, 1.15, 1], opacity: [0, 1]` on the gear.
- `refresh-cw` uses `rotate: [0, 360], scale: [1, 0.85, 1]` on the arrow group.
- `loader-circle` uses `rotate: [0, 720]` with `pathLength: [0.25, 1]`.

**Slide:** translate from an offset into rest, with overshoot.

- **In from aside:** `users`' second person uses `translateX: [3, 0]`; `scroll` text uses `x: [-3, 0], y: [-2, 0]`.
- **Drop in:** `scroll` rolls use `y: [-10, 2, 0]`; `download` uses `translateY: [-6, 2, 0]`; `map-pin` uses `[-6, 1, 0]`.
- **Nudge out and back:** `arrow-right` uses `translateX: [0, 5, -2, 0]`; `chevron-down` uses `translateY: [0, 3, -1, 0]`.
- **Squeeze and return:** `menu` moves its outer bars toward the still middle bar with `translateY: [0, 4, -1, 0]` and `[0, -4, 1, 0]`.
- **Press and settle:** `lock` moves its shackle with `translateY: [0, 1.5, 0]`, never leaving the body.
- Arrows slide along their direction; diagonal arrows move on both axes (`external-link`).

**Swing:** damped rotation about a pivot, `rotate: [-15, 15, -10, 0]`.

- `user-pen` adds `translateX: [0, -5, 0]` for a tool in use.
- `bell` uses `rotate: [0, -18, 14, -8, 0]` from its hanger; the clapper lags by τmin.
- `moon` uses `rotate: [0, -30, 10, 0]`.
- Swing includes sway and rocking (`leaf`, `trees`, `ship`, a ringing `smartphone`).

**Sweep:** a short loop and return, `translateX: [0, -12, 2, 0], translateY: [0, -6, -5, 0]`, with `scale: [1.35, 1]`.

- Use it for looking around (`user-search`, `search` with a smaller loop).

**Flip:** `{ rotateY: [0, 360], scale: [1, 1.35, 1] }`.

- Use it sparingly for a glyph that asks or answers (`question`); `info` slides its glyph instead.

**Breathe:** a scale wave on the subject, `scale: [1, 0.85, 1.15, 1]`.

- Examples are `palette`, `orbit`'s planet `[1, 0.5, 1]`, `heart`'s beat `[1, 1.25, 0.9, 1.15, 1]`, and `folder`'s `scaleY`.
- A badge breathes as an intact seal, without drawing or turning (`badge-check`: `[1, 0.85, 1.1, 1]`).
- Use one axis for squashing or opening: `eye`'s lid `[1, 0.1, 1]`, `mail`'s flap `[-1, 1.15, 1]`, `house`'s door `[0, 1.25, 1]`.

**Blink:** `opacity: [1, 0, 1]` at `{ duration: τmin, delay: τmid }`.

- Only a part that blinks in life uses it; `text-cursor` grows its caret, then blinks it.

**Carry:** rotation continuing between hovers (`orbit` only).

- Do not add another carry; only `orbit` is exempt from returning to rest.

## 5. Meaning Picks the Verb

Pick the verb from the icon's meaning, then apply it to the part carrying that meaning.

| Meaning                  | Verb                              | Reference or exemplar       |
| ------------------------ | --------------------------------- | --------------------------- |
| Add, remove, close       | Pop with a 180° or 360° turn      | user-plus, user-x, plus, x  |
| Confirm, done            | Draw, then late pop               | user-check, check           |
| Favourite, rate          | Pop with a 360° turn              | user-star, star             |
| Settings, sync, loading  | Spin                              | settings, refresh-cw        |
| Edit, write, tool in use | Swing                             | user-pen, pencil            |
| Find, inspect            | Sweep                             | user-search, search         |
| Signal, broadcast, sound | Burst cascade outward             | rss, wifi, volume-2         |
| Direction, transfer      | Slide along the arrow             | arrow-right, download, send |
| Document, container      | Draw at 0.5 opacity; slide lines  | user-contact, file-text     |
| Alert, ring, notify      | Swing from the hanger             | bell                        |
| Help, info               | Ring draw plus flip or slide      | question, info              |
| Create, colour, magic    | Breathe or pop, plus burst        | palette, sparkles           |
| Secure, toggle, reveal   | Squash on one axis, press or lift | lock, eye                   |
| Person                   | Body draw plus head nod           | user, users                 |

## 6. Return to Rest

The checker compares screenshots in Chromium and Firefox and requires **0 differing pixels** after runtime cleanup, except for `orbit`.

- Start every keyframe array explicitly, so a mid-flight re-hover restarts cleanly.
- End every keyframe array at rest: `scale` 1, translation 0, `opacity` 1 and `pathLength` 1.
  - End rotation at a multiple of 360°, or a turn that maps the part onto itself.
  - Two-fold symmetry permits 180°: a line, plus, x, `snowflake` or `atom` orbit.
  - Four-fold symmetry permits 90° (`sparkles`' main spark).
  - A circle pivoting at its centre permits any angle (`circle-user` ends at 115°).
- Runtime cleanup removes final symmetric turns because rotation changes antialiasing.
  - An unsafe 90° or 180° turn would visibly snap back; do not write one.
  - `orbit` keeps its carry because its final turn is not a multiple of 90°.
- `opacity: [1]` is an allowed no-op pinning the end state.
- Do not animate `fill`, `stroke`, `stroke-width` or colour; use geometry, opacity and drawing.
- Every part must be visible and fully drawn at rest.

## 7. Unique but in the Family

- **One meaning verb per icon:** the accent gets it; the subject draws or breathes; the closer nods or pops.
- **Families share a base and vary the accent.**
  - `user-*` icons draw the body and nod the head; only the accent changes.
  - Follow that model for `file-*`, `folder-*`, `mail-*`, `calendar-*`, `bell-*`, `shield-*` and `message-*`.
  - Write the family base first, then pick an accent verb from section 5.
- **Each family member has its own accent gesture.**
  - Vary the turn, beat or verb.
  - `user-plus` and `user-x` share a gesture; reserve that sharing for add/remove pairs.
- **Symmetric pairs mirror offsets** (`download`/`upload`, `arrow-up`/`arrow-down`).
  - Keep their timing related; `upload`'s tray settles earlier than `download`'s.
  - `log-in` and `log-out` share arrow travel; their doors nudge in opposite directions.
  - `log-in`'s door moves at τmid; `log-out`'s moves at 0.
- **Use 2–4 `animate()` calls** unless one expresses the whole gesture.
  - Examples include `chevron-down`, `rss`, `heart`, `send`, `pause` and `hammer`.
- **A shared accent keeps one gesture across families.**
  - Badge and shield alert marks swing like `bell`, and their dots pop at τmid.
  - Their frames differ: `badge-exclamation-point` breathes and `shield-exclamation-point` draws.
  - `triangle-exclamation-point` is the exception: the whole sign sways, its bar draws and its dot pops at τmid.
- **Look-alikes from different families differ in at least one part's verb.**
  - Change the one whose gesture is not its family's base.
  - `text-align-center` grows its lines from the centre; `menu` squeezes its outer bars toward the middle and back.
- **Give the accent the largest overshoot.** The subject draws without one; closers overshoot by at most 2 units or 1.25.

## 8. Judge the Render

Draw and movement may combine when they form one readable gesture.
`loader-circle` spins its arc while it grows: the growing arc expresses loading.
Judge the film, not the property list.

- **Drawing is a choice.** Do not trace away the silhouette that a squeeze, beat, slide or swing needs.
  - `heart`, `folder` and `play` move as complete shapes.
- **Judge uniqueness in the render.** Different keyframes can produce the same gesture.
  - Keep `menu`'s squeeze distinct from `text-align-center`'s growing bars.
- **Do not anticipate with the opposite state.**
  - `face-slightly-smiling` grows its mouth from flat; `scaleY: [-1, …]` would begin with a frown.
- **A physical gesture needs its event:** a clap needs contact.
  - `clapperboard` touches the board once, and the board dips at contact.
- **Count actions a viewer follows**, rather than `animate()` calls.
  - `flower` only grows its petals; extra turns, bursts and centre pops would compete.
- **Secondary draws and closers cost attention.** Delete them when they compete with the meaning verb.
  - `clipboard-list` keeps its clip still, and `battery-full` keeps its terminal still.

## 9. Exemplars

The 19 reference icons are `github`, `orbit`, `palette`, `question`, `rss`, `scroll`, `settings`, `user`, `users`, `user-check`, `circle-user`, `user-cog`, `user-contact`, `user-minus`, `user-pen`, `user-plus`, `user-search`, `user-star` and `user-x`.
Keep these references intact when extending the set.

Read these modules in `anim/` for each verb:

- **Draw:** `house`, `file-text`, `calendar`, `external-link`'s frame.
- **Ring draw:** `circle-check`, `info`.
- **Pop:** `plus`, `x` (staggered bars), `star`, `check` (late pop).
- **Burst:** `wifi`, `rainbow`, `volume-2`, `sun`, `sparkles`.
- **Spin:** `refresh-cw`, `loader-circle`.
- **Slide:** `arrow-right`, `chevron-down`, `download`, `upload`, `send`, `copy`, `map-pin`, `menu`, `lock`.
- **Swing:** `pencil`, `bell`, `moon`, `trash`'s lid.
- **Sweep:** `search`.
- **Breathe:** `heart`, `badge-check`'s seal, `folder`, `eye`, `mail`'s flap.

## 10. Check a New Icon

Run the gallery from the repository root:

```bash
go run ./icons/check -addr 127.0.0.1:18493
```

In another terminal, set `NODE_MODULES` to the absolute path of a directory containing `puppeteer-core`, then run:

```bash
node icons/check/check.mjs --modules "$NODE_MODULES" --only plus,x,check
```

Replace `--only` with your icon names; omit it to check the whole set.
The checker defaults to Chromium at `/usr/bin/chromium` and Firefox at `/usr/bin/firefox-developer-edition`.
`--base` changes the gallery URL; `--hovers N` sets the first-load hover count (default 2).
`--shots <dir>` saves rest and first/repeat image pairs for failures.

The checker enforces:

- A module exists for every requested icon.
- Each hover runs exactly once on first load and after five boosted swaps.
- First-load hovers 2 onward match hover 1 pixel for pixel at 0.1, 0.22 and 0.35 s, except for `orbit`'s random turn.
- Rest matches the authored screenshot with **0 differing pixels**, except for `orbit`'s carry.
- The snap check compares authored rest with the end frame before cleanup, except for `orbit`.
  - At most 16 pixels may have a maximum RGB-channel delta greater than 128.
- Opacity, full drawing, fill, stroke and visibility return to rest, and nothing animates during rest screenshots.
- Reduced motion produces no animation calls.
- Browser console errors, page errors and failed requests fail the check.

Also inspect the movement for the meaning, family and look-alike rules in sections 5–8; the checker cannot judge those.
