import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector("path");
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  animate(glyph, { scaleX: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
