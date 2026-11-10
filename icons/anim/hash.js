import { animate, τmax } from "../motion.js";

export default function (el) {
  const glyph = el.querySelector("g");
  animate(glyph, { scale: [1, 0.9, 1.05, 1] }, { duration: τmax });
}
