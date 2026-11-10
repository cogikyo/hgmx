import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector('path[d="M12 20V4"]');
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const bar = el.querySelector('path[d="M5 4h14"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
