import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="H3"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const arrow = el.querySelector("g");
  animate(arrow, { rotate: [0, -360], scale: [1, 0.85, 1] }, { duration: τmid, delay: τmin });
}
