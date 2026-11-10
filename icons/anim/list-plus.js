import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="H3"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const plus = el.querySelectorAll('path[d="M18 9v6"], path[d="M21 12h-6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
