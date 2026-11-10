import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lever = el.querySelector('path[d="M14 16h4"]');
  animate(lever, { translateY: [0, 2, 0] }, { duration: τmin });

  const toast = el.querySelector('path[d^="M16 10V5.7"]');
  animate(toast, { translateY: [0, -3, 0.5, 0] }, { duration: τmid, delay: τmin });
}
