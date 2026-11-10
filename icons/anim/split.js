import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const stem = el.querySelector('path[d^="M12 22"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const branch = el.querySelector("g");
  animate(branch, { translateX: [0, 1.5, -0.5, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
