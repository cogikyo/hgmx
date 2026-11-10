import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector('path[d^="M9 12"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const leg = el.querySelector('path[d^="m11 12"]');
  animate(leg, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
