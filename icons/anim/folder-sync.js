import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M9 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrows = el.querySelector("g");
  animate(arrows, { rotate: [0, 180], scale: [1, 0.7, 1] }, { duration: τmid, delay: τmin });
}
