import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="H3"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const note = el.querySelector("g");
  animate(note, { rotate: [0, -12, 8, 0] }, { duration: τmid, delay: τmin });
}
