import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector('path[d^="m9 8"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const stem = el.querySelector('path[d^="M12 12"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
