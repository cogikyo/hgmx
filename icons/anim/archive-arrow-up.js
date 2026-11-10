import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [6, -2, 0], opacity: [0, 1] }, { duration: τmax });

  const box = el.querySelectorAll(':scope > path[d$="h2"], :scope > path[d$="h-2"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector("rect");
  animate(lid, { translateY: [0, -2, 0] }, { duration: τmin, delay: τmid });
}
