import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const button = el.querySelectorAll('path[d^="M12 4"], path[d^="M10 4"]');
  animate(button, { translateY: [0, 1.5, 0] }, { duration: τmid });
  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
