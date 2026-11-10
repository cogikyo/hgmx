import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector("path");
  animate(glyph, { scale: [1, 1.2, 1] }, { duration: τmid, delay: τmin });
}
