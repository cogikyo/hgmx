import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const flame = el.querySelector('path[d^="M4.5"]');
  animate(flame, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmid });

  const rocket = el.querySelectorAll('path:not([d^="M4.5"])');
  animate(rocket, { translateX: [0, 3, -1, 0], translateY: [0, -3, 1, 0] }, { duration: τmax });
}
