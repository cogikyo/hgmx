import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const slice = el.querySelector('path[d^="M5.5"]');
  animate(slice, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const butter = el.querySelector('path[d^="M7.5"]');
  animate(butter, { scale: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
