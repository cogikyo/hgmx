import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cover = el.querySelector("rect");
  animate(cover, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lines = el.querySelectorAll('path[d^="M9.5 "]');
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
