import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyph = ["M2 5", "M7 2", "m4 14", "m5 8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(glyph, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const letter = el.querySelectorAll('path[d^="m22 22"], path[d^="M14 18"]');
  animate(letter, { translateX: [-6, 1, 0], translateY: [-6, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
