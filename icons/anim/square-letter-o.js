import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const glyph = el.querySelector("circle");
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: τmid, delay: τmin });
}
