import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M13 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const key = el.querySelector("g");
  animate(key, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
