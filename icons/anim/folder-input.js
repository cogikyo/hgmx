import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M2 9"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
