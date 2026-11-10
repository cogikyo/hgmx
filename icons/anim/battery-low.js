import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("rect");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cell = el.querySelector('path[d="M6 14v-4"]');
  animate(cell, { opacity: [1, 0, 1] }, { duration: τmin, delay: τmid });
}
