import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelectorAll('path[d^="M10 19"], path[d^="M5 19"], path[d^="M5 5.5"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const barb = el.querySelector('path[d^="m22 19"]');
  animate(barb, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
